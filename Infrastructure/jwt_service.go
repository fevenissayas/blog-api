package infrastructure

import (
	domain "blog-api/Domain"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	SecretKey string
}

func NewJWTService() domain.IJWTService {
	secretKey := Env.JWT_SECRET
	if secretKey == "" {
		log.Fatal("JWT_SECRET not set in environment")
	}
	return &JWTService{SecretKey: secretKey}
}

type JwtCustomAccessClaims struct {
	Email string
	Role  domain.Role
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	TokenID string
	UserID  string
	jwt.RegisteredClaims
}

func (j *JWTService) CreateAccessToken(user *domain.User) (string, error) {

	claims := JwtCustomAccessClaims{
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "blog-api",
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.SecretKey))

	if err != nil {
		return "", fmt.Errorf("error creating token: %w", err)
	}
	return tokenString, nil
}
func (j *JWTService) CreateRefreshToken(user *domain.User) (string, *domain.RefreshTokenPayload, error) {
	tokenID := uuid.NewString()
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(7 * 24 * time.Hour)

	claims := JwtCustomRefreshClaims{
		TokenID: tokenID,
		UserID:  user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    "blog-api",
			NotBefore: jwt.NewNumericDate(issuedAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", nil, fmt.Errorf("error creating token: %w", err)
	}

	payload := &domain.RefreshTokenPayload{
		TokenID:   tokenID,
		UserID:    user.ID,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}

	return tokenString, payload, nil
}

func (j *JWTService) ValidateAccessToken(tokenString string) (*domain.AccessTokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomAccessClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid or expired token: %w", err)
	}

	claims, ok := token.Claims.(*JwtCustomAccessClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims: %w", err)
	}

	return &domain.AccessTokenPayload{
		Email: claims.Email,
		Role:  claims.Role,
	}, nil
}
func (j *JWTService) ValidateRefreshToken(tokenString string) (*domain.RefreshTokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomRefreshClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid or expired token: %w", err)
	}

	claims, ok := token.Claims.(*JwtCustomRefreshClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims: %w", err)
	}

	if claims.Issuer != "blog-api" {
		return nil, fmt.Errorf("invalid issuer")
	}

	if claims.IssuedAt == nil || claims.ExpiresAt == nil {
		return nil, fmt.Errorf("token missing time claims")
	}

	return &domain.RefreshTokenPayload{
		TokenID:   claims.TokenID,
		UserID:    claims.UserID,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}
