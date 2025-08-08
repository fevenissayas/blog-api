package infrastructure

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
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

func (s *SMTPEmailService) SendVerificationEmail(ctx context.Context, toEmail string, otp string) error {
	// Create verification email template inline
	verificationTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .otp-code { font-size: 32px; font-weight: bold; text-align: center; 
                    background-color: #fff; padding: 20px; margin: 20px 0;
                    border: 2px dashed #4CAF50; letter-spacing: 5px; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Email Verification Required</h1>
        </div>
        <div class="content">
            <h2>Welcome to Blog API!</h2>
            <p>Thank you for registering with us. To complete your registration, please verify your email address using the OTP code below:</p>
            
            <div class="otp-code">{{.OTP}}</div>
            
            <p><strong>Important:</strong></p>
            <ul>
                <li>This OTP will expire in 15 minutes</li>
                <li>Do not share this code with anyone</li>
                <li>If you didn't create an account, please ignore this email</li>
            </ul>
            
            <p>Once verified, you'll be able to login and start using our blog platform.</p>
        </div>
        <div class="footer">
            <p>This is an automated email. Please do not reply to this message.</p>
        </div>
    </div>
</body>
</html>`

	tmpl, err := template.New("verification").Parse(verificationTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse verification template: %w", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, struct {
		OTP string
	}{
		OTP: otp,
	})
	if err != nil {
		return fmt.Errorf("failed to execute verification template: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Verify Your Email Address - Blog API")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)
	d.SSL = true

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	log.Printf("Verification email sent successfully to %s", toEmail)
	return nil
}

func ParsePort(portStr string, defaultPort int) int {
	if portStr == "" {
		return defaultPort
	}
	p, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Warning: invalid port '%s', defaulting to %d\n", portStr, defaultPort)
		return defaultPort
	}
	return p
}
