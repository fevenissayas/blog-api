package domain

import (
	"context"
	"time"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID             string
	Password       string
	Email          string
	Username       string
	Role           Role
	IsVerified     bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Bio            string
	ProfilePicture string
	ContactInfo    string
}

type IUserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}
type IUserUsecase interface {
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, user *User) (*TokenResponse, error)
}
