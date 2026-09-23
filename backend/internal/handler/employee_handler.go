package handler

import (
	"encoding/json"
	"fmt"
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

	path := "/employees"

	payload, err := utils.EmbedCreatedBy(r.Context(), nil)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Sending the injected payload to the Lambda via API Gateway
	body, statusCode, err := h.APIGateway.Get(path, payload)
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

	payload, err := utils.EmbedCreatedBy(r.Context(), nil)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	body, statusCode, err := h.APIGateway.Get(path, payload)
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

	// Extract payload from request body
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	// 1. Inject created_by from context before unmarshaling to model
	enrichedPayload, err := utils.EmbedCreatedBy(r.Context(), payload)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	// 2. Parse enriched payload to EmployeeCreate struct
	var employee model.EmployeeCreate
	if err := json.Unmarshal(enrichedPayload, &employee); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// 3. Validate employee data (now including the injected created_by)
	if err := utils.ValidateEmployeeCreate(&employee); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Send the validated struct to API Gateway to ensure no missing fields
	finalPayload, _ := json.Marshal(employee)
	fmt.Printf("[CreateEmployee Go Handler] Final JSON payload being sent to API Gateway: %s\n", string(finalPayload))
	body, gwStatusCode, err := h.APIGateway.Post("/employees", finalPayload)
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

	// 1. Inject created_by from context before unmarshaling to model
	enrichedPayload, err := utils.EmbedCreatedBy(r.Context(), payload)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	// 2. Parse payload to validate
	var employeeUpdate model.EmployeeUpdate
	if err := json.Unmarshal(enrichedPayload, &employeeUpdate); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// 3. Validate employee update data
	if err := utils.ValidateEmployeeUpdate(&employeeUpdate); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	path := "/employees/" + url.PathEscape(empId)

	// 4. Send enriched payload to API Gateway
	body, gwStatusCode, err := h.APIGateway.Patch(path, enrichedPayload)
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

	payload, err := utils.EmbedCreatedBy(r.Context(), nil)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	body, statusCode, err := h.APIGateway.Delete(path, payload)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to reach API Gateway")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(body)
}
