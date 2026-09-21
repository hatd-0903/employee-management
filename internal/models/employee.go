package models

import (
	"encoding/json"
	"time"
)

// Salary is kept as the exact decimal text (via json.Number) rather than
// float64 end to end, so it round-trips through the MySQL DECIMAL(12,2)
// column and the JSON API without ever going through binary float rounding.
// It still marshals/unmarshals as a plain JSON number on the wire.
type Salary = json.Number

type Employee struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Age          int        `json:"age"`
	Position     string     `json:"position"`
	DepartmentID int64      `json:"departmentId"`
	Salary       Salary     `json:"salary"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty"`
}

// CreateEmployeeRequest is the payload for POST /employees.
type CreateEmployeeRequest struct {
	Name         string `json:"name"`
	Age          int    `json:"age"`
	Position     string `json:"position"`
	DepartmentID int64  `json:"departmentId"`
	Salary       Salary `json:"salary"`
}

// UpdateEmployeeRequest supports partial updates: only non-nil fields are applied.
type UpdateEmployeeRequest struct {
	Name         *string `json:"name"`
	Age          *int    `json:"age"`
	Position     *string `json:"position"`
	DepartmentID *int64  `json:"departmentId"`
	Salary       *Salary `json:"salary"`
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
