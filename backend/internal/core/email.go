package core

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailSender struct {
	cfg Config
}

func NewEmailSender(cfg Config) *EmailSender {
	return &EmailSender{cfg: cfg}
}

func (s *EmailSender) Send(to, subject, body string) string {
	to = strings.TrimSpace(to)
	if to == "" {
		return "skipped:no-recipient"
	}
	if !s.cfg.SMTPReady() {
		return "logged:mail-not-configured"
	}

	addr := s.cfg.SMTPHost + ":" + s.cfg.SMTPPort
	auth := smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)
	message := strings.Join([]string{
		"From: " + s.cfg.SMTPFrom,
		"To: " + to,
		"Subject: " + subject,
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")
	if err := smtp.SendMail(addr, auth, s.cfg.SMTPFrom, []string{to}, []byte(message)); err != nil {
		return fmt.Sprintf("error:%s", err.Error())
	}
	return "sent"
}
