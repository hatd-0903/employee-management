package services

import (
	"context"
	"strings"
	"time"

	"employee-management/internal/models"
	"employee-management/internal/repositories"
	"employee-management/internal/utils"
)

const dbTimeout = 3 * time.Second

type EmployeeService struct {
	employees   repositories.EmployeeRepository
	departments repositories.DepartmentRepository
}

func NewEmployeeService(employees repositories.EmployeeRepository, departments repositories.DepartmentRepository) *EmployeeService {
	return &EmployeeService{employees: employees, departments: departments}
}

func (s *EmployeeService) Create(ctx context.Context, req models.CreateEmployeeRequest) (*models.Employee, error) {
	if err := utils.ValidateCreateEmployee(req); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	e := &models.Employee{
		Name:         strings.TrimSpace(req.Name),
		Age:          req.Age,
		Position:     strings.TrimSpace(req.Position),
		DepartmentID: req.DepartmentID,
		Salary:       req.Salary,
	}
	if err := s.employees.Create(ctx, e); err != nil {
		return nil, utils.AsAppError(err)
	}
	return e, nil
}

func (s *EmployeeService) GetByID(ctx context.Context, id int64) (*models.Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	e, err := s.employees.GetByID(ctx, id)
	if err != nil {
		return nil, utils.AsAppError(err)
	}
	return e, nil
}

func (s *EmployeeService) List(ctx context.Context, filter models.EmployeeListFilter) (*models.EmployeeListResult, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	employees, total, err := s.employees.List(ctx, filter)
	if err != nil {
		return nil, utils.AsAppError(err)
	}
	return &models.EmployeeListResult{TotalCount: total, Employees: employees}, nil
}

func (s *EmployeeService) Update(ctx context.Context, id int64, req models.UpdateEmployeeRequest) (*models.Employee, error) {
	if err := utils.ValidateUpdateEmployee(req); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	e, err := s.employees.Update(ctx, id, req)
	if err != nil {
		return nil, utils.AsAppError(err)
	}
	return e, nil
}

func (s *EmployeeService) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	if err := s.employees.Delete(ctx, id); err != nil {
		return utils.AsAppError(err)
	}
	return nil
}

func (s *EmployeeService) Search(ctx context.Context, keyword string) ([]*models.Employee, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, models.ErrValidation("keyword is required")
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	employees, err := s.employees.Search(ctx, keyword)
	if err != nil {
		return nil, utils.AsAppError(err)
	}
	return employees, nil
}

func (s *EmployeeService) ListByDepartment(ctx context.Context, departmentID int64) ([]*models.Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	if _, err := s.departments.GetByID(ctx, departmentID); err != nil {
		return nil, utils.AsAppError(err)
	}

	employees, err := s.employees.ListByDepartment(ctx, departmentID)
	if err != nil {
		return nil, utils.AsAppError(err)
	}
	return employees, nil
}
