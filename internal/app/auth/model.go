package auth

type RegisterRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	TenantID    string `json:"tenant_id" validate:"required,uuid"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UserResponse struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenant_id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	CreatedAt   string `json:"created_at"`
}
