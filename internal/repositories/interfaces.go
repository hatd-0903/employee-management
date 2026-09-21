package repositories

import (
	"context"

	"employee-management/internal/models"
)

// EmployeeRepository is the storage abstraction for employees. Defined as an
// interface so services can be tested against an in-memory fake instead of a
// real database.
type EmployeeRepository interface {
	Create(ctx context.Context, e *models.Employee) error
	GetByID(ctx context.Context, id int64) (*models.Employee, error)
	List(ctx context.Context, filter models.EmployeeListFilter) ([]*models.Employee, int, error)
	Update(ctx context.Context, id int64, req models.UpdateEmployeeRequest) (*models.Employee, error)
	Delete(ctx context.Context, id int64) error
	Search(ctx context.Context, keyword string, filter models.EmployeeListFilter) ([]*models.Employee, int, error)
	ListByDepartment(ctx context.Context, departmentID int64, filter models.EmployeeListFilter) ([]*models.Employee, int, error)
}

// DepartmentRepository is the storage abstraction for departments.
type DepartmentRepository interface {
	Create(ctx context.Context, d *models.Department) error
	List(ctx context.Context) ([]*models.Department, error)
	GetByID(ctx context.Context, id int64) (*models.Department, error)
}
