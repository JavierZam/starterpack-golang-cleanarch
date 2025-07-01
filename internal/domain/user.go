package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	TenantID     uuid.UUID `db:"tenant_id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Name         string    `db:"name"`
	PhoneNumber  string    `db:"phone_number"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (u *User) GenerateID() {
	u.ID = uuid.New()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) UpdateTimestamp() {
	u.UpdatedAt = time.Now()
}

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
