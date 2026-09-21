package handlers

import (
	"net/http"

	"employee-management/internal/models"
	"employee-management/internal/services"
	"employee-management/internal/utils"
)

type EmployeeHandler struct {
	service *services.EmployeeService
}

func NewEmployeeHandler(service *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEmployeeRequest
	if appErr := decodeJSON(r, &req); appErr != nil {
		writeError(w, appErr)
		return
	}

	employee, err := h.service.Create(r.Context(), req)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, employee)
}

func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	filter, appErr := utils.ParseEmployeeListFilter(r)
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	result, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *EmployeeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, appErr := utils.ParseID(r.PathValue("id"))
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	employee, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, employee)
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, appErr := utils.ParseID(r.PathValue("id"))
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	var req models.UpdateEmployeeRequest
	if appErr := decodeJSON(r, &req); appErr != nil {
		writeError(w, appErr)
		return
	}

	employee, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, employee)
}

func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, appErr := utils.ParseID(r.PathValue("id"))
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *EmployeeHandler) Search(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")

	employees, err := h.service.Search(r.Context(), keyword)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"employees": employees})
}
