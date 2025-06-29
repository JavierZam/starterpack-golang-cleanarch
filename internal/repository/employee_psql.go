package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"starterpack-golang-cleanarch/internal/domain"

	"github.com/jmoiron/sqlx"
	// Tetap import ini jika GenerateID() di domain masih pakai uuid.New().String()
)

type postgreSQLEmployeeRepository struct {
	db *sqlx.DB
}

func NewPostgreSQLEmployeeRepository(db *sqlx.DB) domain.EmployeeRepository {
	return &postgreSQLEmployeeRepository{db: db}
}

// Save menyimpan data Employee baru ke database.
func (r *postgreSQLEmployeeRepository) Save(ctx context.Context, emp *domain.Employee) error {
	// ID (SERIAL) akan di-generate oleh database, jadi tidak perlu disertakan di sini.
	// PostgreSQL akan mengembalikan ID yang di-generate.
	// FIX: Tambahkan password_hash ke query INSERT
	query := `INSERT INTO users (tenant_id, name, email, phone_number, password_hash, created_at, updated_at)
              VALUES (:tenant_id, :name, :email, :phone_number, :password_hash, :created_at, :updated_at)
              RETURNING id`

	// NamedQueryContext dan Scan untuk mendapatkan ID yang di-generate
	rows, err := r.db.NamedQueryContext(ctx, query, emp)
	if err != nil {
		return fmt.Errorf("employeeRepo.Save: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&emp.ID) // Scan ID yang di-generate kembali ke struct Employee
		if err != nil {
			return fmt.Errorf("employeeRepo.Save: failed to scan ID: %w", err)
		}
	} else {
		return fmt.Errorf("employeeRepo.Save: no ID returned after insert")
	}

	return nil
}

// FindByID mencari Employee berdasarkan ID (int64) dan TenantID (string).
func (r *postgreSQLEmployeeRepository) FindByID(ctx context.Context, tenantID string, id int64) (*domain.Employee, error) {
	var emp domain.Employee
	// FIX: Tambahkan password_hash ke query SELECT
	query := `SELECT id, tenant_id, name, email, phone_number, password_hash, created_at, updated_at
              FROM users WHERE id = $1 AND tenant_id = $2`
	err := r.db.GetContext(ctx, &emp, query, id, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("employeeRepo.FindByID: %w", err)
	}
	return &emp, nil
}

// FindByEmail mencari Employee berdasarkan email dan TenantID (string).
func (r *postgreSQLEmployeeRepository) FindByEmail(ctx context.Context, tenantID string, email string) (*domain.Employee, error) {
	var emp domain.Employee
	// FIX: Tambahkan password_hash ke query SELECT
	query := `SELECT id, tenant_id, name, email, phone_number, password_hash, created_at, updated_at
              FROM users WHERE email = $1 AND tenant_id = $2`
	err := r.db.GetContext(ctx, &emp, query, email, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("employeeRepo.FindByEmail: %w", err)
	}
	return &emp, nil
}

// FindAll retrieves a list of Employees with pagination and filtering.
func (r *postgreSQLEmployeeRepository) FindAll(ctx context.Context, tenantID string, page, limit int, query string) (int64, []*domain.Employee, error) {
	offset := (page - 1) * limit
	var employees []*domain.Employee
	var total int64

	baseQuery := `FROM users WHERE tenant_id = $1`
	args := []interface{}{tenantID}
	argCounter := 2

	if query != "" {
		searchQuery := `(name ILIKE $` + strconv.Itoa(argCounter) + ` OR email ILIKE $` + strconv.Itoa(argCounter+1) + ` OR phone_number ILIKE $` + strconv.Itoa(argCounter+2) + `)`
		baseQuery += " AND " + searchQuery
		args = append(args, "%"+query+"%", "%"+query+"%", "%"+query+"%")
		argCounter += 3
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) %s`, baseQuery)
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return 0, nil, fmt.Errorf("employeeRepo.FindAll count: %w", err)
	}

	// FIX: Tambahkan password_hash ke query SELECT
	dataQuery := fmt.Sprintf(`SELECT id, tenant_id, name, email, phone_number, password_hash, created_at, updated_at
                              %s ORDER BY name ASC LIMIT $%d OFFSET $%d`,
		baseQuery, argCounter, argCounter+1)
	args = append(args, limit, offset)

	err = r.db.SelectContext(ctx, &employees, dataQuery, args...)
	if err != nil {
		return 0, nil, fmt.Errorf("employeeRepo.FindAll data: %w", err)
	}

	return total, employees, nil
}

// Update an existing Employee.
func (r *postgreSQLEmployeeRepository) Update(ctx context.Context, emp *domain.Employee) error {
	query := `UPDATE users SET name = :name, email = :email, phone_number = :phone_number, password_hash = :password_hash, updated_at = :updated_at
              WHERE id = :id AND tenant_id = :tenant_id`
	_, err := r.db.NamedExecContext(ctx, query, emp)
	if err != nil {
		return fmt.Errorf("employeeRepo.Update: %w", err)
	}
	return nil
}

// Delete an Employee by ID and TenantID.
func (r *postgreSQLEmployeeRepository) Delete(ctx context.Context, tenantID string, id int64) error {
	query := `DELETE FROM users WHERE id = $1 AND tenant_id = $2`
	_, err := r.db.ExecContext(ctx, query, id, tenantID)
	if err != nil {
		return fmt.Errorf("employeeRepo.Delete: %w", err)
	}
	return nil
}
