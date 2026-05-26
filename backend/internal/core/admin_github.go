package core

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type githubIssueTarget struct {
	Owner       string
	Repo        string
	IssueNumber int
}

func (t githubIssueTarget) fullName() string {
	return t.Owner + "/" + t.Repo
}

type adminGitHubClient struct {
	token  string
	client *http.Client
}

func newAdminGitHubClient(cfg Config, requireToken bool) (*adminGitHubClient, error) {
	token := strings.TrimSpace(cfg.GitHubToken)
	if requireToken && token == "" {
		return nil, errors.New("GITHUB_TOKEN is required to merge pull requests")
	}
	return &adminGitHubClient{
		token: token,
		client: &http.Client{
			Timeout: 25 * time.Second,
		},
	}, nil
}

func (s *Server) adminTaskPullRequests(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	task, project, ok := s.store.TaskWithProject(r.PathValue("id"))
	if !ok {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	target, err := githubIssueTargetForTask(task, project)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	client, err := newAdminGitHubClient(s.cfg, false)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	pulls, err := client.listPullRequestsLinkedToIssue(r.Context(), target)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, AdminTaskPullRequestsResponse{
		TaskID:       task.ID,
		IssueNumber:  target.IssueNumber,
		IssueURL:     task.IssueURL,
		Repository:   target.fullName(),
		PullRequests: pulls,
	})
}

func (s *Server) mergeAdminTaskPullRequest(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	task, project, ok := s.store.TaskWithProject(r.PathValue("id"))
	if !ok {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	if task.Status == TaskAccepted {
		writeError(w, http.StatusConflict, "task already accepted")
		return
	}
	pullNumber, err := strconv.Atoi(strings.TrimSpace(r.PathValue("number")))
	if err != nil || pullNumber <= 0 {
		writeError(w, http.StatusBadRequest, "pull request number is invalid")
		return
	}
	target, err := githubIssueTargetForTask(task, project)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	client, err := newAdminGitHubClient(s.cfg, true)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	pull, err := client.pullRequest(r.Context(), target, pullNumber)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	if pull.Draft {
		writeError(w, http.StatusConflict, "draft pull requests cannot be merged")
		return
	}
	if !pull.Merged {
		if !strings.EqualFold(pull.State, "open") {
			writeError(w, http.StatusConflict, "pull request is closed without being merged")
			return
		}
		if err := client.mergePullRequest(r.Context(), target, pullNumber); err != nil {
			writeError(w, http.StatusBadGateway, err.Error())
			return
		}
		if refreshed, err := client.pullRequest(r.Context(), target, pullNumber); err == nil {
			pull = refreshed
		}
	}

	req, err := acceptRequestForPullAuthor(task, pull.Author)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	accepted, err := s.store.AcceptTask(task.ID, req)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, AdminMergeTaskPullRequestResponse{
		Task:        accepted,
		PullRequest: pull,
		WorkerID:    req.WorkerID,
	})
}

func githubIssueTargetForTask(task *Task, project *Project) (githubIssueTarget, error) {
	if task == nil {
		return githubIssueTarget{}, errors.New("task not found")
	}
	if target, err := parseGitHubIssueURL(task.IssueURL); err == nil {
		return target, nil
	}
	if project != nil && strings.EqualFold(project.RepoProvider, "github") && task.IssueNumber > 0 {
		for _, candidate := range []string{project.RepoURL, project.BountyRepoName} {
			owner, repo, err := parseGitHubRepo(candidate)
			if err == nil {
				return githubIssueTarget{Owner: owner, Repo: repo, IssueNumber: task.IssueNumber}, nil
			}
		}
	}
	return githubIssueTarget{}, errors.New("task is not tied to a GitHub issue")
}

func parseGitHubIssueURL(value string) (githubIssueTarget, error) {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return githubIssueTarget{}, errors.New("issue url is required")
	}
	parsed, err := url.Parse(raw)
	if err != nil || !strings.EqualFold(parsed.Hostname(), "github.com") {
		return githubIssueTarget{}, errors.New("issue url must be a GitHub URL")
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 4 || !strings.EqualFold(parts[2], "issues") {
		return githubIssueTarget{}, errors.New("issue url must look like https://github.com/owner/repo/issues/123")
	}
	number, err := strconv.Atoi(parts[3])
	if err != nil || number <= 0 {
		return githubIssueTarget{}, errors.New("issue number is invalid")
	}
	owner, repo, err := cleanRepoParts(parts[:2])
	if err != nil {
		return githubIssueTarget{}, err
	}
	return githubIssueTarget{Owner: owner, Repo: repo, IssueNumber: number}, nil
}

func acceptRequestForPullAuthor(task *Task, author string) (AcceptTaskRequest, error) {
	workerID, err := githubWorkerID(author)
	if err != nil {
		return AcceptTaskRequest{}, err
	}
	req := AcceptTaskRequest{
		WorkerKind: task.RequiredWorkerKind,
		WorkerID:   workerID,
	}
	if req.WorkerKind != WorkerHuman {
		req.AgentType = strings.TrimSpace(task.SuggestedAgentType)
		if req.AgentType == "" {
			req.AgentType = "github-pr"
		}
	}
	return req, nil
}

func githubWorkerID(login string) (string, error) {
	login = strings.TrimPrefix(strings.TrimSpace(login), "@")
	if login == "" {
		return "", errors.New("pull request author is required")
	}
	return "github:" + login, nil
}

func (c *adminGitHubClient) listPullRequestsLinkedToIssue(ctx context.Context, target githubIssueTarget) ([]AdminTaskPullRequest, error) {
	seen := map[int]bool{}
	numbers := []int{}
	collect := func(number int) {
		if number <= 0 || seen[number] {
			return
		}
		seen[number] = true
		numbers = append(numbers, number)
	}

	var firstErr error
	if timelineNumbers, err := c.timelinePullNumbers(ctx, target); err == nil {
		for _, number := range timelineNumbers {
			collect(number)
		}
	} else {
		firstErr = err
	}
	if searchNumbers, err := c.searchPullNumbers(ctx, target); err == nil {
		for _, number := range searchNumbers {
			collect(number)
		}
	} else if firstErr == nil {
		firstErr = err
	}
	if len(numbers) == 0 && firstErr != nil {
		return nil, firstErr
	}

	pulls := make([]AdminTaskPullRequest, 0, len(numbers))
	for _, number := range numbers {
		pull, err := c.pullRequest(ctx, target, number)
		if err != nil {
			return nil, err
		}
		pulls = append(pulls, pull)
	}
	sort.SliceStable(pulls, func(i, j int) bool {
		return pulls[i].UpdatedAt.After(pulls[j].UpdatedAt)
	})
	return pulls, nil
}

func (c *adminGitHubClient) timelinePullNumbers(ctx context.Context, target githubIssueTarget) ([]int, error) {
	endpoint := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/issues/%d/timeline?per_page=100",
		url.PathEscape(target.Owner),
		url.PathEscape(target.Repo),
		target.IssueNumber,
	)
	var rows []githubTimelineEvent
	if err := c.githubJSON(ctx, http.MethodGet, endpoint, nil, &rows); err != nil {
		return nil, err
	}
	numbers := []int{}
	for _, row := range rows {
		if row.Source == nil {
			continue
		}
		if row.Source.PullRequest != nil && row.Source.PullRequest.Number > 0 {
			numbers = append(numbers, row.Source.PullRequest.Number)
			continue
		}
		if row.Source.Issue != nil && row.Source.Issue.PullRequest != nil {
			numbers = append(numbers, row.Source.Issue.Number)
		}
	}
	return numbers, nil
}

func (c *adminGitHubClient) searchPullNumbers(ctx context.Context, target githubIssueTarget) ([]int, error) {
	query := fmt.Sprintf("repo:%s/%s type:pr linked:issue #%d", target.Owner, target.Repo, target.IssueNumber)
	endpoint := "https://api.github.com/search/issues?q=" + url.QueryEscape(query) + "&per_page=50"
	var response githubIssueSearchResponse
	if err := c.githubJSON(ctx, http.MethodGet, endpoint, nil, &response); err != nil {
		return nil, err
	}
	numbers := []int{}
	for _, item := range response.Items {
		if item.PullRequest == nil {
			continue
		}
		numbers = append(numbers, item.Number)
	}
	return numbers, nil
}

func (c *adminGitHubClient) pullRequest(ctx context.Context, target githubIssueTarget, number int) (AdminTaskPullRequest, error) {
	endpoint := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/pulls/%d",
		url.PathEscape(target.Owner),
		url.PathEscape(target.Repo),
		number,
	)
	var row githubPullRequestRow
	if err := c.githubJSON(ctx, http.MethodGet, endpoint, nil, &row); err != nil {
		return AdminTaskPullRequest{}, err
	}
	return row.adminRow(), nil
}

func (c *adminGitHubClient) mergePullRequest(ctx context.Context, target githubIssueTarget, number int) error {
	endpoint := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/pulls/%d/merge",
		url.PathEscape(target.Owner),
		url.PathEscape(target.Repo),
		number,
	)
	payload := map[string]any{
		"merge_method": "merge",
		"commit_title": fmt.Sprintf("Merge PR #%d through MergeOS admin", number),
	}
	var result struct {
		Merged  bool   `json:"merged"`
		Message string `json:"message"`
	}
	if err := c.githubJSON(ctx, http.MethodPut, endpoint, payload, &result); err != nil {
		return err
	}
	if !result.Merged {
		message := strings.TrimSpace(result.Message)
		if message == "" {
			message = "GitHub did not merge the pull request"
		}
		return errors.New(message)
	}
	return nil
}

func (c *adminGitHubClient) githubJSON(ctx context.Context, method, endpoint string, body any, out any) error {
	var payload bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&payload).Encode(body); err != nil {
			return err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, &payload)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("github request failed (%d): %s", resp.StatusCode, readBody(resp.Body))
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

type githubTimelineEvent struct {
	Event  string `json:"event"`
	Source *struct {
		Type        string             `json:"type"`
		Issue       *githubLinkedIssue `json:"issue"`
		PullRequest *githubLinkedIssue `json:"pull_request"`
	} `json:"source"`
}

type githubIssueSearchResponse struct {
	Items []githubLinkedIssue `json:"items"`
}

type githubLinkedIssue struct {
	Number      int         `json:"number"`
	HTMLURL     string      `json:"html_url"`
	PullRequest interface{} `json:"pull_request"`
}

type githubPullRequestRow struct {
	Number         int        `json:"number"`
	Title          string     `json:"title"`
	State          string     `json:"state"`
	HTMLURL        string     `json:"html_url"`
	Draft          bool       `json:"draft"`
	Merged         bool       `json:"merged"`
	MergeableState string     `json:"mergeable_state"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	MergedAt       *time.Time `json:"merged_at"`
	User           *struct {
		Login string `json:"login"`
	} `json:"user"`
	Base *struct {
		Ref string `json:"ref"`
	} `json:"base"`
	Head *struct {
		Ref string `json:"ref"`
	} `json:"head"`
}

func (row githubPullRequestRow) adminRow() AdminTaskPullRequest {
	result := AdminTaskPullRequest{
		Number:         row.Number,
		Title:          row.Title,
		State:          row.State,
		HTMLURL:        row.HTMLURL,
		Draft:          row.Draft,
		Merged:         row.Merged,
		MergeableState: row.MergeableState,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
		MergedAt:       row.MergedAt,
	}
	if row.User != nil {
		result.Author = row.User.Login
	}
	if row.Base != nil {
		result.BaseRef = row.Base.Ref
	}
	if row.Head != nil {
		result.HeadRef = row.Head.Ref
	}
	return result
}
