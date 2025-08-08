package domain

import "context"

type IEmailService interface {
	SendPasswordResetEmail(ctx context.Context, toEmail string, resetToken string) error
	SendVerificationEmail(ctx context.Context, toEmail string, otp string) error
}
