package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"auth_service/internal/config"
	"auth_service/internal/db"
	"auth_service/internal/repository"
	"auth_service/internal/utils"
)

// SignUpRequest для регистрации
type SignUpRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type AuthService struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.TokenRepository
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (s *AuthService) SignUp(ctx context.Context, req db.SignUpRequest) (*db.User, *db.Tokens, error) {
	existingUser, err := s.userRepo.GetByPhoneNumber(ctx, req.PhoneNumber)
	if existingUser != nil {
		return nil, nil, fmt.Errorf("phone number already registered")
	}

	if req.Email != "" {
		existingUser, err = s.userRepo.GetByEmail(ctx, req.Email)
		if existingUser != nil {
			return nil, nil, fmt.Errorf("email already registered")
		}
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &db.User{
		Name:        strings.TrimSpace(req.Name),
		PhoneNumber: req.PhoneNumber,
		Email:       sql.NullString{String: req.Email, Valid: req.Email != ""},
		Password:    hashedPassword,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	tokens, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user, tokens, nil
}

func (s *AuthService) SignIn(ctx context.Context, req db.LoginRequest) (*db.User, *db.Tokens, error) {
	user, err := s.userRepo.GetByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return nil, nil, errors.New("invalid phone number or password")
	}

	if user == nil {
		return nil, nil, errors.New("invalid phone number or password")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, nil, errors.New("invalid phone number or password")
	}

	tokens, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user, tokens, nil
}

func (s *AuthService) Logout(ctx context.Context, userID int64, accessToken string) error {
	err := s.tokenRepo.DeleteRefreshToken(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	ttl := parseDuration(config.App.JWT.AccessTTL)
	err = s.tokenRepo.StoreBlacklistedToken(ctx, accessToken, ttl)

	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*db.Tokens, error) {
	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	blacklisted, err := s.tokenRepo.IsTokenBlacklisted(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	if blacklisted {
		return nil, errors.New("token is blacklisted")
	}

	storedToken, err := s.tokenRepo.GetRefreshToken(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("refresh token not found")
	}

	if storedToken != refreshToken {
		return nil, errors.New("refresh token mismatch")
	}

	tokens, err := s.generateTokens(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	return tokens, nil
}

func (s *AuthService) generateTokens(userID int64) (*db.Tokens, error) {
	accessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)

	}

	ctx := context.Background()
	err = s.tokenRepo.StoreRefreshToken(ctx, userID, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &db.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func parseDuration(durationStr string) time.Duration {
	dur, err := time.ParseDuration(durationStr)
	if err != nil {
		return 15 * time.Minute // default
	}
	return dur
}
