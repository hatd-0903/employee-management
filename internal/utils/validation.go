package utils

import (
	"strconv"
	"strings"

	"employee-management/internal/models"
)

func ValidateCreateEmployee(req models.CreateEmployeeRequest) *models.AppError {
	if strings.TrimSpace(req.Name) == "" {
		return models.ErrValidation("name is required")
	}
	if req.Age <= 0 {
		return models.ErrValidation("age must be greater than 0")
	}
	if strings.TrimSpace(req.Position) == "" {
		return models.ErrValidation("position is required")
	}
	if req.DepartmentID <= 0 {
		return models.ErrValidation("departmentId is required")
	}
	if req.Salary <= 0 {
		return models.ErrValidation("salary must be greater than 0")
	}
	return nil
}

func ValidateUpdateEmployee(req models.UpdateEmployeeRequest) *models.AppError {
	if req.Name == nil && req.Age == nil && req.Position == nil && req.DepartmentID == nil && req.Salary == nil {
		return models.ErrValidation("at least one field must be provided")
	}
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return models.ErrValidation("name cannot be empty")
	}
	if req.Age != nil && *req.Age <= 0 {
		return models.ErrValidation("age must be greater than 0")
	}
	if req.Position != nil && strings.TrimSpace(*req.Position) == "" {
		return models.ErrValidation("position cannot be empty")
	}
	if req.DepartmentID != nil && *req.DepartmentID <= 0 {
		return models.ErrValidation("departmentId must be valid")
	}
	if req.Salary != nil && *req.Salary <= 0 {
		return models.ErrValidation("salary must be greater than 0")
	}
	return nil
}

func ValidateCreateDepartment(req models.CreateDepartmentRequest) *models.AppError {
	if strings.TrimSpace(req.Name) == "" {
		return models.ErrValidation("name is required")
	}
	return nil
}

// ParseID validates a path parameter as a positive int64 id.
func ParseID(raw string) (int64, *models.AppError) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, models.ErrValidation("invalid id")
	}
	return id, nil
}
