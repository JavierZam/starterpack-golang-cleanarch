package auth

import (
	"net/http"

	"starterpack-golang-cleanarch/internal/utils/errors"
)

var (
	ErrUserNotFound        = errors.New("USER_NOT_FOUND", "User with given email not found", http.StatusNotFound, nil, nil)
	ErrUserAlreadyExists   = errors.New("USER_ALREADY_EXISTS", "User with this email already exists", http.StatusConflict, nil, nil)
	ErrInvalidCredentials  = errors.New("INVALID_CREDENTIALS", "Invalid email or password", http.StatusUnauthorized, nil, nil)
	ErrInvalidToken        = errors.New("INVALID_TOKEN", "Invalid or expired token", http.StatusUnauthorized, nil, nil)
	ErrRefreshTokenExpired = errors.New("REFRESH_TOKEN_EXPIRED", "Refresh token has expired, please login again", http.StatusUnauthorized, nil, nil)
)
