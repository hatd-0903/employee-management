package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"employee-management/internal/models"
)

type mysqlEmployeeRepository struct {
	db *sql.DB
}

// NewMySQLEmployeeRepository builds an EmployeeRepository backed by MySQL.
func NewMySQLEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &mysqlEmployeeRepository{db: db}
}

func (r *mysqlEmployeeRepository) Create(ctx context.Context, e *models.Employee) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var departmentExists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM departments WHERE id = ?)`, e.DepartmentID).Scan(&departmentExists)
	if err != nil {
		return fmt.Errorf("check department: %w", err)
	}
	if !departmentExists {
		return models.ErrValidation(fmt.Sprintf("department %d does not exist", e.DepartmentID))
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO employees (name, age, position, department_id, salary)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, e.Name, e.Age, e.Position, e.DepartmentID, e.Salary)
	if err != nil {
		return fmt.Errorf("insert employee: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("read last insert id: %w", err)
	}

	row := tx.QueryRowContext(ctx, `
		SELECT id, name, age, position, department_id, salary, created_at, updated_at
		FROM employees WHERE id = ?
	`, id)
	if err := scanEmployee(row, e); err != nil {
		return fmt.Errorf("read created employee: %w", err)
	}

	return tx.Commit()
}

func (r *mysqlEmployeeRepository) GetByID(ctx context.Context, id int64) (*models.Employee, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, name, age, position, department_id, salary, created_at, updated_at
		FROM employees WHERE id = ? AND deleted_at IS NULL
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare select: %w", err)
	}
	defer stmt.Close()

	var e models.Employee
	err = scanEmployee(stmt.QueryRowContext(ctx, id), &e)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound(fmt.Sprintf("employee %d not found", id))
	}
	if err != nil {
		return nil, fmt.Errorf("select employee: %w", err)
	}
	return &e, nil
}

func (r *mysqlEmployeeRepository) List(ctx context.Context, filter models.EmployeeListFilter) ([]*models.Employee, int, error) {
	where := "WHERE deleted_at IS NULL"
	args := []any{}
	if filter.DepartmentID != nil {
		where += " AND department_id = ?"
		args = append(args, *filter.DepartmentID)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM employees " + where
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count employees: %w", err)
	}

	listQuery := fmt.Sprintf(`
		SELECT id, name, age, position, department_id, salary, created_at, updated_at
		FROM employees %s
		ORDER BY id ASC
		LIMIT ? OFFSET ?
	`, where)
	listArgs := append(append([]any{}, args...), filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list employees: %w", err)
	}
	defer rows.Close()

	employees := make([]*models.Employee, 0)
	for rows.Next() {
		var e models.Employee
		if err := scanEmployee(rows, &e); err != nil {
			return nil, 0, fmt.Errorf("scan employee: %w", err)
		}
		employees = append(employees, &e)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate employees: %w", err)
	}

	return employees, total, nil
}

func (r *mysqlEmployeeRepository) Update(ctx context.Context, id int64, req models.UpdateEmployeeRequest) (*models.Employee, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var exists bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM employees WHERE id = ? AND deleted_at IS NULL)`, id).Scan(&exists); err != nil {
		return nil, fmt.Errorf("check employee: %w", err)
	}
	if !exists {
		return nil, models.ErrNotFound(fmt.Sprintf("employee %d not found", id))
	}

	sets := make([]string, 0, 5)
	args := make([]any, 0, 6)

	if req.Name != nil {
		sets = append(sets, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Age != nil {
		sets = append(sets, "age = ?")
		args = append(args, *req.Age)
	}
	if req.Position != nil {
		sets = append(sets, "position = ?")
		args = append(args, *req.Position)
	}
	if req.DepartmentID != nil {
		var departmentExists bool
		if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM departments WHERE id = ?)`, *req.DepartmentID).Scan(&departmentExists); err != nil {
			return nil, fmt.Errorf("check department: %w", err)
		}
		if !departmentExists {
			return nil, models.ErrValidation(fmt.Sprintf("department %d does not exist", *req.DepartmentID))
		}
		sets = append(sets, "department_id = ?")
		args = append(args, *req.DepartmentID)
	}
	if req.Salary != nil {
		sets = append(sets, "salary = ?")
		args = append(args, *req.Salary)
	}

	if len(sets) > 0 {
		args = append(args, id)
		query := fmt.Sprintf(`UPDATE employees SET %s WHERE id = ?`, strings.Join(sets, ", "))
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("prepare update: %w", err)
		}
		defer stmt.Close()
		if _, err := stmt.ExecContext(ctx, args...); err != nil {
			return nil, fmt.Errorf("update employee: %w", err)
		}
	}

	var e models.Employee
	row := tx.QueryRowContext(ctx, `
		SELECT id, name, age, position, department_id, salary, created_at, updated_at
		FROM employees WHERE id = ?
	`, id)
	if err := scanEmployee(row, &e); err != nil {
		return nil, fmt.Errorf("read updated employee: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	return &e, nil
}

func (r *mysqlEmployeeRepository) Delete(ctx context.Context, id int64) error {
	stmt, err := r.db.PrepareContext(ctx, `
		UPDATE employees SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL
	`)
	if err != nil {
		return fmt.Errorf("prepare delete: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("delete employee: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read rows affected: %w", err)
	}
	if affected == 0 {
		return models.ErrNotFound(fmt.Sprintf("employee %d not found", id))
	}
	return nil
}

func (r *mysqlEmployeeRepository) Search(ctx context.Context, keyword string) ([]*models.Employee, error) {
	like := "%" + keyword + "%"
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, age, position, department_id, salary, created_at, updated_at
		FROM employees
		WHERE deleted_at IS NULL AND (name LIKE ? OR position LIKE ?)
		ORDER BY id ASC
	`, like, like)
	if err != nil {
		return nil, fmt.Errorf("search employees: %w", err)
	}
	defer rows.Close()

	employees := make([]*models.Employee, 0)
	for rows.Next() {
		var e models.Employee
		if err := scanEmployee(rows, &e); err != nil {
			return nil, fmt.Errorf("scan employee: %w", err)
		}
		employees = append(employees, &e)
	}
	return employees, rows.Err()
}

func (r *mysqlEmployeeRepository) ListByDepartment(ctx context.Context, departmentID int64) ([]*models.Employee, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT e.id, e.name, e.age, e.position, e.department_id, e.salary, e.created_at, e.updated_at
		FROM employees e
		INNER JOIN departments d ON d.id = e.department_id
		WHERE d.id = ? AND e.deleted_at IS NULL
		ORDER BY e.id ASC
	`, departmentID)
	if err != nil {
		return nil, fmt.Errorf("list employees by department: %w", err)
	}
	defer rows.Close()

	employees := make([]*models.Employee, 0)
	for rows.Next() {
		var e models.Employee
		if err := scanEmployee(rows, &e); err != nil {
			return nil, fmt.Errorf("scan employee: %w", err)
		}
		employees = append(employees, &e)
	}
	return employees, rows.Err()
}

// rowScanner abstracts over *sql.Row and *sql.Rows so scanEmployee works for both.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanEmployee(row rowScanner, e *models.Employee) error {
	return row.Scan(&e.ID, &e.Name, &e.Age, &e.Position, &e.DepartmentID, &e.Salary, &e.CreatedAt, &e.UpdatedAt)
}
