package services

import (
	"context"
	"strings"

	"employee-management/internal/models"
	"employee-management/internal/repositories"
	"employee-management/internal/utils"
)

type DepartmentService struct {
	departments repositories.DepartmentRepository
}

func NewDepartmentService(departments repositories.DepartmentRepository) *DepartmentService {
	return &DepartmentService{departments: departments}
}

func (s *DepartmentService) Create(ctx context.Context, req models.CreateDepartmentRequest) (*models.Department, error) {
	if err := utils.ValidateCreateDepartment(req); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	d := &models.Department{Name: strings.TrimSpace(req.Name)}
	if err := s.departments.Create(ctx, d); err != nil {
		return nil, utils.AsAppError(err)
	}
	return d, nil
}

func (s *DepartmentService) List(ctx context.Context) ([]*models.Department, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	departments, err := s.departments.List(ctx)
	if err != nil {
		return nil, utils.AsAppError(err)
	}
	return departments, nil
}

func (s *DepartmentService) GetByID(ctx context.Context, id int64) (*models.Department, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	d, err := s.departments.GetByID(ctx, id)
	if err != nil {
		return nil, utils.AsAppError(err)
	}
	return d, nil
}
