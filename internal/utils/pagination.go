package utils

import (
	"net/http"
	"strconv"

	"employee-management/internal/models"
)

const (
	defaultLimit = 10
	maxLimit     = 100
)

func ParseEmployeeListFilter(r *http.Request) (models.EmployeeListFilter, *models.AppError) {
	q := r.URL.Query()
	filter := models.EmployeeListFilter{Limit: defaultLimit, Offset: 0}

	if raw := q.Get("limit"); raw != "" {
		limit, err := strconv.Atoi(raw)
		if err != nil || limit <= 0 {
			return filter, models.ErrValidation("limit must be a positive integer")
		}
		if limit > maxLimit {
			limit = maxLimit
		}
		filter.Limit = limit
	}

	if raw := q.Get("offset"); raw != "" {
		offset, err := strconv.Atoi(raw)
		if err != nil || offset < 0 {
			return filter, models.ErrValidation("offset must be a non-negative integer")
		}
		filter.Offset = offset
	}

	if raw := q.Get("departmentId"); raw != "" {
		departmentID, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || departmentID <= 0 {
			return filter, models.ErrValidation("departmentId must be a positive integer")
		}
		filter.DepartmentID = &departmentID
	}

	return filter, nil
}
