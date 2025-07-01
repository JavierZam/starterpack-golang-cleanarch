package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"starterpack-golang-cleanarch/internal/domain" // Pastikan baris import ini ada dan benar

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type postgreSQLUserRepository struct {
	db *sqlx.DB
}

func NewPostgreSQLUserRepository(db *sqlx.DB) domain.UserRepository {
	return &postgreSQLUserRepository{db: db}
}

func (r *postgreSQLUserRepository) Save(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, tenant_id, email, password_hash, name, phone_number, role, created_at, updated_at)
              VALUES (:id, :tenant_id, :email, :password_hash, :name, :phone_number, :role, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("userRepo.Save: %w", err)
	}
	return nil
}

func (r *postgreSQLUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, tenant_id, email, password_hash, name, phone_number, role, created_at, updated_at
              FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("userRepo.FindByEmail: %w", err)
	}
	return &user, nil
}

func (r *postgreSQLUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, tenant_id, email, password_hash, name, phone_number, role, created_at, updated_at
              FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("userRepo.FindByID: %w", err)
	}
	return &user, nil
}

func (r *postgreSQLUserRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()
	query := `UPDATE users SET email = :email, password_hash = :password_hash, name = :name, phone_number = :phone_number, role = :role, updated_at = :updated_at
              WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("userRepo.Update: %w", err)
	}
	return nil
}

func (r *postgreSQLUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("userRepo.Delete: %w", err)
	}
	return nil
}
