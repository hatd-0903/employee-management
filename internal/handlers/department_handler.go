package handlers

import (
	"net/http"

	"employee-management/internal/models"
	"employee-management/internal/services"
	"employee-management/internal/utils"
)

type DepartmentHandler struct {
	departments *services.DepartmentService
	employees   *services.EmployeeService
}

func NewDepartmentHandler(departments *services.DepartmentService, employees *services.EmployeeService) *DepartmentHandler {
	return &DepartmentHandler{departments: departments, employees: employees}
}

func (h *DepartmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateDepartmentRequest
	if appErr := decodeJSON(r, &req); appErr != nil {
		writeError(w, appErr)
		return
	}

	department, err := h.departments.Create(r.Context(), req)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, department)
}

func (h *DepartmentHandler) List(w http.ResponseWriter, r *http.Request) {
	departments, err := h.departments.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"departments": departments})
}

func (h *DepartmentHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	id, appErr := utils.ParseID(r.PathValue("id"))
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	employees, err := h.employees.ListByDepartment(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"employees": employees})
}
