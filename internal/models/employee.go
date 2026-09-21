package models

import "time"

type Employee struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Age          int        `json:"age"`
	Position     string     `json:"position"`
	DepartmentID int64      `json:"departmentId"`
	Salary       float64    `json:"salary"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty"`
}

// CreateEmployeeRequest is the payload for POST /employees.
type CreateEmployeeRequest struct {
	Name         string  `json:"name"`
	Age          int     `json:"age"`
	Position     string  `json:"position"`
	DepartmentID int64   `json:"departmentId"`
	Salary       float64 `json:"salary"`
}

// UpdateEmployeeRequest supports partial updates: only non-nil fields are applied.
type UpdateEmployeeRequest struct {
	Name         *string  `json:"name"`
	Age          *int     `json:"age"`
	Position     *string  `json:"position"`
	DepartmentID *int64   `json:"departmentId"`
	Salary       *float64 `json:"salary"`
}

// EmployeeListFilter carries pagination and filter parameters for GET /employees.
type EmployeeListFilter struct {
	Limit        int
	Offset       int
	DepartmentID *int64
}

type EmployeeListResult struct {
	TotalCount int         `json:"totalCount"`
	Employees  []*Employee `json:"employees"`
}
