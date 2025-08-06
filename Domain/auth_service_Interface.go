package domain

import "time"

type IJWTService interface {
	// GenerateToken(user User) (string, error)
	CreateAccessToken(user *User) (accessToken string, err error)
	CreateRefreshToken(user *User) (refreshToken string, payload *RefreshTokenPayload, err error)
	ValidateAccessToken(tokenString string) (*AccessTokenPayload, error)
	ValidateRefreshToken(tokenString string) (*RefreshTokenPayload, error)
	// ExtractIDFromToken(requestToken string, secret string) (string, error)
}

type IPasswordService interface {
	Hash(password string) (string, error)
	Compare(hashed, plain string) error
	ValidateStrength(password string) error
	GenerateRandomToken() (string,error) 
}
type RefreshTokenPayload struct {
	TokenID   string
	UserID    string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type AccessTokenPayload struct {
	UserID string
	Email  string
	Role   Role
}