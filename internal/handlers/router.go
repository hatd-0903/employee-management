package handlers

import "net/http"

func NewRouter(employees *EmployeeHandler, departments *DepartmentHandler, export *ExportHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthCheck)

	mux.HandleFunc("POST /employees", employees.Create)
	mux.HandleFunc("GET /employees", employees.List)
	mux.HandleFunc("GET /employees/search", employees.Search)
	mux.HandleFunc("POST /employees/export", export.Export)
	mux.HandleFunc("GET /employees/{id}", employees.GetByID)
	mux.HandleFunc("PUT /employees/{id}", employees.Update)
	mux.HandleFunc("DELETE /employees/{id}", employees.Delete)

	mux.HandleFunc("POST /departments", departments.Create)
	mux.HandleFunc("GET /departments", departments.List)
	mux.HandleFunc("GET /departments/{id}/employees", departments.ListEmployees)

	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
