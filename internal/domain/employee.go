package domain

import (
	"context"
	"time"

	"github.com/google/uuid" // Tetap digunakan untuk generate UUID string untuk TenantID jika diperlukan.
)

// Employee represents the core business entity for an employee.
type Employee struct {
	ID           int64     `db:"id"`        // Changed to int64 to match SERIAL PRIMARY KEY
	TenantID     string    `db:"tenant_id"` // Changed to string to match VARCHAR(36)
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PhoneNumber  string    `db:"phone_number"`
	PasswordHash string    `db:"password_hash"` // <-- Added this field
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// GenerateID now populates ID (left to DB auto-increment) and creates UUID string for TenantID.
func (e *Employee) GenerateID() {
	// ID (SERIAL) akan di-generate oleh database saat insert, jadi tidak perlu diisi di sini.
	// Jika TenantID masih kosong, generate UUID baru sebagai string.
	if e.TenantID == "" {
		e.TenantID = uuid.New().String()
	}
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
}

// UpdateTimestamp updates the UpdatedAt field.
func (e *Employee) UpdateTimestamp() {
	e.UpdatedAt = time.Now()
}

// EmployeeRepository defines the interface for data access operations for Employee.
type EmployeeRepository interface {
	Save(ctx context.Context, emp *Employee) error
	FindByID(ctx context.Context, tenantID string, id int64) (*Employee, error)                              // ID changed to int64, TenantID to string
	FindByEmail(ctx context.Context, tenantID string, email string) (*Employee, error)                       // TenantID to string
	FindAll(ctx context.Context, tenantID string, page, limit int, query string) (int64, []*Employee, error) // TenantID to string
	Update(ctx context.Context, emp *Employee) error
	Delete(ctx context.Context, tenantID string, id int64) error // ID changed to int64, TenantID to string
}
