package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"employee-management/internal/models"

	"github.com/go-sql-driver/mysql"
)

type mysqlDepartmentRepository struct {
	db *sql.DB
}

// NewMySQLDepartmentRepository builds a DepartmentRepository backed by MySQL.
func NewMySQLDepartmentRepository(db *sql.DB) DepartmentRepository {
	return &mysqlDepartmentRepository{db: db}
}

func (r *mysqlDepartmentRepository) Create(ctx context.Context, d *models.Department) error {
	stmt, err := r.db.PrepareContext(ctx, `INSERT INTO departments (name) VALUES (?)`)
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, d.Name)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return models.ErrConflict(fmt.Sprintf("department %q already exists", d.Name))
		}
		return fmt.Errorf("insert department: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("read last insert id: %w", err)
	}

	row := r.db.QueryRowContext(ctx, `SELECT id, name, created_at FROM departments WHERE id = ?`, id)
	if err := row.Scan(&d.ID, &d.Name, &d.CreatedAt); err != nil {
		return fmt.Errorf("read created department: %w", err)
	}
	return nil
}

func (r *mysqlDepartmentRepository) List(ctx context.Context) ([]*models.Department, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, created_at FROM departments ORDER BY id ASC`)
	if err != nil {
		return nil, fmt.Errorf("list departments: %w", err)
	}
	defer rows.Close()

	departments := make([]*models.Department, 0)
	for rows.Next() {
		var d models.Department
		if err := rows.Scan(&d.ID, &d.Name, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan department: %w", err)
		}
		departments = append(departments, &d)
	}
	return departments, rows.Err()
}

func (r *mysqlDepartmentRepository) GetByID(ctx context.Context, id int64) (*models.Department, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, created_at FROM departments WHERE id = ?`, id)
	var d models.Department
	err := row.Scan(&d.ID, &d.Name, &d.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound(fmt.Sprintf("department %d not found", id))
	}
	if err != nil {
		return nil, fmt.Errorf("select department: %w", err)
	}
	return &d, nil
}
