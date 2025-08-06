package infrastructure

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	domain "blog-api/Domain"

	"gopkg.in/gomail.v2"
)

type SMTPEmailService struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

func NewSMTPEmailService(from, host string, port int, username, password string) domain.IEmailService {
	return &SMTPEmailService{
		From:     from,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (s *SMTPEmailService) SendPasswordResetEmail(ctx context.Context, toEmail string, resetToken string) error {

	tmpl, err := template.ParseFiles("Infrastructure/templates/reset_email.html")
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, struct {
		ResetLink string
	}{
		ResetLink: fmt.Sprintf("http://localhost:3000/reset-password?token=%s", resetToken),
	})
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Password Reset Request")
	m.SetBody("text/html", body.String())

	
	d := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)
	d.SSL = true 

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
