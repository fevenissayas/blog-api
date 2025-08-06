package infrastructure

import (
	domain "blog-api/Domain"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	passwordvalidator "github.com/wagslane/go-password-validator"
)

const minEntropyBits = 60 //for password validation

type passwordService struct{}

func NewPasswordService() domain.IPasswordService {
	return &passwordService{}
}

func (p *passwordService) Hash(password string) (string, error) {
	const hashCost = bcrypt.DefaultCost
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(passwordBytes), nil
}

func (p *passwordService) Compare(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

func (p *passwordService) ValidateStrength(password string) error {
	err := passwordvalidator.Validate(password, minEntropyBits)

	if err != nil {
		return errors.New("password too weak")
	}

	return nil
}

func (p *passwordService) GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)

	return token, nil
}
