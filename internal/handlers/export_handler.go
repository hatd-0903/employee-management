package handlers

import (
	"net/http"

	"employee-management/internal/services"
)

type ExportHandler struct {
	service *services.ExportService
}

func NewExportHandler(service *services.ExportService) *ExportHandler {
	return &ExportHandler{service: service}
}

func (h *ExportHandler) Export(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.ExportAll(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}
