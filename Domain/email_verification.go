package domain

import (
	"context"
	"time"
)

type EmailVerification struct {
	ID        string
	Email     string
	OTP       string
	ExpiresAt time.Time
	CreatedAt time.Time
	Used      bool
}

type IEmailVerificationRepository interface {
	Store(ctx context.Context, verification *EmailVerification) error
	GetByEmail(ctx context.Context, email string) (*EmailVerification, error)
	GetByOTP(ctx context.Context, otp, email string) (*EmailVerification, error)
	MarkUsed(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
}

type VerifyEmailInput struct {
	Email string
	OTP   string
}
