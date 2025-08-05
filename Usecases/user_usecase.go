package usecases

import (
	domain "blog-api/Domain"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type UserUsecase struct {
	userRepository   domain.IUserRepository
	refreshTokenRepo domain.IRefreshTokenRepository
	JWTService       domain.IJWTService
	passwordService  domain.IPasswordService
	contextTimeout   time.Duration
}

func NewUserUseCase(userRepo domain.IUserRepository, refreshRepo domain.IRefreshTokenRepository,
	jwt domain.IJWTService, passwordService domain.IPasswordService, timeout time.Duration) domain.IUserUsecase {
	return &UserUsecase{
		userRepository:   userRepo,
		JWTService:       jwt,
		refreshTokenRepo: refreshRepo,
		passwordService:  passwordService,
		contextTimeout:   timeout,
	}
}

func (uc *UserUsecase) Register(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if user.Email == "" || user.Password == "" || user.Username == "" {
		return domain.ErrInvalidInput
	}

	if err := uc.passwordService.ValidateStrength(user.Password); err != nil {
		return err
	}

	//check if email or username already exists
	emailExists, err := uc.userRepository.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if emailExists {
		return domain.ErrEmailTaken
	}

	usernameExists, err := uc.userRepository.ExistsByUsername(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("failed to check username: %w", err)
	}
	if usernameExists {
		return domain.ErrUsernameTaken
	}

	//hash password
	hashedPassword, err := uc.passwordService.Hash(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	//save user
	_, err = uc.userRepository.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (uc *UserUsecase) Login(ctx context.Context, user *domain.User) (*domain.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if user.Email == "" || user.Password == "" {
		return nil, domain.ErrInvalidInput
	}

	dbUser, err := uc.userRepository.GetByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUnauthorized
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if dbUser == nil {
		return nil, domain.ErrUnauthorized
	}
	if err := uc.passwordService.Compare(dbUser.Password, user.Password); err != nil {
		return nil, errors.New("incorrect email or password")
	}

	//generate access token
	accessToken, err := uc.JWTService.CreateAccessToken(dbUser)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	err = uc.refreshTokenRepo.DeleteAllTokensForUser(ctx, dbUser.ID)
	if err != nil {
		log.Printf("warning: failed to delete existing refresh tokens for user %s: %v", dbUser.ID, err)
	} else {
		log.Printf("successfully deleted all existing refresh tokens for user %s", dbUser.ID)
	}

	//generate refresh token
	refreshToken, refreshTokenPayload, err := uc.JWTService.CreateRefreshToken(dbUser)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	//store refresh token
	rt := &domain.RefreshToken{
		ID:        refreshTokenPayload.TokenID,
		Token:     refreshToken,
		UserID:    dbUser.ID,
		CreatedAt: refreshTokenPayload.IssuedAt,
		UpdatedAt: refreshTokenPayload.IssuedAt,
		ExpiresAt: refreshTokenPayload.ExpiresAt,
	}

	if err := uc.refreshTokenRepo.StoreToken(ctx, rt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}


func (uc *UserUsecase) Logout(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if userID == "" {
		return domain.ErrInvalidInput
	}
	
	_, err := uc.userRepository.GetByID(ctx, userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	//delete all refresh tokens for the user
	err = uc.refreshTokenRepo.DeleteAllTokensForUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to logout user: %w", err)
	}

	return nil
}
func (uc *UserUsecase) Promote (ctx context.Context, username string) error{
	user,err := uc.userRepository.GetByUsername(ctx,username)
	if err != nil{
		return err
	}
	if user.Role == "admin"{
		return errors.New("user is already an admin")
	}
    return uc.userRepository.Promote(ctx,user)
}
