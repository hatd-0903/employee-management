package handlers

import (
	"encoding/json"
	"net/http"

	"employee-management/internal/models"
	"employee-management/internal/utils"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, err error) {
	appErr := utils.AsAppError(err)
	writeJSON(w, appErr.Code, appErr)
}

func decodeJSON(r *http.Request, dst any) *models.AppError {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return models.ErrValidation("invalid request body: " + err.Error())
	}
	return nil
}
