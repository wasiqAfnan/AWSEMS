package handler

import (
	"net/http"
	"net/url"

	"awsems/internal/apigateway"
	"awsems/internal/utils"
)

type EmployeeHandler struct {
	APIGateway *apigateway.Client
}

func NewEmployeeHandler(apiGatewayClient *apigateway.Client) *EmployeeHandler {
	return &EmployeeHandler{
		APIGateway: apiGatewayClient,
	}
}

// GetEmployees proxies the request to the Lambda via API Gateway and
// streams the raw JSON response back to the client as-is.
func (h *EmployeeHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	body, statusCode, err := h.APIGateway.Get("/employees")
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to reach API Gateway")
		return
	}

	// The Lambda already returns a properly structured JSON response.
	// Forward it directly without re-encoding.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(body)
}

// SearchEmployees proxies the search request to the Lambda via API Gateway.
func (h *EmployeeHandler) SearchEmployees(w http.ResponseWriter, r *http.Request) {
	// log.Println("Searching employees")
	if r.Method != http.MethodGet {
		utils.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		utils.Error(w, http.StatusBadRequest, "Missing search query parameter 'q'")
		return
	}

	path := "/employees/search?q=" + url.QueryEscape(query)
	body, statusCode, err := h.APIGateway.Get(path)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to reach API Gateway")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(body)
}
