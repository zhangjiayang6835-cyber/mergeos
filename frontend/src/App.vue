<template>
  <div v-if="showPublicShell" class="public-shell">
    <header class="public-topbar">
      <div class="brand-lockup compact">
        <div class="brand-mark">
          <PanelsTopLeft :size="23" />
        </div>
        <div>
          <p class="eyebrow">{{ t('nav.marketplace') }}</p>
          <h1>MergeOS</h1>
        </div>
      </div>
      <nav class="top-tabs public-tabs">
        <button :class="{ active: publicTab === 'talent' }" @click="publicTab = 'talent'">
          <UserRound :size="16" />
          <span>{{ t('nav.talent') }}</span>
        </button>
        <button :class="{ active: publicTab === 'project' }" @click="publicTab = 'project'">
          <FilePlus2 :size="16" />
          <span>{{ t('nav.project') }}</span>
        </button>
        <button :class="{ active: publicTab === 'repo' }" @click="publicTab = 'repo'">
          <GitBranch :size="16" />
          <span>{{ t('nav.repo') }}</span>
        </button>
      </nav>
      <div class="topbar-spacer" />
      <label class="language-control" :title="t('language.title')">
        <Languages :size="16" />
        <select v-model="language" :aria-label="t('language.label')">
          <option v-for="option in languageOptions" :key="option.code" :value="option.code">
            {{ option.flag }} {{ option.label }}
          </option>
        </select>
      </label>
      <button class="action-button ghost public-auth-button" @click="openAuth('login')">
        <LogIn :size="16" />
        <span>{{ t('nav.login') }}</span>
      </button>
      <button class="action-button solid public-auth-button" @click="openAuth('register')">
        <span>{{ t('nav.startOrder') }}</span>
      </button>
    </header>

    <main class="public-main">
      <section class="public-hero">
        <div>
          <p class="eyebrow">{{ t('hero.eyebrow') }}</p>
          <h2>{{ t('hero.title') }}</h2>
          <p>{{ t('hero.body') }}</p>
          <div class="public-actions">
            <button class="action-button solid" @click="publicTab = 'repo'">
              <GitBranch :size="17" />
              <span>{{ t('hero.scoreRepo') }}</span>
            </button>
            <button class="action-button ghost" @click="publicTab = 'talent'">
              <UserRound :size="17" />
              <span>{{ t('hero.findTalent') }}</span>
            </button>
          </div>
        </div>
        <div class="public-signal-board">
          <article>
            <span>{{ t('signal.authLabel') }}</span>
            <strong>{{ t('signal.authValue') }}</strong>
          </article>
          <article>
            <span>{{ t('signal.repoLabel') }}</span>
            <strong>{{ t('signal.repoValue') }}</strong>
          </article>
          <article>
            <span>{{ t('signal.scoringLabel') }}</span>
            <strong>{{ t('signal.scoringValue') }}</strong>
          </article>
        </div>
      </section>

      <section v-if="publicTab === 'talent'" class="public-board talent-board">
        <div class="checkout-panel">
          <div class="panel-heading">
            <UserRound :size="18" />
            <span>{{ t('talent.heading') }}</span>
          </div>
          <div class="talent-controls">
            <label>
              {{ t('talent.search') }}
              <input v-model="talentQuery" :placeholder="t('talent.searchPlaceholder')" />
            </label>
            <label>
              {{ t('talent.skill') }}
              <select v-model="talentSkill">
                <option value="all">{{ t('talent.allSkills') }}</option>
                <option value="frontend">{{ t('skill.frontend') }}</option>
                <option value="backend">{{ t('skill.backend') }}</option>
                <option value="design">{{ t('skill.design') }}</option>
                <option value="qa">{{ t('skill.qa') }}</option>
              </select>
            </label>
          </div>
          <div class="talent-grid">
            <article v-for="talent in filteredTalents" :key="talent.id" class="talent-card">
              <div class="talent-head">
                <span>{{ talent.initials }}</span>
                <div>
                  <strong>{{ talent.name }}</strong>
                  <small>{{ t(talent.roleKey) }}</small>
                </div>
              </div>
              <p>{{ t(talent.summaryKey) }}</p>
              <div class="talent-tags">
                <span v-for="skill in talent.skills" :key="skill">{{ t(`skill.${skill}`) }}</span>
              </div>
              <div class="talent-stats">
                <span>{{ talent.rating }} {{ t('talent.rating') }}</span>
                <span>{{ talent.completed }} {{ t('talent.tasks') }}</span>
                <span>{{ money(talent.rate_cents) }}/{{ t('talent.taskUnit') }}</span>
              </div>
              <button class="action-button ghost card-action" @click="startHireTalent(talent)">
                <CheckCircle2 :size="16" />
                <span>{{ t('talent.invite') }}</span>
              </button>
            </article>
          </div>
        </div>
      </section>

      <section v-if="publicTab === 'project'" class="public-board project-intake-board">
        <div class="checkout-panel">
          <div class="panel-heading">
            <FilePlus2 :size="18" />
            <span>{{ t('project.heading') }}</span>
          </div>
          <div class="project-preview-grid">
            <label>
              {{ t('project.name') }}
              <input v-model="projectForm.title" />
            </label>
            <label>
              {{ t('project.siteType') }}
              <select v-model="projectForm.site_type">
                <option value="Landing page">{{ t('site.landing') }}</option>
                <option value="Business website">{{ t('site.business') }}</option>
                <option value="SaaS website">{{ t('site.saas') }}</option>
                <option value="Storefront">{{ t('site.storefront') }}</option>
                <option value="Web app shell">{{ t('site.webapp') }}</option>
              </select>
            </label>
            <label>
              {{ t('project.budget') }}
              <input v-model.number="budgetUsd" min="100" type="number" />
            </label>
            <label>
              {{ t('project.timeline') }}
              <select v-model="projectForm.timeline">
                <option value="7 days">{{ t('timeline.7') }}</option>
                <option value="14 days">{{ t('timeline.14') }}</option>
                <option value="30 days">{{ t('timeline.30') }}</option>
              </select>
            </label>
          </div>
          <label>
            {{ t('project.brief') }}
            <textarea v-model="projectForm.brief" rows="5" />
          </label>
          <div class="estimate-strip">
            <article>
              <span>{{ t('project.estimatedTasks') }}</span>
              <strong>6</strong>
            </article>
            <article>
              <span>{{ t('project.workerMix') }}</span>
              <strong>{{ t('project.workerMixValue') }}</strong>
            </article>
            <article>
              <span>{{ t('project.checkout') }}</span>
              <strong>{{ money(projectForm.budget_cents) }}</strong>
            </article>
          </div>
          <button class="action-button solid full-action" @click="openAuth('register')">
            <CreditCard :size="17" />
            <span>{{ t('project.continueCheckout') }}</span>
          </button>
        </div>
      </section>

      <section v-if="publicTab === 'repo'" class="public-board repo-import-board">
        <div class="checkout-panel repo-import-panel">
          <div class="panel-heading">
            <GitBranch :size="18" />
            <span>{{ t('repo.heading') }}</span>
          </div>
          <div class="repo-import-form">
            <label>
              {{ t('repo.url') }}
              <input v-model="repoForm.repo_url" placeholder="https://github.com/owner/repo" />
            </label>
            <button class="action-button solid" :disabled="repoImportBusy" @click="importRepoIssues">
              <RefreshCw :size="17" />
              <span>{{ repoImportBusy ? t('repo.loading') : t('repo.load') }}</span>
            </button>
          </div>
          <p v-if="repoImportError" class="error-line">{{ repoImportError }}</p>
          <div v-if="repoImportResult" class="repo-import-summary">
            <article>
              <span>{{ t('repo.repo') }}</span>
              <strong>{{ repoImportResult.owner }}/{{ repoImportResult.name }}</strong>
            </article>
            <article>
              <span>{{ t('repo.openIssues') }}</span>
              <strong>{{ repoImportResult.issue_count }}</strong>
            </article>
            <article>
              <span>{{ t('repo.selectedEstimate') }}</span>
              <strong>{{ money(selectedRepoIssueTotal) }}</strong>
            </article>
          </div>
          <div v-if="repoImportResult?.issues?.length" class="issue-score-list">
            <article
              v-for="issue in repoImportResult.issues"
              :key="issue.number"
              :class="['issue-score-card', { selected: selectedRepoIssueNumbers.includes(issue.number) }]"
            >
              <button class="issue-score-main" @click="toggleRepoIssue(issue)">
                <span class="score-badge">{{ issue.score }}</span>
                <span>
                  <strong>#{{ issue.number }} {{ issue.title }}</strong>
                  <small>{{ issueComplexity(issue.complexity) }} / {{ workerKindLabel(issue.required_worker_kind) }} / {{ money(issue.estimated_cents) }}</small>
                </span>
              </button>
              <div class="issue-score-meta">
                <span v-for="label in issue.labels.slice(0, 4)" :key="label">{{ label }}</span>
                <span>{{ issue.suggested_agent_type || t('repo.humanReview') }}</span>
              </div>
              <p>{{ issueReasonText(issue.reasons) }}</p>
              <a :href="issue.url" target="_blank" rel="noreferrer">
                <ExternalLink :size="15" />
                <span>{{ t('repo.openIssue') }}</span>
              </a>
            </article>
          </div>
          <div v-else class="empty-canvas repo-empty">
            <Database :size="30" />
            <strong>{{ t('repo.emptyTitle') }}</strong>
            <span>{{ t('repo.emptyBody') }}</span>
          </div>
        </div>
        <aside class="order-panel repo-order-preview">
          <div class="panel-heading">
            <CreditCard :size="18" />
            <span>{{ t('repo.orderHeading') }}</span>
          </div>
          <div class="manifest-grid">
            <span>{{ t('repo.selected') }}</span>
            <strong>{{ selectedRepoIssueNumbers.length }} {{ t('repo.issues') }}</strong>
            <span>{{ t('repo.estimated') }}</span>
            <strong>{{ money(selectedRepoIssueTotal) }}</strong>
            <span>{{ t('repo.auth') }}</span>
            <strong>{{ t('repo.authValue') }}</strong>
            <span>{{ t('repo.privateRepo') }}</span>
            <strong>{{ t('repo.privateRepoValue') }}</strong>
          </div>
          <button class="action-button solid full-action" :disabled="selectedRepoIssueNumbers.length === 0" @click="openAuth('register')">
            <CreditCard :size="17" />
            <span>{{ t('repo.checkoutSelected') }}</span>
          </button>
        </aside>
      </section>
    </main>
  </div>

  <div v-else-if="!user" class="auth-shell">
    <section class="auth-panel">
      <div class="brand-lockup">
        <div class="brand-mark">
          <PanelsTopLeft :size="25" />
        </div>
        <div>
          <p class="eyebrow">Client portal</p>
          <h1>MergeOS</h1>
        </div>
      </div>
      <button class="panel-action auth-back-button" @click="backToPublic">
        <span>Back to marketplace</span>
      </button>

      <div class="auth-copy">
        <h2>Fund a website, get a private bounty repo, track every paid task.</h2>
        <p>Register as a client, verify PayPal or crypto payment, and MergeOS converts the budget into MERGE credits for human and agent delivery.</p>
      </div>

      <div class="segmented">
        <button :class="{ active: authMode === 'register' }" @click="authMode = 'register'">Register</button>
        <button :class="{ active: authMode === 'login' }" @click="authMode = 'login'">Login</button>
      </div>

      <form class="auth-form" @submit.prevent="submitAuth">
        <label v-if="authMode === 'register'">
          Full name
          <input v-model="authForm.name" autocomplete="name" />
        </label>
        <label v-if="authMode === 'register'">
          Company
          <input v-model="authForm.company_name" autocomplete="organization" />
        </label>
        <label>
          Email
          <input v-model="authForm.email" autocomplete="email" type="email" />
        </label>
        <label>
          Password
          <input v-model="authForm.password" autocomplete="current-password" type="password" />
        </label>
        <button class="primary-button" :disabled="authBusy">
          <LogIn :size="17" />
          <span>{{ authBusy ? 'Working...' : authMode === 'register' ? 'Create account' : 'Login' }}</span>
        </button>
        <p v-if="errorMessage" class="error-line">{{ errorMessage }}</p>
      </form>

      <div class="auth-runtime">
        <span>{{ runtimeConfig?.payment_mode || 'loading payment' }}</span>
        <span>{{ runtimeConfig?.repo_provider || 'loading repo' }}</span>
        <span>{{ runtimeConfig?.smtp_ready ? 'smtp ready' : 'email log mode' }}</span>
      </div>
    </section>
  </div>

  <div v-else-if="isAdmin" class="admin-shell">
    <header class="topbar">
      <div class="brand-lockup compact">
        <div class="brand-mark">
          <PanelsTopLeft :size="23" />
        </div>
        <div>
          <p class="eyebrow">Admin console</p>
          <h1>MergeOS</h1>
        </div>
      </div>
      <nav class="top-tabs">
        <button :class="{ active: adminTab === 'overview' }" @click="adminTab = 'overview'">
          <LayoutDashboard :size="16" />
          <span>Overview</span>
        </button>
        <button :class="{ active: adminTab === 'projects' }" @click="adminTab = 'projects'">
          <FolderKanban :size="16" />
          <span>Projects</span>
        </button>
        <button :class="{ active: adminTab === 'users' }" @click="adminTab = 'users'">
          <UserRound :size="16" />
          <span>Users</span>
        </button>
        <button :class="{ active: adminTab === 'ledger' }" @click="adminTab = 'ledger'">
          <WalletCards :size="16" />
          <span>Ledger</span>
        </button>
        <button :class="{ active: adminTab === 'inbox' }" @click="adminTab = 'inbox'">
          <Mail :size="16" />
          <span>Email</span>
        </button>
      </nav>
      <div class="topbar-spacer" />
      <div class="status-pill">
        <ShieldCheck :size="16" />
        <span>{{ adminSummary?.repo_provider || runtimeConfig?.repo_provider || 'admin ready' }}</span>
      </div>
      <button class="icon-button" title="Refresh admin data" @click="refreshAll">
        <RefreshCw :size="18" />
      </button>
      <button class="icon-button" title="Logout" @click="logout">
        <LogOut :size="18" />
      </button>
    </header>

    <aside class="admin-sidebar">
      <div class="panel-heading">
        <ShieldCheck :size="18" />
        <span>Admin</span>
      </div>
      <div class="profile-card">
        <strong>{{ user.name }}</strong>
        <span>{{ user.company_name || 'MergeOS' }}</span>
        <small>{{ user.email }}</small>
      </div>
      <div class="runtime-card">
        <div>
          <span>Payment</span>
          <strong>{{ adminSummary?.payment_mode || runtimeConfig?.payment_mode || 'loading' }}</strong>
        </div>
        <div>
          <span>Repo</span>
          <strong>{{ adminSummary?.repo_provider || runtimeConfig?.repo_provider || 'loading' }}</strong>
        </div>
        <div>
          <span>Email</span>
          <strong>{{ adminSummary?.smtp_ready ? 'smtp' : 'log' }}</strong>
        </div>
        <div>
          <span>Token</span>
          <strong>{{ tokenSymbol }}</strong>
        </div>
      </div>
      <div class="project-list compact-list">
        <div class="panel-heading">
          <FolderKanban :size="18" />
          <span>Projects</span>
        </div>
        <button
          v-for="project in projects"
          :key="project.id"
          :class="['project-row', { selected: adminCurrentProject?.id === project.id }]"
          @click="selectAdminProject(project)"
        >
          <span>
            <strong>{{ project.title }}</strong>
            <small>{{ project.client_email }}</small>
          </span>
          <b>{{ money(project.budget_cents) }}</b>
        </button>
      </div>
    </aside>

    <main class="portal-main admin-main">
      <section class="summary-strip">
        <article>
          <span>Funded</span>
          <strong>{{ money(adminSummary?.total_budget_cents) }}</strong>
        </article>
        <article>
          <span>Work pool</span>
          <strong>{{ money(adminSummary?.work_pool_cents) }}</strong>
        </article>
        <article>
          <span>Open tasks</span>
          <strong>{{ adminSummary?.open_task_count || 0 }}</strong>
        </article>
        <article>
          <span>Clients</span>
          <strong>{{ adminSummary?.client_count || 0 }}</strong>
        </article>
      </section>

      <section v-if="adminTab === 'overview'" class="admin-board">
        <div class="checkout-panel">
          <div class="panel-heading">
            <LayoutDashboard :size="18" />
            <span>Operations</span>
          </div>
          <div class="admin-metric-grid">
            <article>
              <span>Projects</span>
              <strong>{{ adminSummary?.project_count || 0 }}</strong>
            </article>
            <article>
              <span>Paid tasks</span>
              <strong>{{ adminSummary?.accepted_task_count || 0 }}</strong>
            </article>
            <article>
              <span>Fees</span>
              <strong>{{ money(adminSummary?.platform_fee_cents) }}</strong>
            </article>
            <article>
              <span>Paid out</span>
              <strong>{{ money(adminSummary?.paid_task_cents) }}</strong>
            </article>
            <article>
              <span>Users</span>
              <strong>{{ adminSummary?.user_count || 0 }}</strong>
            </article>
            <article>
              <span>Files</span>
              <strong>{{ adminSummary?.attachment_count || 0 }}</strong>
            </article>
          </div>
        </div>

        <div class="checkout-panel ssl-panel">
          <div class="panel-heading action-heading">
            <ShieldCheck :size="18" />
            <span>SSL review</span>
            <button class="panel-action" :disabled="sslReviewBusy" title="Review SSL now" @click="reviewSSL">
              <RefreshCw :size="15" />
              <span>{{ sslReviewBusy ? 'Checking...' : 'Review' }}</span>
            </button>
          </div>
          <div v-if="sslReviews.length" class="ssl-domain-list">
            <article v-for="review in sslReviews" :key="review.domain" class="ssl-domain-row">
              <div class="ssl-domain-head">
                <strong>{{ review.domain }}</strong>
                <span :class="['ssl-state', review.status]">{{ sslStatusLabel(review.status) }}</span>
              </div>
              <div class="ssl-facts">
                <span>Expires</span>
                <strong>{{ sslDaysText(review) }}</strong>
                <span>Issuer</span>
                <strong>{{ review.issuer || 'n/a' }}</strong>
                <span>Valid until</span>
                <strong>{{ formatDate(review.not_after) }}</strong>
                <span>Last check</span>
                <strong>{{ formatDate(review.last_checked_at) }}</strong>
              </div>
              <p v-if="review.error" class="ssl-error">{{ review.error }}</p>
            </article>
          </div>
          <p v-else class="muted-line">No SSL domains configured.</p>
        </div>

        <div class="project-list">
          <div class="panel-heading">
            <CheckCircle2 :size="18" />
            <span>Open task queue</span>
          </div>
          <button
            v-for="task in adminOpenTasks"
            :key="task.id"
            :class="['task-row', { selected: adminSelectedTask?.id === task.id }]"
            @click="selectAdminTask(task)"
          >
            <span :class="['status-dot', task.status]" />
            <span>{{ task.title }}</span>
            <strong>{{ money(task.reward_cents) }}</strong>
          </button>
        </div>
      </section>

      <section v-if="adminTab === 'projects'" class="admin-board">
        <div class="project-list">
          <div class="panel-heading">
            <FolderKanban :size="18" />
            <span>Funded projects</span>
          </div>
          <button
            v-for="project in projects"
            :key="project.id"
            :class="['admin-project-row', { selected: adminCurrentProject?.id === project.id }]"
            @click="selectAdminProject(project)"
          >
            <span>
              <strong>{{ project.title }}</strong>
              <small>{{ project.client_name }} / {{ project.client_email }}</small>
            </span>
            <span>{{ project.payment_provider }}</span>
            <b>{{ money(project.budget_cents) }}</b>
          </button>
        </div>

        <div class="checkout-panel">
          <div class="panel-heading">
            <SplitSquareVertical :size="18" />
            <span>Project detail</span>
          </div>
          <div v-if="adminCurrentProject" class="admin-detail">
            <p class="eyebrow">{{ adminCurrentProject.site_type }} / {{ adminCurrentProject.timeline }}</p>
            <h2>{{ adminCurrentProject.title }}</h2>
            <div class="manifest-grid">
              <span>Client</span>
              <strong>{{ adminCurrentProject.client_name }}</strong>
              <span>Company</span>
              <strong>{{ adminCurrentProject.company_name || 'n/a' }}</strong>
              <span>Budget</span>
              <strong>{{ money(adminCurrentProject.budget_cents) }}</strong>
              <span>Work pool</span>
              <strong>{{ money(adminCurrentProject.work_pool_cents) }}</strong>
              <span>Files</span>
              <strong>{{ attachmentCountForProject(adminCurrentProject.id) }}</strong>
              <span>Created</span>
              <strong>{{ formatDate(adminCurrentProject.created_at) }}</strong>
            </div>
            <a v-if="adminCurrentProject.repo_url" class="approval-link" :href="adminCurrentProject.repo_url" target="_blank" rel="noreferrer">
              <ExternalLink :size="16" />
              <span>Open repo</span>
            </a>
            <div class="task-list">
              <button
                v-for="task in adminProjectTasks"
                :key="task.id"
                :class="['task-row', { selected: adminSelectedTask?.id === task.id }]"
                @click="selectAdminTask(task)"
              >
                <span :class="['status-dot', task.status]" />
                <span>{{ task.title }}</span>
                <strong>{{ money(task.reward_cents) }}</strong>
              </button>
            </div>
          </div>
        </div>
      </section>

      <section v-if="adminTab === 'users'" class="checkout-panel">
        <div class="panel-heading">
          <UserRound :size="18" />
          <span>Users</span>
        </div>
        <div class="admin-table">
          <div class="admin-table-head">
            <span>User</span>
            <span>Role</span>
            <span>Projects</span>
            <span>Funded</span>
            <span>Last login</span>
          </div>
          <div v-for="row in adminUsers" :key="row.id" class="admin-table-row">
            <span>
              <strong>{{ row.name }}</strong>
              <small>{{ row.email }}</small>
            </span>
            <span>{{ row.role }}</span>
            <span>{{ row.project_count }}</span>
            <span>{{ money(row.total_budget_cents) }}</span>
            <span>{{ formatDate(row.last_login_at) }}</span>
          </div>
        </div>
      </section>

      <section v-if="adminTab === 'ledger'" class="checkout-panel">
        <div class="panel-heading">
          <WalletCards :size="18" />
          <span>Ledger</span>
        </div>
        <div class="admin-table ledger-table">
          <div class="admin-table-head">
            <span>#</span>
            <span>Type</span>
            <span>Amount</span>
            <span>Reference</span>
            <span>Hash</span>
          </div>
          <div v-for="entry in adminLedgerRows" :key="entry.sequence" class="admin-table-row">
            <span>{{ entry.sequence }}</span>
            <span>{{ entry.type }}</span>
            <span>{{ money(entry.amount_cents) }}</span>
            <span>{{ entry.reference }}</span>
            <span>{{ shortHash(entry.entry_hash) }}</span>
          </div>
        </div>
      </section>

      <section v-if="adminTab === 'inbox'" class="inbox-grid">
        <div class="email-list">
          <div class="panel-heading">
            <Mail :size="18" />
            <span>Customer emails</span>
          </div>
          <article v-for="note in notifications" :key="note.id" class="email-card">
            <span>{{ note.status }}</span>
            <strong>{{ note.subject }}</strong>
            <p>{{ note.body }}</p>
          </article>
        </div>
      </section>
    </main>

    <aside class="inspector admin-inspector">
      <div class="panel-heading">
        <SplitSquareVertical :size="18" />
        <span>Task control</span>
      </div>
      <div v-if="adminSelectedTask" class="task-inspector">
        <p class="eyebrow">{{ adminSelectedTask.status }} / {{ adminSelectedTask.required_worker_kind }}</p>
        <h3>{{ adminSelectedTask.title }}</h3>
        <p>{{ adminSelectedTask.acceptance }}</p>
        <a v-if="adminSelectedTask.issue_url" :href="adminSelectedTask.issue_url" target="_blank" rel="noreferrer">
          <ExternalLink :size="16" />
          <span>Open issue</span>
        </a>
        <div class="manifest-grid">
          <span>Project</span>
          <strong>{{ projectTitle(adminSelectedTask.project_id) }}</strong>
          <span>Worker</span>
          <strong>{{ adminSelectedTask.worker_id || 'pending' }}</strong>
          <span>Reward</span>
          <strong>{{ money(adminSelectedTask.reward_cents) }} {{ tokenSymbol }}</strong>
          <span>Proof</span>
          <strong>{{ adminSelectedTask.proof_hash ? shortHash(adminSelectedTask.proof_hash) : 'pending' }}</strong>
        </div>
        <label>
          Worker kind
          <select v-model="workerForm.worker_kind">
            <option value="human">Human</option>
            <option value="agent">Agent</option>
            <option value="hybrid">Hybrid</option>
          </select>
        </label>
        <label>
          Worker ID
          <input v-model="workerForm.worker_id" placeholder="github:alice or agent:web-001" />
        </label>
        <label>
          Agent type
          <input v-model="workerForm.agent_type" :disabled="workerForm.worker_kind === 'human'" placeholder="frontend-agent" />
        </label>
        <button class="primary-button" :disabled="adminSelectedTask.status === 'accepted' || accepting" @click="acceptAdminSelectedTask">
          <CheckCircle2 :size="17" />
          <span>{{ adminSelectedTask.status === 'accepted' ? 'Paid' : 'Accept and pay' }}</span>
        </button>
        <p v-if="errorMessage" class="error-line">{{ errorMessage }}</p>
      </div>
    </aside>
  </div>

  <div v-else class="app-shell">
    <header class="topbar">
      <div class="brand-lockup compact">
        <div class="brand-mark">
          <PanelsTopLeft :size="23" />
        </div>
        <div>
          <p class="eyebrow">Private client workspace</p>
          <h1>MergeOS</h1>
        </div>
      </div>
      <nav class="top-tabs">
        <button :class="{ active: portalTab === 'workspace' }" @click="portalTab = 'workspace'">
          <LayoutDashboard :size="16" />
          <span>Workspace</span>
        </button>
        <button :class="{ active: portalTab === 'billing' }" @click="portalTab = 'billing'">
          <WalletCards :size="16" />
          <span>Billing</span>
        </button>
        <button :class="{ active: portalTab === 'inbox' }" @click="portalTab = 'inbox'">
          <Mail :size="16" />
          <span>Email</span>
        </button>
      </nav>
      <div class="topbar-spacer" />
      <div class="status-pill">
        <ShieldCheck :size="16" />
        <span>{{ statusLabel }}</span>
      </div>
      <button class="icon-button" title="Refresh workspace" @click="refreshAll">
        <RefreshCw :size="18" />
      </button>
      <button class="icon-button" title="Logout" @click="logout">
        <LogOut :size="18" />
      </button>
    </header>

    <aside class="customer-panel">
      <div class="panel-heading">
        <UserRound :size="18" />
        <span>Customer</span>
      </div>

      <div class="profile-card">
        <strong>{{ user.name }}</strong>
        <span>{{ user.company_name || 'Independent client' }}</span>
        <small>{{ user.email }}</small>
      </div>

      <div class="runtime-card">
        <div>
          <span>Payment</span>
          <strong>{{ runtimeConfig?.payment_mode || 'loading' }}</strong>
        </div>
        <div>
          <span>Repo</span>
          <strong>{{ runtimeConfig?.repo_provider || 'loading' }}</strong>
        </div>
        <div>
          <span>Email</span>
          <strong>{{ runtimeConfig?.smtp_ready ? 'smtp' : 'log' }}</strong>
        </div>
        <div>
          <span>Token</span>
          <strong>{{ tokenSymbol }}</strong>
        </div>
      </div>

      <label>
        Contact name
        <input v-model="projectForm.client_name" />
      </label>
      <label>
        Contact email
        <input v-model="projectForm.client_email" type="email" />
      </label>
      <label>
        Phone
        <input v-model="projectForm.phone" />
      </label>
      <label>
        Company
        <input v-model="projectForm.company_name" />
      </label>
    </aside>

    <main class="portal-main">
      <section class="summary-strip">
        <article>
          <span>Total funded</span>
          <strong>{{ money(totalBudget) }}</strong>
        </article>
        <article>
          <span>MERGE reserved</span>
          <strong>{{ money(totalPool) }}</strong>
        </article>
        <article>
          <span>Open tasks</span>
          <strong>{{ openTasks.length }}</strong>
        </article>
        <article>
          <span>Paid tasks</span>
          <strong>{{ acceptedTasks.length }}</strong>
        </article>
      </section>

      <section v-if="portalTab === 'workspace'" class="workspace-grid">
        <div class="canvas-column">
          <div class="canvas-toolbar">
            <div>
              <p class="eyebrow">Website order</p>
              <h2>{{ currentProject?.title || 'Create a funded website project' }}</h2>
            </div>
            <div class="toolbar-metrics">
              <span>{{ currentTasks.length }} issues</span>
              <span>{{ currentProject?.repo_provider || runtimeConfig?.repo_provider || 'local-git' }}</span>
            </div>
          </div>

          <section class="builder-canvas">
            <div class="canvas-section hero-section">
              <div class="section-handle">BRIEF</div>
              <div>
                <p class="eyebrow">{{ projectForm.site_type }} / {{ projectForm.timeline }}</p>
                <h3>{{ currentProject?.company_name || projectForm.company_name }} delivery room</h3>
                <p>{{ currentProject?.brief || projectForm.brief }}</p>
              </div>
              <div class="quote-block">
                <span>{{ currentProject?.payment_provider || 'checkout pending' }}</span>
                <strong>{{ money(currentProject?.budget_cents || projectForm.budget_cents) }}</strong>
              </div>
            </div>

            <div v-if="currentProject?.attachments?.length" class="attachment-preview">
              <button
                v-for="attachment in currentProject.attachments"
                :key="attachment.id"
                type="button"
                class="attachment-chip"
                @click="openAttachment(attachment)"
              >
                <FileImage v-if="attachment.is_image" :size="18" />
                <Paperclip v-else :size="18" />
                <span>{{ attachment.original_name }}</span>
              </button>
            </div>

            <div v-if="currentTasks.length" class="canvas-grid">
              <button
                v-for="task in currentTasks"
                :key="task.id"
                :class="['task-tile', { selected: selectedTask?.id === task.id, accepted: task.status === 'accepted' }]"
                @click="selectTask(task)"
              >
                <span class="issue-label">Issue #{{ task.issue_number }}</span>
                <strong>{{ task.title }}</strong>
                <small>{{ task.required_worker_kind }} / {{ money(task.reward_cents) }} {{ tokenSymbol }}</small>
              </button>
            </div>

            <div v-else class="empty-canvas">
              <Database :size="30" />
              <strong>No funded bounty yet</strong>
              <span>{{ runtimeConfig?.dev_payment_enabled ? `Use ${runtimeConfig.dev_payment_code} as the local payment reference.` : 'Configure PayPal or crypto credentials first.' }}</span>
            </div>
          </section>

          <section class="repo-strip">
            <div class="strip-header">
              <GitBranch :size="18" />
              <a v-if="currentProject?.repo_url" :href="currentProject.repo_url" target="_blank" rel="noreferrer">
                {{ currentProject.bounty_repo_name }}
              </a>
              <span v-else>mergeos-bounties repo pending</span>
            </div>
            <div class="task-list">
              <button
                v-for="task in currentTasks"
                :key="task.id"
                :class="['task-row', { selected: selectedTask?.id === task.id }]"
                @click="selectTask(task)"
              >
                <span :class="['status-dot', task.status]" />
                <span>{{ task.title }}</span>
                <strong>{{ money(task.reward_cents) }}</strong>
              </button>
            </div>
          </section>
        </div>

        <aside class="order-panel">
          <div class="panel-heading">
            <FilePlus2 :size="18" />
            <span>New project</span>
          </div>
          <label>
            Project name
            <input v-model="projectForm.title" />
          </label>
          <label>
            Site type
            <select v-model="projectForm.site_type">
              <option>Landing page</option>
              <option>Business website</option>
              <option>SaaS website</option>
              <option>Storefront</option>
              <option>Web app shell</option>
            </select>
          </label>
          <label>
            Package
            <select v-model="projectForm.package_tier">
              <option>Launch</option>
              <option>Growth</option>
              <option>Scale</option>
            </select>
          </label>
          <label>
            Timeline
            <select v-model="projectForm.timeline">
              <option>7 days</option>
              <option>14 days</option>
              <option>30 days</option>
            </select>
          </label>
          <label>
            Budget USD
            <input v-model.number="budgetUsd" min="100" type="number" />
          </label>
          <label>
            Brief
            <textarea v-model="projectForm.brief" rows="5" />
          </label>
          <label class="upload-control">
            Reference files
            <input type="file" multiple @change="uploadProjectFiles" />
            <span class="upload-surface">
              <UploadCloud :size="18" />
              <span>{{ uploadBusy ? 'Uploading...' : 'Add images or files' }}</span>
            </span>
          </label>
          <div v-if="uploadedAttachments.length" class="attachment-list pending">
            <div v-for="attachment in uploadedAttachments" :key="attachment.id" class="attachment-row">
              <FileImage v-if="attachment.is_image" :size="17" />
              <Paperclip v-else :size="17" />
              <span>
                <strong>{{ attachment.original_name }}</strong>
                <small>{{ attachment.content_type }} / {{ fileSize(attachment.size_bytes) }}</small>
              </span>
              <button class="remove-attachment" title="Remove file" @click="removeUploadedAttachment(attachment.id)">
                <X :size="15" />
              </button>
            </div>
          </div>
          <button class="primary-button" :disabled="creating || uploadBusy" @click="createProject">
            <CreditCard :size="17" />
            <span>{{ creating ? 'Funding...' : 'Verify payment and create repo' }}</span>
          </button>
          <p v-if="errorMessage" class="error-line">{{ errorMessage }}</p>
        </aside>
      </section>

      <section v-if="portalTab === 'billing'" class="billing-grid">
        <div class="checkout-panel">
          <div class="panel-heading">
            <WalletCards :size="18" />
            <span>Checkout</span>
          </div>
          <div class="payment-choice">
            <button :class="{ active: projectForm.payment_method === 'paypal' }" @click="projectForm.payment_method = 'paypal'">PayPal</button>
            <button :class="{ active: projectForm.payment_method === 'crypto' }" @click="projectForm.payment_method = 'crypto'">Crypto</button>
          </div>
          <button
            v-if="projectForm.payment_method === 'paypal'"
            class="secondary-button"
            :disabled="!runtimeConfig?.paypal_ready || preparingPayPal"
            @click="preparePayPalOrder"
          >
            <WalletCards :size="17" />
            <span>{{ preparingPayPal ? 'Creating order...' : 'Create PayPal order' }}</span>
          </button>
          <a v-if="paypalOrder?.approval_url" class="approval-link" :href="paypalOrder.approval_url" target="_blank" rel="noreferrer">
            <ExternalLink :size="16" />
            <span>Open PayPal approval</span>
          </a>
          <label>
            Payment reference
            <input v-model="projectForm.payment_reference" :placeholder="paymentReferencePlaceholder" />
          </label>
          <div v-if="projectForm.payment_method === 'crypto'" class="receiver-card">
            <span>Receiver</span>
            <strong>{{ runtimeConfig?.crypto_receiver || 'Configure CRYPTO_RECEIVER in backend env' }}</strong>
          </div>
          <div class="billing-ledger">
            <div v-for="entry in recentLedger" :key="entry.sequence" class="ledger-line">
              <span>#{{ entry.sequence }} {{ entry.type }}</span>
              <strong>{{ money(entry.amount_cents) }}</strong>
            </div>
          </div>
        </div>

        <div class="project-list">
          <div class="panel-heading">
            <FolderKanban :size="18" />
            <span>Funded projects</span>
          </div>
          <button
            v-for="project in projects"
            :key="project.id"
            :class="['project-row', { selected: currentProject?.id === project.id }]"
            @click="selectProject(project)"
          >
            <span>
              <strong>{{ project.title }}</strong>
              <small>{{ project.bounty_repo_name }}</small>
            </span>
            <b>{{ money(project.budget_cents) }}</b>
          </button>
        </div>
      </section>

      <section v-if="portalTab === 'inbox'" class="inbox-grid">
        <div class="email-list">
          <div class="panel-heading">
            <Mail :size="18" />
            <span>Customer emails</span>
          </div>
          <article v-for="note in notifications" :key="note.id" class="email-card">
            <span>{{ note.status }}</span>
            <strong>{{ note.subject }}</strong>
            <p>{{ note.body }}</p>
          </article>
        </div>
      </section>
    </main>

    <aside class="inspector">
      <div class="panel-heading">
        <SplitSquareVertical :size="18" />
        <span>Task inspector</span>
      </div>

      <div class="repo-summary">
        <p class="eyebrow">Child bounty repo</p>
        <h3>{{ currentProject?.bounty_repo_name || 'Not created' }}</h3>
        <p>{{ currentProject?.repo_provider || runtimeConfig?.repo_provider || 'local-git' }}</p>
        <a v-if="currentProject?.repo_url" :href="currentProject.repo_url" target="_blank" rel="noreferrer">
          <ExternalLink :size="16" />
          <span>Open repo</span>
        </a>
        <div v-if="currentProject?.attachments?.length" class="repo-attachments">
          <button
            v-for="attachment in currentProject.attachments"
            :key="attachment.id"
            type="button"
            @click="openAttachment(attachment)"
          >
            <FileImage v-if="attachment.is_image" :size="16" />
            <Paperclip v-else :size="16" />
            <span>{{ attachment.original_name }}</span>
          </button>
        </div>
      </div>

      <div v-if="selectedTask" class="task-inspector">
        <p class="eyebrow">Selected issue</p>
        <h3>{{ selectedTask.title }}</h3>
        <p>{{ selectedTask.acceptance }}</p>
        <a v-if="selectedTask.issue_url" :href="selectedTask.issue_url" target="_blank" rel="noreferrer">
          <ExternalLink :size="16" />
          <span>Open issue</span>
        </a>

        <div class="manifest-grid">
          <span>Required</span>
          <strong>{{ selectedTask.required_worker_kind }}</strong>
          <span>Suggested</span>
          <strong>{{ selectedTask.suggested_agent_type || 'human-review' }}</strong>
          <span>Reward</span>
          <strong>{{ money(selectedTask.reward_cents) }} {{ tokenSymbol }}</strong>
          <span>Proof</span>
          <strong>{{ selectedTask.proof_hash ? shortHash(selectedTask.proof_hash) : 'pending' }}</strong>
        </div>

        <label>
          Worker kind
          <select v-model="workerForm.worker_kind">
            <option value="human">Human</option>
            <option value="agent">Agent</option>
            <option value="hybrid">Hybrid</option>
          </select>
        </label>
        <label>
          Worker ID
          <input v-model="workerForm.worker_id" placeholder="github:alice or agent:web-001" />
        </label>
        <label>
          Agent type
          <input v-model="workerForm.agent_type" :disabled="workerForm.worker_kind === 'human'" placeholder="frontend-agent" />
        </label>
        <button class="primary-button" :disabled="selectedTask.status === 'accepted' || accepting" @click="acceptSelectedTask">
          <CheckCircle2 :size="17" />
          <span>{{ selectedTask.status === 'accepted' ? 'Paid' : 'Accept and pay' }}</span>
        </button>
      </div>
    </aside>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue';
import {
  CheckCircle2,
  CreditCard,
  Database,
  ExternalLink,
  FileImage,
  FilePlus2,
  FolderKanban,
  GitBranch,
  Languages,
  LayoutDashboard,
  LogIn,
  LogOut,
  Mail,
  PanelsTopLeft,
  Paperclip,
  RefreshCw,
  ShieldCheck,
  SplitSquareVertical,
  UploadCloud,
  UserRound,
  WalletCards,
  X,
} from '@lucide/vue/dist/esm/lucide-vue.mjs';

const runtimeConfig = ref(null);
const user = ref(null);
const publicTab = ref('talent');
const authVisible = ref(false);
const authMode = ref('register');
const hasWindow = typeof window !== 'undefined';
const token = ref(hasWindow ? localStorage.getItem('mergeos_token') || '' : '');
const language = ref('en');
const authBusy = ref(false);
const talentQuery = ref('');
const talentSkill = ref('all');
const projects = ref([]);
const tasks = ref([]);
const ledger = ref([]);
const notifications = ref([]);
const adminSummary = ref(null);
const adminUsers = ref([]);
const attachments = ref([]);
const sslReviews = ref([]);
const selectedProjectId = ref('');
const selectedTask = ref(null);
const selectedAdminProjectId = ref('');
const adminSelectedTask = ref(null);
const portalTab = ref('workspace');
const adminTab = ref('overview');
const creating = ref(false);
const accepting = ref(false);
const preparingPayPal = ref(false);
const uploadBusy = ref(false);
const sslReviewBusy = ref(false);
const errorMessage = ref('');
const paypalOrder = ref(null);
const uploadedAttachments = ref([]);
const repoImportBusy = ref(false);
const repoImportError = ref('');
const repoImportResult = ref(null);
const selectedRepoIssueNumbers = ref([]);

const languageOptions = [
  { code: 'en', flag: '🇺🇸', label: 'English' },
  { code: 'zh', flag: '🇨🇳', label: '中文' },
  { code: 'ja', flag: '🇯🇵', label: '日本語' },
  { code: 'ko', flag: '🇰🇷', label: '한국어' },
  { code: 'vi', flag: '🇻🇳', label: 'Tiếng Việt' },
];

const localeByLanguage = {
  en: 'en-US',
  zh: 'zh-CN',
  ja: 'ja-JP',
  ko: 'ko-KR',
  vi: 'vi-VN',
};

const translations = {
  en: {
    'nav.marketplace': 'Public marketplace',
    'nav.talent': 'Talent',
    'nav.project': 'New project',
    'nav.repo': 'Fix repo issues',
    'nav.login': 'Login',
    'nav.startOrder': 'Start order',
    'language.label': 'Language',
    'language.title': 'Choose interface language',
    'hero.eyebrow': 'Hire humans and agents by outcome',
    'hero.title': 'Find talent, fund a new build, or turn repo issues into scored paid work.',
    'hero.body': 'Guests can browse talent and analyze public GitHub issues first. Account creation only appears when you are ready to hire, checkout, connect private repos, or track delivery.',
    'hero.scoreRepo': 'Score repo issues',
    'hero.findTalent': 'Find talent',
    'signal.authLabel': 'Auth gate',
    'signal.authValue': 'Only at checkout',
    'signal.repoLabel': 'Repo import',
    'signal.repoValue': 'Public GitHub issues',
    'signal.scoringLabel': 'Scoring',
    'signal.scoringValue': 'Complexity, bounty, worker mix',
    'talent.heading': 'Talent marketplace',
    'talent.search': 'Search',
    'talent.searchPlaceholder': 'frontend, checkout, QA...',
    'talent.skill': 'Skill',
    'talent.allSkills': 'All skills',
    'talent.rating': 'rating',
    'talent.tasks': 'tasks',
    'talent.taskUnit': 'task',
    'talent.invite': 'Invite to order',
    'skill.frontend': 'Frontend',
    'skill.backend': 'Backend',
    'skill.design': 'Design',
    'skill.qa': 'QA',
    'skill.ui': 'UI',
    'skill.accessibility': 'Accessibility',
    'skill.payments': 'Payments',
    'skill.api': 'API',
    'skill.copy': 'Copy',
    'talentProfile.frontend.role': 'Frontend delivery pod',
    'talentProfile.frontend.summary': 'Ships responsive Vue/React surfaces, checkout states, dashboards, and accessibility passes.',
    'talentProfile.backend.role': 'Go API and payment engineers',
    'talentProfile.backend.summary': 'Handles auth, webhooks, payment verification, data migrations, and proof-ledger fixes.',
    'talentProfile.design.role': 'Human review team',
    'talentProfile.design.summary': 'Checks product clarity, copy, responsive layout, contrast, and release readiness.',
    'talentProfile.repo.role': 'Issue triage and PR agents',
    'talentProfile.repo.summary': 'Imports issue context, opens scoped PRs, and pairs with human review for risky changes.',
    'project.heading': 'Project order preview',
    'project.name': 'Project name',
    'project.siteType': 'Site type',
    'project.budget': 'Budget USD',
    'project.timeline': 'Timeline',
    'project.brief': 'Brief',
    'project.estimatedTasks': 'Estimated tasks',
    'project.workerMix': 'Worker mix',
    'project.workerMixValue': 'Human + agent',
    'project.checkout': 'Checkout',
    'project.continueCheckout': 'Continue to checkout',
    'site.landing': 'Landing page',
    'site.business': 'Business website',
    'site.saas': 'SaaS website',
    'site.storefront': 'Storefront',
    'site.webapp': 'Web app shell',
    'timeline.7': '7 days',
    'timeline.14': '14 days',
    'timeline.30': '30 days',
    'repo.heading': 'Fix issues in an existing repo',
    'repo.url': 'GitHub repo URL',
    'repo.loading': 'Loading issues...',
    'repo.load': 'Load and score issues',
    'repo.repo': 'Repo',
    'repo.openIssues': 'Open issues',
    'repo.selectedEstimate': 'Selected estimate',
    'repo.humanReview': 'human-review',
    'repo.openIssue': 'Open issue',
    'repo.emptyTitle': 'Paste a public GitHub repo',
    'repo.emptyBody': 'MergeOS will import open issues, skip pull requests, score complexity, estimate bounty, and suggest human or agent work.',
    'repo.orderHeading': 'Issue order',
    'repo.selected': 'Selected',
    'repo.issues': 'issues',
    'repo.estimated': 'Estimated',
    'repo.auth': 'Auth',
    'repo.authValue': 'Required at checkout',
    'repo.privateRepo': 'Private repo',
    'repo.privateRepoValue': 'Connect GitHub after login',
    'repo.checkoutSelected': 'Checkout selected issues',
    'complexity.low': 'low',
    'complexity.medium': 'medium',
    'complexity.high': 'high',
    'worker.human': 'human',
    'worker.agent': 'agent',
    'worker.hybrid': 'hybrid',
    'reason.openGitHubIssue': 'open GitHub issue',
    'reason.detailedIssueBody': 'detailed issue body',
    'reason.clearReproductionContext': 'clear reproduction context',
    'reason.activeDiscussion': 'active discussion',
    'reason.securityOrAuthRisk': 'security or auth risk',
    'reason.productionRisk': 'production risk',
    'reason.bugFix': 'bug fix',
    'reason.backendSurface': 'backend surface',
    'reason.frontendSurface': 'frontend surface',
    'reason.scopeExpansion': 'scope expansion',
    'reason.smallEditorialTask': 'small editorial task',
    'reason.lowComplexityLabel': 'low complexity label',
  },
  zh: {
    'nav.marketplace': '公开市场',
    'nav.talent': '人才',
    'nav.project': '新项目',
    'nav.repo': '修复仓库问题',
    'nav.login': '登录',
    'nav.startOrder': '开始下单',
    'language.label': '语言',
    'language.title': '选择界面语言',
    'hero.eyebrow': '按结果雇用人类与智能体',
    'hero.title': '寻找人才、资助新项目，或把仓库 issue 转成可评分的付费任务。',
    'hero.body': '访客可以先浏览人才并分析公开 GitHub issue。只有在雇用、结账、连接私有仓库或跟踪交付时才需要创建账号。',
    'hero.scoreRepo': '评分仓库 issue',
    'hero.findTalent': '寻找人才',
    'signal.authLabel': '登录门槛',
    'signal.authValue': '仅结账时需要',
    'signal.repoLabel': '仓库导入',
    'signal.repoValue': '公开 GitHub issue',
    'signal.scoringLabel': '评分',
    'signal.scoringValue': '复杂度、赏金、工作类型',
    'talent.heading': '人才市场',
    'talent.search': '搜索',
    'talent.searchPlaceholder': '前端、结账、QA...',
    'talent.skill': '技能',
    'talent.allSkills': '全部技能',
    'talent.rating': '评分',
    'talent.tasks': '任务',
    'talent.taskUnit': '任务',
    'talent.invite': '邀请到订单',
    'skill.frontend': '前端',
    'skill.backend': '后端',
    'skill.design': '设计',
    'skill.qa': '质量检查',
    'skill.ui': '界面',
    'skill.accessibility': '无障碍',
    'skill.payments': '支付',
    'skill.api': 'API',
    'skill.copy': '文案',
    'talentProfile.frontend.role': '前端交付团队',
    'talentProfile.frontend.summary': '交付响应式 Vue/React 界面、结账状态、仪表盘与无障碍检查。',
    'talentProfile.backend.role': 'Go API 与支付工程师',
    'talentProfile.backend.summary': '处理认证、webhook、支付验证、数据迁移和证明账本修复。',
    'talentProfile.design.role': '人工设计审核团队',
    'talentProfile.design.summary': '检查产品清晰度、文案、响应式布局、对比度和发布准备度。',
    'talentProfile.repo.role': 'Issue 分诊与 PR 智能体',
    'talentProfile.repo.summary': '导入 issue 上下文，提交限定范围 PR，并为高风险改动配合人工审核。',
    'project.heading': '项目订单预览',
    'project.name': '项目名称',
    'project.siteType': '网站类型',
    'project.budget': '预算 USD',
    'project.timeline': '周期',
    'project.brief': '需求说明',
    'project.estimatedTasks': '预计任务',
    'project.workerMix': '工作组合',
    'project.workerMixValue': '人类 + 智能体',
    'project.checkout': '结账',
    'project.continueCheckout': '继续结账',
    'site.landing': '落地页',
    'site.business': '企业网站',
    'site.saas': 'SaaS 网站',
    'site.storefront': '店铺网站',
    'site.webapp': 'Web App 外壳',
    'timeline.7': '7 天',
    'timeline.14': '14 天',
    'timeline.30': '30 天',
    'repo.heading': '修复现有仓库中的 issue',
    'repo.url': 'GitHub 仓库 URL',
    'repo.loading': '正在加载 issue...',
    'repo.load': '加载并评分 issue',
    'repo.repo': '仓库',
    'repo.openIssues': '开放 issue',
    'repo.selectedEstimate': '已选预估',
    'repo.humanReview': '人工审核',
    'repo.openIssue': '打开 issue',
    'repo.emptyTitle': '粘贴公开 GitHub 仓库',
    'repo.emptyBody': 'MergeOS 会导入开放 issue，跳过 PR，评分复杂度，预估赏金，并建议人类或智能体执行。',
    'repo.orderHeading': 'Issue 订单',
    'repo.selected': '已选',
    'repo.issues': '个 issue',
    'repo.estimated': '预估',
    'repo.auth': '登录',
    'repo.authValue': '结账时需要',
    'repo.privateRepo': '私有仓库',
    'repo.privateRepoValue': '登录后连接 GitHub',
    'repo.checkoutSelected': '结账所选 issue',
    'complexity.low': '低',
    'complexity.medium': '中',
    'complexity.high': '高',
    'worker.human': '人类',
    'worker.agent': '智能体',
    'worker.hybrid': '混合',
    'reason.openGitHubIssue': '开放 GitHub issue',
    'reason.detailedIssueBody': 'issue 内容详细',
    'reason.clearReproductionContext': '复现上下文清楚',
    'reason.activeDiscussion': '讨论活跃',
    'reason.securityOrAuthRisk': '安全或认证风险',
    'reason.productionRisk': '生产风险',
    'reason.bugFix': '缺陷修复',
    'reason.backendSurface': '后端范围',
    'reason.frontendSurface': '前端范围',
    'reason.scopeExpansion': '范围扩展',
    'reason.smallEditorialTask': '小型文案任务',
    'reason.lowComplexityLabel': '低复杂度标签',
  },
  ja: {
    'nav.marketplace': '公開マーケット',
    'nav.talent': 'タレント',
    'nav.project': '新規プロジェクト',
    'nav.repo': 'Issue 修正',
    'nav.login': 'ログイン',
    'nav.startOrder': '注文を開始',
    'language.label': '言語',
    'language.title': '表示言語を選択',
    'hero.eyebrow': '成果単位で人とエージェントを採用',
    'hero.title': 'タレントを探し、新規制作を依頼し、リポジトリの issue を採点済みの有償タスクに変換。',
    'hero.body': 'ゲストは先にタレント閲覧と公開 GitHub issue の分析ができます。アカウント作成は採用、決済、非公開リポジトリ接続、納品追跡の時だけ必要です。',
    'hero.scoreRepo': 'Issue を採点',
    'hero.findTalent': 'タレントを探す',
    'signal.authLabel': '認証',
    'signal.authValue': '決済時のみ',
    'signal.repoLabel': 'Repo インポート',
    'signal.repoValue': '公開 GitHub issue',
    'signal.scoringLabel': '採点',
    'signal.scoringValue': '複雑度、報酬、作業種別',
    'talent.heading': 'タレントマーケット',
    'talent.search': '検索',
    'talent.searchPlaceholder': 'フロントエンド、決済、QA...',
    'talent.skill': 'スキル',
    'talent.allSkills': 'すべてのスキル',
    'talent.rating': '評価',
    'talent.tasks': 'タスク',
    'talent.taskUnit': 'タスク',
    'talent.invite': '注文に招待',
    'skill.frontend': 'フロントエンド',
    'skill.backend': 'バックエンド',
    'skill.design': 'デザイン',
    'skill.qa': 'QA',
    'skill.ui': 'UI',
    'skill.accessibility': 'アクセシビリティ',
    'skill.payments': '決済',
    'skill.api': 'API',
    'skill.copy': 'コピー',
    'talentProfile.frontend.role': 'フロントエンド納品チーム',
    'talentProfile.frontend.summary': 'レスポンシブな Vue/React 画面、決済状態、ダッシュボード、アクセシビリティ確認を納品します。',
    'talentProfile.backend.role': 'Go API と決済エンジニア',
    'talentProfile.backend.summary': '認証、webhook、決済検証、データ移行、証跡台帳の修正を担当します。',
    'talentProfile.design.role': '人間のレビュー チーム',
    'talentProfile.design.summary': '製品の明確さ、コピー、レスポンシブ表示、コントラスト、リリース準備を確認します。',
    'talentProfile.repo.role': 'Issue 分類と PR エージェント',
    'talentProfile.repo.summary': 'Issue の文脈を取り込み、範囲を絞った PR を作成し、リスクの高い変更は人間がレビューします。',
    'project.heading': 'プロジェクト注文プレビュー',
    'project.name': 'プロジェクト名',
    'project.siteType': 'サイト種別',
    'project.budget': '予算 USD',
    'project.timeline': '期間',
    'project.brief': '概要',
    'project.estimatedTasks': '推定タスク',
    'project.workerMix': '作業構成',
    'project.workerMixValue': '人 + エージェント',
    'project.checkout': '決済',
    'project.continueCheckout': '決済へ進む',
    'site.landing': 'ランディングページ',
    'site.business': '企業サイト',
    'site.saas': 'SaaS サイト',
    'site.storefront': 'ストアサイト',
    'site.webapp': 'Web アプリ シェル',
    'timeline.7': '7日',
    'timeline.14': '14日',
    'timeline.30': '30日',
    'repo.heading': '既存リポジトリの issue を修正',
    'repo.url': 'GitHub リポジトリ URL',
    'repo.loading': 'Issue を読み込み中...',
    'repo.load': 'Issue を読み込んで採点',
    'repo.repo': 'Repo',
    'repo.openIssues': '未解決 issue',
    'repo.selectedEstimate': '選択分の見積',
    'repo.humanReview': '人間レビュー',
    'repo.openIssue': 'Issue を開く',
    'repo.emptyTitle': '公開 GitHub repo を貼り付け',
    'repo.emptyBody': 'MergeOS は未解決 issue を取り込み、PR を除外し、複雑度を採点し、報酬を見積もり、人またはエージェントの作業を提案します。',
    'repo.orderHeading': 'Issue 注文',
    'repo.selected': '選択済み',
    'repo.issues': 'issue',
    'repo.estimated': '見積',
    'repo.auth': '認証',
    'repo.authValue': '決済時に必要',
    'repo.privateRepo': '非公開 repo',
    'repo.privateRepoValue': 'ログイン後に GitHub 接続',
    'repo.checkoutSelected': '選択した issue を決済',
    'complexity.low': '低',
    'complexity.medium': '中',
    'complexity.high': '高',
    'worker.human': '人間',
    'worker.agent': 'エージェント',
    'worker.hybrid': 'ハイブリッド',
    'reason.openGitHubIssue': '公開 GitHub issue',
    'reason.detailedIssueBody': '詳細な issue 本文',
    'reason.clearReproductionContext': '再現条件が明確',
    'reason.activeDiscussion': '議論が活発',
    'reason.securityOrAuthRisk': 'セキュリティまたは認証リスク',
    'reason.productionRisk': '本番リスク',
    'reason.bugFix': 'バグ修正',
    'reason.backendSurface': 'バックエンド領域',
    'reason.frontendSurface': 'フロントエンド領域',
    'reason.scopeExpansion': 'スコープ拡張',
    'reason.smallEditorialTask': '小さな編集タスク',
    'reason.lowComplexityLabel': '低複雑度ラベル',
  },
  ko: {
    'nav.marketplace': '공개 마켓플레이스',
    'nav.talent': '인재',
    'nav.project': '새 프로젝트',
    'nav.repo': 'Repo issue 수정',
    'nav.login': '로그인',
    'nav.startOrder': '주문 시작',
    'language.label': '언어',
    'language.title': '인터페이스 언어 선택',
    'hero.eyebrow': '성과 기준으로 사람과 에이전트 고용',
    'hero.title': '인재를 찾고, 새 빌드를 펀딩하고, repo issue를 점수화된 유료 작업으로 전환하세요.',
    'hero.body': '게스트는 먼저 인재를 둘러보고 공개 GitHub issue를 분석할 수 있습니다. 계정은 고용, 결제, 비공개 repo 연결, 납품 추적 시에만 필요합니다.',
    'hero.scoreRepo': 'Repo issue 점수화',
    'hero.findTalent': '인재 찾기',
    'signal.authLabel': '인증',
    'signal.authValue': '결제 시에만',
    'signal.repoLabel': 'Repo 가져오기',
    'signal.repoValue': '공개 GitHub issue',
    'signal.scoringLabel': '점수화',
    'signal.scoringValue': '복잡도, 보상, 작업 유형',
    'talent.heading': '인재 마켓플레이스',
    'talent.search': '검색',
    'talent.searchPlaceholder': '프론트엔드, 결제, QA...',
    'talent.skill': '스킬',
    'talent.allSkills': '전체 스킬',
    'talent.rating': '평점',
    'talent.tasks': '작업',
    'talent.taskUnit': '작업',
    'talent.invite': '주문에 초대',
    'skill.frontend': '프론트엔드',
    'skill.backend': '백엔드',
    'skill.design': '디자인',
    'skill.qa': 'QA',
    'skill.ui': 'UI',
    'skill.accessibility': '접근성',
    'skill.payments': '결제',
    'skill.api': 'API',
    'skill.copy': '카피',
    'talentProfile.frontend.role': '프론트엔드 납품 팀',
    'talentProfile.frontend.summary': '반응형 Vue/React 화면, 결제 상태, 대시보드, 접근성 점검을 제공합니다.',
    'talentProfile.backend.role': 'Go API 및 결제 엔지니어',
    'talentProfile.backend.summary': '인증, webhook, 결제 검증, 데이터 마이그레이션, 증명 원장 수정을 처리합니다.',
    'talentProfile.design.role': '휴먼 리뷰 팀',
    'talentProfile.design.summary': '제품 명확성, 카피, 반응형 레이아웃, 대비, 릴리스 준비 상태를 점검합니다.',
    'talentProfile.repo.role': 'Issue 분류 및 PR 에이전트',
    'talentProfile.repo.summary': 'Issue 맥락을 가져오고 범위가 명확한 PR을 만들며 위험한 변경은 휴먼 리뷰와 함께 처리합니다.',
    'project.heading': '프로젝트 주문 미리보기',
    'project.name': '프로젝트 이름',
    'project.siteType': '사이트 유형',
    'project.budget': '예산 USD',
    'project.timeline': '일정',
    'project.brief': '브리프',
    'project.estimatedTasks': '예상 작업',
    'project.workerMix': '작업 구성',
    'project.workerMixValue': '사람 + 에이전트',
    'project.checkout': '결제',
    'project.continueCheckout': '결제로 계속',
    'site.landing': '랜딩 페이지',
    'site.business': '비즈니스 웹사이트',
    'site.saas': 'SaaS 웹사이트',
    'site.storefront': '스토어프론트',
    'site.webapp': '웹 앱 셸',
    'timeline.7': '7일',
    'timeline.14': '14일',
    'timeline.30': '30일',
    'repo.heading': '기존 repo issue 수정',
    'repo.url': 'GitHub repo URL',
    'repo.loading': 'Issue 로딩 중...',
    'repo.load': 'Issue 로드 및 점수화',
    'repo.repo': 'Repo',
    'repo.openIssues': '열린 issue',
    'repo.selectedEstimate': '선택 예상 비용',
    'repo.humanReview': '휴먼 리뷰',
    'repo.openIssue': 'Issue 열기',
    'repo.emptyTitle': '공개 GitHub repo 붙여넣기',
    'repo.emptyBody': 'MergeOS는 열린 issue를 가져오고 PR을 제외한 뒤 복잡도를 점수화하고 보상을 추정하며 사람 또는 에이전트 작업을 제안합니다.',
    'repo.orderHeading': 'Issue 주문',
    'repo.selected': '선택됨',
    'repo.issues': '개 issue',
    'repo.estimated': '예상',
    'repo.auth': '인증',
    'repo.authValue': '결제 시 필요',
    'repo.privateRepo': '비공개 repo',
    'repo.privateRepoValue': '로그인 후 GitHub 연결',
    'repo.checkoutSelected': '선택한 issue 결제',
    'complexity.low': '낮음',
    'complexity.medium': '중간',
    'complexity.high': '높음',
    'worker.human': '사람',
    'worker.agent': '에이전트',
    'worker.hybrid': '하이브리드',
    'reason.openGitHubIssue': '열린 GitHub issue',
    'reason.detailedIssueBody': '상세한 issue 본문',
    'reason.clearReproductionContext': '명확한 재현 맥락',
    'reason.activeDiscussion': '활발한 논의',
    'reason.securityOrAuthRisk': '보안 또는 인증 위험',
    'reason.productionRisk': '프로덕션 위험',
    'reason.bugFix': '버그 수정',
    'reason.backendSurface': '백엔드 영역',
    'reason.frontendSurface': '프론트엔드 영역',
    'reason.scopeExpansion': '범위 확장',
    'reason.smallEditorialTask': '작은 편집 작업',
    'reason.lowComplexityLabel': '낮은 복잡도 라벨',
  },
  vi: {
    'nav.marketplace': 'Chợ công khai',
    'nav.talent': 'Talent',
    'nav.project': 'Project mới',
    'nav.repo': 'Fix issue repo',
    'nav.login': 'Đăng nhập',
    'nav.startOrder': 'Bắt đầu đặt hàng',
    'language.label': 'Ngôn ngữ',
    'language.title': 'Chọn ngôn ngữ giao diện',
    'hero.eyebrow': 'Thuê người và agent theo kết quả',
    'hero.title': 'Tìm talent, đặt build mới, hoặc biến issue trong repo thành task có điểm và bounty.',
    'hero.body': 'Khách có thể xem talent và phân tích issue GitHub public trước. Chỉ cần tạo tài khoản khi muốn thuê, checkout, nối repo private hoặc theo dõi bàn giao.',
    'hero.scoreRepo': 'Chấm điểm issue repo',
    'hero.findTalent': 'Tìm talent',
    'signal.authLabel': 'Cổng auth',
    'signal.authValue': 'Chỉ ở checkout',
    'signal.repoLabel': 'Import repo',
    'signal.repoValue': 'Issue GitHub public',
    'signal.scoringLabel': 'Chấm điểm',
    'signal.scoringValue': 'Độ khó, bounty, loại worker',
    'talent.heading': 'Chợ Talent',
    'talent.search': 'Tìm kiếm',
    'talent.searchPlaceholder': 'frontend, checkout, QA...',
    'talent.skill': 'Kỹ năng',
    'talent.allSkills': 'Tất cả kỹ năng',
    'talent.rating': 'đánh giá',
    'talent.tasks': 'task',
    'talent.taskUnit': 'task',
    'talent.invite': 'Mời vào đơn',
    'skill.frontend': 'Frontend',
    'skill.backend': 'Backend',
    'skill.design': 'Thiết kế',
    'skill.qa': 'QA',
    'skill.ui': 'UI',
    'skill.accessibility': 'Accessibility',
    'skill.payments': 'Thanh toán',
    'skill.api': 'API',
    'skill.copy': 'Copy',
    'talentProfile.frontend.role': 'Nhóm giao diện frontend',
    'talentProfile.frontend.summary': 'Làm giao diện Vue/React responsive, trạng thái checkout, dashboard và kiểm tra accessibility.',
    'talentProfile.backend.role': 'Kỹ sư Go API và thanh toán',
    'talentProfile.backend.summary': 'Xử lý auth, webhook, xác minh thanh toán, migration dữ liệu và sửa proof ledger.',
    'talentProfile.design.role': 'Nhóm review con người',
    'talentProfile.design.summary': 'Kiểm tra độ rõ sản phẩm, copy, layout responsive, tương phản và sẵn sàng release.',
    'talentProfile.repo.role': 'Agent triage issue và PR',
    'talentProfile.repo.summary': 'Import ngữ cảnh issue, mở PR có scope rõ, và ghép review con người cho thay đổi rủi ro.',
    'project.heading': 'Preview đơn project',
    'project.name': 'Tên project',
    'project.siteType': 'Loại site',
    'project.budget': 'Ngân sách USD',
    'project.timeline': 'Timeline',
    'project.brief': 'Mô tả',
    'project.estimatedTasks': 'Task dự kiến',
    'project.workerMix': 'Worker mix',
    'project.workerMixValue': 'Người + agent',
    'project.checkout': 'Checkout',
    'project.continueCheckout': 'Tiếp tục checkout',
    'site.landing': 'Landing page',
    'site.business': 'Website doanh nghiệp',
    'site.saas': 'Website SaaS',
    'site.storefront': 'Cửa hàng',
    'site.webapp': 'Web app shell',
    'timeline.7': '7 ngày',
    'timeline.14': '14 ngày',
    'timeline.30': '30 ngày',
    'repo.heading': 'Fix issue trong repo có sẵn',
    'repo.url': 'URL repo GitHub',
    'repo.loading': 'Đang load issue...',
    'repo.load': 'Load và chấm điểm issue',
    'repo.repo': 'Repo',
    'repo.openIssues': 'Issue mở',
    'repo.selectedEstimate': 'Ước tính đã chọn',
    'repo.humanReview': 'human-review',
    'repo.openIssue': 'Mở issue',
    'repo.emptyTitle': 'Dán repo GitHub public',
    'repo.emptyBody': 'MergeOS sẽ import issue mở, bỏ qua pull request, chấm độ khó, ước tính bounty và gợi ý người hoặc agent làm.',
    'repo.orderHeading': 'Đơn issue',
    'repo.selected': 'Đã chọn',
    'repo.issues': 'issue',
    'repo.estimated': 'Ước tính',
    'repo.auth': 'Auth',
    'repo.authValue': 'Cần ở checkout',
    'repo.privateRepo': 'Repo private',
    'repo.privateRepoValue': 'Nối GitHub sau đăng nhập',
    'repo.checkoutSelected': 'Checkout issue đã chọn',
    'complexity.low': 'thấp',
    'complexity.medium': 'vừa',
    'complexity.high': 'cao',
    'worker.human': 'người',
    'worker.agent': 'agent',
    'worker.hybrid': 'lai',
    'reason.openGitHubIssue': 'issue GitHub đang mở',
    'reason.detailedIssueBody': 'mô tả issue chi tiết',
    'reason.clearReproductionContext': 'ngữ cảnh tái hiện rõ',
    'reason.activeDiscussion': 'thảo luận sôi nổi',
    'reason.securityOrAuthRisk': 'rủi ro bảo mật hoặc auth',
    'reason.productionRisk': 'rủi ro production',
    'reason.bugFix': 'sửa bug',
    'reason.backendSurface': 'phạm vi backend',
    'reason.frontendSurface': 'phạm vi frontend',
    'reason.scopeExpansion': 'mở rộng scope',
    'reason.smallEditorialTask': 'task biên tập nhỏ',
    'reason.lowComplexityLabel': 'nhãn độ khó thấp',
  },
};

const issueReasonKeys = {
  'open GitHub issue': 'reason.openGitHubIssue',
  'detailed issue body': 'reason.detailedIssueBody',
  'clear reproduction context': 'reason.clearReproductionContext',
  'active discussion': 'reason.activeDiscussion',
  'security or auth risk': 'reason.securityOrAuthRisk',
  'production risk': 'reason.productionRisk',
  'bug fix': 'reason.bugFix',
  'backend surface': 'reason.backendSurface',
  'frontend surface': 'reason.frontendSurface',
  'scope expansion': 'reason.scopeExpansion',
  'small editorial task': 'reason.smallEditorialTask',
  'low complexity label': 'reason.lowComplexityLabel',
};

const talentProfiles = [
  {
    id: 'frontend-ana',
    initials: 'AA',
    name: 'An Agent Studio',
    roleKey: 'talentProfile.frontend.role',
    summaryKey: 'talentProfile.frontend.summary',
    skills: ['frontend', 'ui', 'accessibility'],
    rating: '4.9',
    completed: 128,
    rate_cents: 18000,
  },
  {
    id: 'backend-ledger',
    initials: 'BL',
    name: 'Backend Ledger Crew',
    roleKey: 'talentProfile.backend.role',
    summaryKey: 'talentProfile.backend.summary',
    skills: ['backend', 'payments', 'api'],
    rating: '4.8',
    completed: 96,
    rate_cents: 22000,
  },
  {
    id: 'design-qa',
    initials: 'DQ',
    name: 'Design QA Review',
    roleKey: 'talentProfile.design.role',
    summaryKey: 'talentProfile.design.summary',
    skills: ['design', 'qa', 'copy'],
    rating: '4.7',
    completed: 74,
    rate_cents: 14000,
  },
  {
    id: 'repo-fix',
    initials: 'RF',
    name: 'Repo Fix Agents',
    roleKey: 'talentProfile.repo.role',
    summaryKey: 'talentProfile.repo.summary',
    skills: ['backend', 'frontend', 'qa'],
    rating: '4.8',
    completed: 151,
    rate_cents: 20000,
  },
];

const authForm = reactive({
  name: 'Thanh Truc Client',
  company_name: 'Thanh Truc Solutions',
  email: 'client@mergeos.local',
  password: 'mergeos123',
});

const projectForm = reactive({
  client_name: 'Thanh Truc Client',
  company_name: 'Thanh Truc Solutions',
  client_email: 'client@mergeos.local',
  phone: '+84 900 000 000',
  title: 'Elementor-style company website',
  site_type: 'Business website',
  package_tier: 'Growth',
  timeline: '14 days',
  brief: 'Build a polished responsive website with services, portfolio, lead form, payment-ready checkout, and customer dashboard preview.',
  budget_cents: 240000,
  payment_method: 'paypal',
  payment_reference: 'LOCAL-PAID',
});

const workerForm = reactive({
  worker_kind: 'agent',
  worker_id: 'agent:mergeos-web-001',
  agent_type: 'frontend-agent',
});

const repoForm = reactive({
  repo_url: 'https://github.com/vuejs/core',
});

const isAdmin = computed(() => user.value?.role === 'admin');
const showPublicShell = computed(() => !user.value && !authVisible.value);
const activeLocale = computed(() => localeByLanguage[language.value] || localeByLanguage.en);
const filteredTalents = computed(() => {
  const query = talentQuery.value.trim().toLowerCase();
  return talentProfiles.filter((talent) => {
    const matchesSkill = talentSkill.value === 'all' || talent.skills.includes(talentSkill.value);
    const haystack = `${talent.name} ${t(talent.roleKey)} ${t(talent.summaryKey)} ${talent.skills.map((skill) => t(`skill.${skill}`)).join(' ')}`.toLowerCase();
    return matchesSkill && (!query || haystack.includes(query));
  });
});
const currentProject = computed(() => {
  if (selectedProjectId.value) {
    return projects.value.find((project) => project.id === selectedProjectId.value) || projects.value[projects.value.length - 1] || null;
  }
  return projects.value[projects.value.length - 1] || null;
});
const currentTasks = computed(() => {
  if (!currentProject.value) return [];
  return tasks.value.filter((task) => task.project_id === currentProject.value.id);
});
const acceptedTasks = computed(() => tasks.value.filter((task) => task.status === 'accepted'));
const openTasks = computed(() => tasks.value.filter((task) => task.status !== 'accepted'));
const recentLedger = computed(() => ledger.value.slice(-8).reverse());
const adminCurrentProject = computed(() => {
  if (selectedAdminProjectId.value) {
    return projects.value.find((project) => project.id === selectedAdminProjectId.value) || projects.value[0] || null;
  }
  return projects.value[0] || null;
});
const adminProjectTasks = computed(() => {
  if (!adminCurrentProject.value) return [];
  return tasks.value.filter((task) => task.project_id === adminCurrentProject.value.id);
});
const adminOpenTasks = computed(() => tasks.value.filter((task) => task.status !== 'accepted'));
const adminLedgerRows = computed(() => ledger.value.slice().reverse());
const totalBudget = computed(() => projects.value.reduce((sum, project) => sum + project.budget_cents, 0));
const totalPool = computed(() => projects.value.reduce((sum, project) => sum + project.work_pool_cents, 0));
const selectedRepoIssueTotal = computed(() => {
  if (!repoImportResult.value) return 0;
  return repoImportResult.value.issues
    .filter((issue) => selectedRepoIssueNumbers.value.includes(issue.number))
    .reduce((sum, issue) => sum + issue.estimated_cents, 0);
});
const tokenSymbol = computed(() => runtimeConfig.value?.token_symbol || 'MERGE');
const statusLabel = computed(() => {
  if (currentProject.value) return `${currentProject.value.payment_provider} verified`;
  return runtimeConfig.value?.repo_provider || 'ready';
});
const paymentReferencePlaceholder = computed(() => {
  if (projectForm.payment_method === 'paypal') return 'PayPal order id';
  if (runtimeConfig.value?.crypto_asset === 'erc20') return 'EVM tx hash for stablecoin transfer';
  return 'EVM tx hash';
});
const budgetUsd = computed({
  get: () => Math.round(projectForm.budget_cents / 100),
  set: (value) => {
    projectForm.budget_cents = Math.max(100, Number(value || 100)) * 100;
  },
});

function t(key) {
  return translations[language.value]?.[key] || translations.en[key] || key;
}

function issueComplexity(value) {
  return t(`complexity.${value || 'medium'}`);
}

function workerKindLabel(value) {
  return t(`worker.${value || 'hybrid'}`);
}

function issueReasonText(reasons = []) {
  return reasons.map((reason) => t(issueReasonKeys[reason] || reason)).join(', ');
}

function money(cents) {
  return new Intl.NumberFormat(activeLocale.value, {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 0,
  }).format((cents || 0) / 100);
}

function fileSize(bytes) {
  if (!bytes) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB'];
  const power = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
  const value = bytes / (1024 ** power);
  return `${value.toFixed(value >= 10 || power === 0 ? 0 : 1)} ${units[power]}`;
}

function formatDate(value) {
  if (!value) return 'n/a';
  return new Intl.DateTimeFormat(activeLocale.value, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value));
}

function sslStatusLabel(status) {
  const labels = {
    ok: 'Valid',
    warning: 'Expiring',
    expired: 'Expired',
    error: 'Issue',
    pending: 'Pending',
  };
  return labels[status] || status || 'Pending';
}

function sslDaysText(review) {
  if (!review?.not_after) return 'waiting';
  const days = Number(review.days_remaining || 0);
  if (days < 0) return `${Math.abs(days)} days expired`;
  if (days === 0) return 'expires today';
  return `${days} days left`;
}

function projectTitle(projectId) {
  return projects.value.find((project) => project.id === projectId)?.title || projectId;
}

function attachmentCountForProject(projectId) {
  return attachments.value.filter((attachment) => attachment.project_id === projectId).length;
}

function shortHash(hash) {
  return `${hash.slice(0, 8)}...${hash.slice(-6)}`;
}

async function api(path, options = {}) {
  const isFormData = typeof FormData !== 'undefined' && options.body instanceof FormData;
  const response = await fetch(path, {
    ...options,
    headers: {
      ...(isFormData ? {} : { 'Content-Type': 'application/json' }),
      ...(token.value ? { Authorization: `Bearer ${token.value}` } : {}),
      ...(options.headers || {}),
    },
  });
  const payload = await response.json();
  if (!response.ok) {
    if (response.status === 401) clearSession();
    throw new Error(payload.error || 'Request failed');
  }
  return payload;
}

function openAuth(mode = 'register') {
  authMode.value = mode;
  authVisible.value = true;
  errorMessage.value = '';
}

function backToPublic() {
  authVisible.value = false;
  errorMessage.value = '';
}

function startHireTalent(talent) {
  projectForm.brief = `Invite ${talent.name} for ${t(talent.roleKey)}. ${t(talent.summaryKey)}`;
  openAuth('register');
}

async function importRepoIssues() {
  repoImportBusy.value = true;
  repoImportError.value = '';
  try {
    const result = await api('/api/public/repo/issues', {
      method: 'POST',
      body: JSON.stringify({ repo_url: repoForm.repo_url }),
    });
    repoImportResult.value = result;
    selectedRepoIssueNumbers.value = result.issues.slice(0, 3).map((issue) => issue.number);
  } catch (error) {
    repoImportError.value = error.message;
  } finally {
    repoImportBusy.value = false;
  }
}

function toggleRepoIssue(issue) {
  if (selectedRepoIssueNumbers.value.includes(issue.number)) {
    selectedRepoIssueNumbers.value = selectedRepoIssueNumbers.value.filter((number) => number !== issue.number);
    return;
  }
  selectedRepoIssueNumbers.value = selectedRepoIssueNumbers.value.concat(issue.number);
}

async function loadConfig() {
  runtimeConfig.value = await api('/api/config');
  if (runtimeConfig.value.dev_payment_enabled && !projectForm.payment_reference) {
    projectForm.payment_reference = runtimeConfig.value.dev_payment_code;
  }
}

async function submitAuth() {
  authBusy.value = true;
  errorMessage.value = '';
  try {
    const path = authMode.value === 'register' ? '/api/auth/register' : '/api/auth/login';
    const body = authMode.value === 'register'
      ? authForm
      : { email: authForm.email, password: authForm.password };
    const auth = await api(path, { method: 'POST', body: JSON.stringify(body) });
    setSession(auth);
    if (!isAdmin.value) syncProjectContact();
    await refreshProtected();
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    authBusy.value = false;
  }
}

function setSession(auth) {
  token.value = auth.token;
  user.value = auth.user;
  if (hasWindow) localStorage.setItem('mergeos_token', auth.token);
}

function clearSession() {
  token.value = '';
  user.value = null;
  authVisible.value = false;
  projects.value = [];
  tasks.value = [];
  ledger.value = [];
  notifications.value = [];
  adminSummary.value = null;
  adminUsers.value = [];
  attachments.value = [];
  sslReviews.value = [];
  selectedAdminProjectId.value = '';
  adminSelectedTask.value = null;
  uploadedAttachments.value = [];
  if (hasWindow) localStorage.removeItem('mergeos_token');
}

function syncProjectContact() {
  if (!user.value) return;
  projectForm.client_name = user.value.name;
  projectForm.company_name = user.value.company_name;
  projectForm.client_email = user.value.email;
}

async function restoreSession() {
  if (!token.value) return;
  try {
    user.value = await api('/api/auth/me');
    if (!isAdmin.value) syncProjectContact();
    await refreshProtected();
  } catch {
    clearSession();
  }
}

async function refreshAll() {
  await loadConfig();
  if (user.value) {
    await refreshProtected();
  }
}

async function refreshProtected() {
  if (isAdmin.value) {
    const [summary, userRows, projectRows, taskRows, ledgerRows, noteRows, attachmentRows, sslRows] = await Promise.all([
      api('/api/admin/summary'),
      api('/api/admin/users'),
      api('/api/admin/projects'),
      api('/api/admin/tasks'),
      api('/api/admin/ledger'),
      api('/api/admin/notifications'),
      api('/api/admin/attachments'),
      api('/api/admin/ssl'),
    ]);
    adminSummary.value = summary;
    adminUsers.value = userRows;
    projects.value = projectRows;
    tasks.value = taskRows;
    ledger.value = ledgerRows;
    notifications.value = noteRows.slice().reverse();
    attachments.value = attachmentRows;
    sslReviews.value = sslRows.length ? sslRows : (summary.ssl_reviews || []);
    if (!selectedAdminProjectId.value && projectRows.length) {
      selectedAdminProjectId.value = projectRows[0].id;
    }
    reconcileAdminSelection();
    return;
  }

  const [projectRows, taskRows, ledgerRows, noteRows] = await Promise.all([
    api('/api/projects'),
    api('/api/tasks'),
    api('/api/ledger'),
    api('/api/notifications'),
  ]);
  projects.value = projectRows;
  tasks.value = taskRows;
  ledger.value = ledgerRows;
  notifications.value = noteRows.slice().reverse();
  if (!selectedProjectId.value && projectRows.length) {
    selectedProjectId.value = projectRows[projectRows.length - 1].id;
  }
  reconcileSelection();
}

async function reviewSSL() {
  sslReviewBusy.value = true;
  errorMessage.value = '';
  try {
    const rows = await api('/api/admin/ssl/review', { method: 'POST' });
    sslReviews.value = rows;
    if (adminSummary.value) {
      adminSummary.value = { ...adminSummary.value, ssl_reviews: rows };
    }
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    sslReviewBusy.value = false;
  }
}

async function preparePayPalOrder() {
  preparingPayPal.value = true;
  errorMessage.value = '';
  try {
    const order = await api('/api/payments/paypal/orders', {
      method: 'POST',
      body: JSON.stringify({
        amount_cents: projectForm.budget_cents,
        description: projectForm.title,
        return_url: hasWindow ? window.location.href : 'http://127.0.0.1:5173',
        cancel_url: hasWindow ? window.location.href : 'http://127.0.0.1:5173',
      }),
    });
    paypalOrder.value = order;
    projectForm.payment_reference = order.order_id;
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    preparingPayPal.value = false;
  }
}

async function createProject() {
  creating.value = true;
  errorMessage.value = '';
  try {
    const project = await api('/api/projects', {
      method: 'POST',
      body: JSON.stringify({
        ...projectForm,
        attachment_ids: uploadedAttachments.value.map((attachment) => attachment.id),
      }),
    });
    selectedProjectId.value = project.id;
    uploadedAttachments.value = [];
    await refreshProtected();
    portalTab.value = 'workspace';
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    creating.value = false;
  }
}

async function uploadProjectFiles(event) {
  const input = event.target;
  const files = Array.from(input.files || []);
  if (!files.length) return;

  uploadBusy.value = true;
  errorMessage.value = '';
  try {
    const formData = new FormData();
    files.forEach((file) => formData.append('files', file));
    const attachments = await api('/api/uploads', {
      method: 'POST',
      body: formData,
    });
    uploadedAttachments.value = uploadedAttachments.value.concat(attachments);
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    uploadBusy.value = false;
    input.value = '';
  }
}

function removeUploadedAttachment(id) {
  uploadedAttachments.value = uploadedAttachments.value.filter((attachment) => attachment.id !== id);
}

async function openAttachment(attachment) {
  errorMessage.value = '';
  try {
    const response = await fetch(attachment.url, {
      headers: {
        ...(token.value ? { Authorization: `Bearer ${token.value}` } : {}),
      },
    });
    if (!response.ok) {
      throw new Error('Could not open file');
    }
    const blob = await response.blob();
    const blobURL = URL.createObjectURL(blob);
    if (hasWindow) {
      const opened = window.open(blobURL, '_blank', 'noopener,noreferrer');
      if (!opened) {
        const link = document.createElement('a');
        link.href = blobURL;
        link.download = attachment.original_name || 'attachment';
        link.click();
      }
      window.setTimeout(() => URL.revokeObjectURL(blobURL), 60000);
    }
  } catch (error) {
    errorMessage.value = error.message;
  }
}

function selectAdminProject(project) {
  selectedAdminProjectId.value = project.id;
  reconcileAdminSelection();
}

function selectAdminTask(task) {
  adminSelectedTask.value = task;
  workerForm.worker_kind = task.required_worker_kind;
  workerForm.worker_id = task.required_worker_kind === 'human' ? 'github:admin-reviewer' : 'agent:mergeos-admin-001';
  workerForm.agent_type = task.required_worker_kind === 'human' ? '' : (task.suggested_agent_type || 'custom-agent');
  errorMessage.value = '';
}

async function acceptAdminSelectedTask() {
  if (!adminSelectedTask.value) return;
  accepting.value = true;
  errorMessage.value = '';
  try {
    await api(`/api/tasks/${adminSelectedTask.value.id}/accept`, {
      method: 'POST',
      body: JSON.stringify(workerForm),
    });
    await refreshProtected();
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    accepting.value = false;
  }
}

function selectProject(project) {
  selectedProjectId.value = project.id;
  reconcileSelection();
}

function selectTask(task) {
  selectedTask.value = task;
  workerForm.worker_kind = task.required_worker_kind;
  workerForm.worker_id = task.required_worker_kind === 'human' ? 'github:client-reviewer' : 'agent:mergeos-web-001';
  workerForm.agent_type = task.required_worker_kind === 'human' ? '' : (task.suggested_agent_type || 'custom-agent');
  errorMessage.value = '';
}

async function acceptSelectedTask() {
  if (!selectedTask.value) return;
  accepting.value = true;
  errorMessage.value = '';
  try {
    await api(`/api/tasks/${selectedTask.value.id}/accept`, {
      method: 'POST',
      body: JSON.stringify(workerForm),
    });
    await refreshProtected();
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    accepting.value = false;
  }
}

async function logout() {
  try {
    await api('/api/auth/logout', { method: 'POST', body: JSON.stringify({}) });
  } finally {
    clearSession();
  }
}

function reconcileSelection() {
  if (selectedTask.value) {
    const fresh = currentTasks.value.find((task) => task.id === selectedTask.value.id);
    selectedTask.value = fresh || currentTasks.value[0] || null;
    return;
  }
  if (currentTasks.value.length) {
    selectTask(currentTasks.value[0]);
  }
}

function reconcileAdminSelection() {
  if (adminSelectedTask.value) {
    const fresh = tasks.value.find((task) => task.id === adminSelectedTask.value.id);
    adminSelectedTask.value = fresh || adminProjectTasks.value[0] || adminOpenTasks.value[0] || null;
    return;
  }
  if (adminProjectTasks.value.length) {
    selectAdminTask(adminProjectTasks.value[0]);
    return;
  }
  if (adminOpenTasks.value.length) {
    selectAdminTask(adminOpenTasks.value[0]);
  }
}

function normalizeLanguage(value) {
  const normalized = String(value || '').trim().toLowerCase();
  if (normalized.startsWith('zh')) return 'zh';
  if (normalized.startsWith('ja')) return 'ja';
  if (normalized.startsWith('ko')) return 'ko';
  if (normalized.startsWith('vi')) return 'vi';
  return 'en';
}

function loadPreferredLanguage() {
  if (!hasWindow) return;
  const saved = localStorage.getItem('mergeos_language');
  language.value = saved ? normalizeLanguage(saved) : normalizeLanguage(navigator.language);
}

watch(currentTasks, reconcileSelection);
watch(adminProjectTasks, reconcileAdminSelection);
watch(language, (value) => {
  if (!hasWindow) return;
  const normalized = normalizeLanguage(value);
  if (normalized !== value) {
    language.value = normalized;
    return;
  }
  localStorage.setItem('mergeos_language', normalized);
  document.documentElement.lang = localeByLanguage[normalized]?.split('-')[0] || 'en';
});

onMounted(async () => {
  loadPreferredLanguage();
  await loadConfig();
  await restoreSession();
});
</script>
