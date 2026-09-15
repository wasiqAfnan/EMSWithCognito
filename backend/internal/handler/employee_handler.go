package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"awsems/internal/apigateway"
	"awsems/internal/model"
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

// CreateEmployee reads the employee data from the request body, proxies it to the Lambda via API Gateway,
// and streams the raw JSON response back to the client.
func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	// Parse payload to validate
	var employee model.EmployeeCreate
	if err := json.Unmarshal(payload, &employee); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Validate employee data
	if err := utils.ValidateEmployeeCreate(&employee); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	body, gwStatusCode, err := h.APIGateway.Post("/employees", payload)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to reach API Gateway")
		return
	}

	// The Lambda might return a 200 HTTP status but contain an error status code in the JSON body.
	var resp struct {
		StatusCode  int `json:"statusCode"`
		Status_Code int `json:"status_code"` // handle both cases
	}
	finalStatusCode := gwStatusCode
	if err := json.Unmarshal(body, &resp); err == nil {
		if resp.StatusCode >= 400 {
			finalStatusCode = resp.StatusCode
		} else if resp.Status_Code >= 400 {
			finalStatusCode = resp.Status_Code
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(finalStatusCode)
	w.Write(body)
}

// UpdateEmployee reads the update payload, validates it, proxies it to Lambda via API Gateway,
// and streams the raw JSON response back to the client.
func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	empId := r.PathValue("empId")
	if empId == "" {
		utils.Error(w, http.StatusBadRequest, "Missing empId in path")
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	// Parse payload to validate
	var employeeUpdate model.EmployeeUpdate
	if err := json.Unmarshal(payload, &employeeUpdate); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Validate employee update data
	if err := utils.ValidateEmployeeUpdate(&employeeUpdate); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	path := "/employees/" + url.PathEscape(empId)
	body, gwStatusCode, err := h.APIGateway.Patch(path, payload)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to reach API Gateway")
		return
	}

	var resp struct {
		StatusCode  int `json:"statusCode"`
		Status_Code int `json:"status_code"`
	}
	finalStatusCode := gwStatusCode
	if err := json.Unmarshal(body, &resp); err == nil {
		if resp.StatusCode >= 400 {
			finalStatusCode = resp.StatusCode
		} else if resp.Status_Code >= 400 {
			finalStatusCode = resp.Status_Code
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(finalStatusCode)
	w.Write(body)
}

// DeleteEmployee proxies the delete request to the Lambda via API Gateway.
func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	empId := r.PathValue("empId")
	if empId == "" {
		utils.Error(w, http.StatusBadRequest, "Missing empId in path")
		return
	}

	path := "/employees/" + url.PathEscape(empId)
	body, statusCode, err := h.APIGateway.Delete(path)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to reach API Gateway")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(body)
}
