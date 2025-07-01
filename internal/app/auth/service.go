package auth

import (
	"context"
	"fmt"
	"time"

	"starterpack-golang-cleanarch/internal/domain"
	"starterpack-golang-cleanarch/internal/utils"
	globalErrors "starterpack-golang-cleanarch/internal/utils/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo domain.UserRepository
}

func NewAuthService(repo domain.UserRepository) *AuthService {
	return &AuthService{userRepo: repo}
}

func (s *AuthService) RegisterUser(ctx context.Context, req RegisterRequest) (*UserResponse, error) {
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, globalErrors.NewInternalServerError(fmt.Errorf("failed to check existing user: %w", err), "Internal error during user registration check.")
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, globalErrors.NewInternalServerError(fmt.Errorf("failed to hash password: %w", err), "Internal error during password hashing.")
	}

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		PhoneNumber:  req.PhoneNumber,
		TenantID:     uuid.MustParse(req.TenantID),
		Role:         "user",
	}
	user.GenerateID()

	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, globalErrors.NewInternalServerError(fmt.Errorf("failed to save user: %w", err), "Internal error saving user.")
	}

	resp := &UserResponse{
		ID:          user.ID.String(),
		TenantID:    user.TenantID.String(),
		Email:       user.Email,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt.Format(utils.ISO8601TimeFormat),
	}

	return resp, nil
}

func (s *AuthService) LoginUser(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, globalErrors.NewInternalServerError(fmt.Errorf("failed to find user by email: %w", err), "Internal error during login.")
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.String(), user.TenantID.String(), user.Role)
	if err != nil {
		return nil, globalErrors.NewInternalServerError(fmt.Errorf("failed to generate access token: %w", err), "Internal error generating token.")
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, globalErrors.NewInternalServerError(fmt.Errorf("failed to generate refresh token: %w", err), "Internal error generating token.")
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserResponse{
			ID:          user.ID.String(),
			TenantID:    user.TenantID.String(),
			Email:       user.Email,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt.Format(utils.ISO8601TimeFormat),
		},
	}, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, req RefreshTokenRequest) (*AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrRefreshTokenExpired
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, globalErrors.NewBadRequest("Invalid user ID in refresh token claims", nil)
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, globalErrors.NewInternalServerError(fmt.Errorf("failed to find user for refresh token: %w", err), "Internal error during token refresh.")
	}
	if user == nil {
		return nil, ErrInvalidToken
	}

	newAccessToken, err := utils.GenerateAccessToken(user.ID.String(), user.TenantID.String(), user.Role)
	if err != nil {
		return nil, globalErrors.NewInternalServerError(fmt.Errorf("failed to generate new access token: %w", err), "Internal error generating token.")
	}
	newRefreshToken, err := utils.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, globalErrors.NewInternalServerError(fmt.Errorf("failed to generate new refresh token: %w", err), "Internal error generating token.")
	}

	return &AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User: UserResponse{
			ID:          user.ID.String(),
			TenantID:    user.TenantID.String(),
			Email:       user.Email,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt.Format(utils.ISO8601TimeFormat),
		},
	}, nil
}
