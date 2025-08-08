package usecases

import (
	domain "blog-api/Domain"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"
)

type UserUsecase struct {
	userRepository         domain.IUserRepository
	emailVerificationRepo  domain.IEmailVerificationRepository
	refreshTokenRepo       domain.IRefreshTokenRepository
	passwordResetTokenRepo domain.IPasswordResetTokenRepository
	JWTService             domain.IJWTService
	emailService           domain.IEmailService
	passwordService        domain.IPasswordService
	contextTimeout         time.Duration
}

func NewUserUseCase(userRepo domain.IUserRepository, emailVerificationRepo domain.IEmailVerificationRepository,
	refreshRepo domain.IRefreshTokenRepository, resetTokenRepo domain.IPasswordResetTokenRepository,
	jwt domain.IJWTService, passwordService domain.IPasswordService, emailService domain.IEmailService,
	timeout time.Duration) domain.IUserUsecase {
	return &UserUsecase{
		userRepository:         userRepo,
		emailVerificationRepo:  emailVerificationRepo,
		refreshTokenRepo:       refreshRepo,
		passwordResetTokenRepo: resetTokenRepo,
		JWTService:             jwt,
		passwordService:        passwordService,
		emailService:           emailService,
		contextTimeout:         timeout,
	}
}

func (uc *UserUsecase) generateOTP() (string, error) {
	// Generate a 6-digit OTP
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
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

	// Check if email or username already exists
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

	// Hash password
	hashedPassword, err := uc.passwordService.Hash(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	// Generate OTP
	otp, err := uc.generateOTP()
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Save user (not verified yet)
	user.IsVerified = false
	savedUser, err := uc.userRepository.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Store email verification
	verification := &domain.EmailVerification{
		Email:     user.Email,
		OTP:       otp,
		ExpiresAt: time.Now().Add(15 * time.Minute), // OTP expires in 15 minutes
		CreatedAt: time.Now(),
		Used:      false,
	}

	err = uc.emailVerificationRepo.Store(ctx, verification)
	if err != nil {
		return fmt.Errorf("failed to store verification: %w", err)
	}

	// Send verification email
	err = uc.emailService.SendVerificationEmail(ctx, user.Email, otp)
	if err != nil {
		log.Printf("Failed to send verification email to %s: %v", user.Email, err)
		// Don't fail registration if email sending fails
	}

	log.Printf("User %s registered successfully with ID: %s", savedUser.Email, savedUser.ID)
	return nil
}

func (uc *UserUsecase) VerifyEmail(ctx context.Context, input domain.VerifyEmailInput) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if input.Email == "" || input.OTP == "" {
		return domain.ErrInvalidInput
	}

	// Find verification record
	verification, err := uc.emailVerificationRepo.GetByOTP(ctx, input.OTP, input.Email)
	if err != nil {
		return fmt.Errorf("invalid or expired OTP: %w", err)
	}

	// Check if already used
	if verification.Used {
		return fmt.Errorf("OTP has already been used")
	}

	// Check if expired
	if time.Now().After(verification.ExpiresAt) {
		return fmt.Errorf("OTP has expired")
	}

	// Get user by email
	user, err := uc.userRepository.GetByEmail(ctx, input.Email)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Mark user as verified
	err = uc.userRepository.VerifyUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}

	// Mark verification as used
	err = uc.emailVerificationRepo.MarkUsed(ctx, verification.ID)
	if err != nil {
		log.Printf("Warning: failed to mark verification as used: %v", err)
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

	// Check if user is verified
	if !dbUser.IsVerified {
		return nil, fmt.Errorf("email not verified. Please verify your email before logging in")
	}

	if err := uc.passwordService.Compare(dbUser.Password, user.Password); err != nil {
		return nil, errors.New("incorrect email or password")
	}

	// Generate access token
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

	// Generate refresh token
	refreshToken, refreshTokenPayload, err := uc.JWTService.CreateRefreshToken(dbUser)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token
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

	// Delete all refresh tokens for the user
	err = uc.refreshTokenRepo.DeleteAllTokensForUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to logout user: %w", err)
	}

	return nil
}

func (uc *UserUsecase) Promote(ctx context.Context, username string) error {
	user, err := uc.userRepository.GetByUsername(ctx, username)
	if err != nil {
		return err
	}
	if user.Role == domain.RoleAdmin {
		return errors.New("user is already an admin")
	}
	return uc.userRepository.Promote(ctx, user)
}

func (uc *UserUsecase) UpdateProfile(ctx context.Context, userID string, bio, profilePicture, contactInfo string) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if userID == "" {
		return domain.ErrInvalidInput
	}

	user, err := uc.userRepository.GetByID(ctx, userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	updated := false
	if bio != "" && user.Bio != bio {
		user.Bio = bio
		updated = true
	}
	if profilePicture != "" && user.ProfilePicture != profilePicture {
		user.ProfilePicture = profilePicture
		updated = true
	}
	if contactInfo != "" && user.ContactInfo != contactInfo {
		user.ContactInfo = contactInfo
		updated = true
	}
	if !updated {
		return nil
	}

	user.UpdatedAt = time.Now()

	return uc.userRepository.Update(ctx, user)
}

func (uc *UserUsecase) RequestPasswordReset(ctx context.Context, input domain.RequestPasswordResetInput) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if input.Email == "" {
		return "", domain.ErrInvalidInput
	}

	user, err := uc.userRepository.GetByEmail(ctx, input.Email)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	rawToken, err := uc.passwordService.GenerateRandomToken()

	if err != nil {
		return "", fmt.Errorf("Error: %w", err)
	}

	hashedToken, err := uc.passwordService.Hash(rawToken)
	if err != nil {
		return "", err
	}

	resetToken := &domain.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
		Used:      false,
	}

	err = uc.passwordResetTokenRepo.Store(ctx, resetToken)
	if err != nil {
		return "", fmt.Errorf("failed to store reset token: %w", err)
	}

	err = uc.emailService.SendPasswordResetEmail(ctx, user.Email, rawToken)
	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return rawToken, nil
}

func (uc *UserUsecase) ResetPassword(ctx context.Context, input domain.ResetPasswordInput) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if input.Token == "" || input.NewPassword == "" {
		return domain.ErrInvalidInput
	}

	tokenRecord, err := uc.passwordResetTokenRepo.GetByTokenHash(ctx, input.Token)
	if err != nil {
		return fmt.Errorf("invalid or expired token: %w", err)
	}

	if err := uc.passwordService.ValidateStrength(input.NewPassword); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	hashedPassword, err := uc.passwordService.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	err = uc.userRepository.UpdatePassword(ctx, tokenRecord.UserID, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	err = uc.passwordResetTokenRepo.MarkUsed(ctx, tokenRecord.ID)
	if err != nil {
		return fmt.Errorf("failed to mark token used: %w", err)
	}

	return nil
}
