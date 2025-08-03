package domain

import "errors"

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrEmailTaken    = errors.New("email is already registered")
	ErrUsernameTaken = errors.New("username is already taken")
	ErrUserNotFound  = errors.New("user not found")
	ErrUnauthorized  = errors.New("incorrect email or password")
	ErrInternal      = errors.New("internal error")
)
