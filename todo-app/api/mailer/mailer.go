package mailer

import (
	"net/smtp"
	"os"
)

type Mailer interface {
	Send(to, subject, body string) error
}

type SMTPMailer struct {
	host string
	port string
	user string
	pass string
	from string
}

func NewMailerFromEnv() Mailer {
	return &SMTPMailer{
		host: getEnv("SMTP_HOST", "mailpit"),
		port: getEnv("SMTP_PORT", "1025"),
		user: getEnv("SMTP_USER", ""),
		pass: getEnv("SMTP_PASSWORD", ""),
		from: getEnv("SMTP_FROM", "noreply@example.com"),
	}
}

func (m *SMTPMailer) Send(to, subject, body string) error {
	addr := m.host + ":" + m.port
	var auth smtp.Auth
	if m.user != "" {
		auth = smtp.PlainAuth("", m.user, m.pass, m.host)
	}
	msg := []byte(
		"From: " + m.from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" + body,
	)
	return smtp.SendMail(addr, auth, m.from, []string{to}, msg)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
