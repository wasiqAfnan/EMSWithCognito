package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"awsems/internal/apigateway"
	"awsems/internal/cognito"
	"awsems/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	apiGatewayClient *apigateway.Client
	jwtVerifier      *cognito.JWTVerifier
}

func NewAuthHandler(apiGatewayClient *apigateway.Client, jwtVerifier *cognito.JWTVerifier) *AuthHandler {
	return &AuthHandler{
		apiGatewayClient: apiGatewayClient,
		jwtVerifier:      jwtVerifier,
	}
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// 1. Read sub from context (set by auth middleware)
	sub, ok := r.Context().Value("user_sub").(string)
	if !ok || sub == "" {
		utils.Error(w, http.StatusUnauthorized, "sub missing from context")
		return
	}

	// 2. Call API Gateway to get user by sub
	userJSON, statusCode, err := h.apiGatewayClient.Get("/users/"+sub, nil)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to call backend services")
		return
	}

	if statusCode == 200 {
		// User exists! Return immediately
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(userJSON)
		return
	}

	if statusCode != 404 {
		// Some other error
		utils.Error(w, statusCode, "Unexpected error fetching user")
		return
	}

	// 3. User not found (404), read and verify X-Id-Token for provisioning
	idTokenStr := r.Header.Get("X-Id-Token")
	if idTokenStr == "" {
		utils.Error(w, http.StatusBadRequest, "Missing X-Id-Token header for user provisioning")
		return
	}

	parsedIDToken, err := h.jwtVerifier.VerifyIDToken(idTokenStr)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Invalid ID token for provisioning")
		return
	}

	idClaims, ok := parsedIDToken.Claims.(jwt.MapClaims)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "Invalid ID token claims")
		return
	}

	email, _ := idClaims["email"].(string)
	name, _ := idClaims["name"].(string)
	phoneNumber, _ := idClaims["phone_number"].(string)

	if email == "" || name == "" || phoneNumber == "" {
		utils.Error(w, http.StatusBadRequest, "ID token missing required claims (email, name)")
		return
	}

	// Prepare CreateUser payload
	payload := map[string]interface{}{
		"cognito_sub":  sub,
		"email":        email,
		"name":         name,
		"phone_number": phoneNumber,
		"role":         "USER",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to marshal payload")
		return
	}
	fmt.Printf("[CreateUser] Final JSON payload being sent to API Gateway: %s\n", string(payloadBytes))
	// 4. Call API Gateway to create user
	createdJSON, createStatus, err := h.apiGatewayClient.Post("/users", payloadBytes)
	if err != nil || createStatus != 201 {
		utils.Error(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create user in backend: status=%d, err=%v, body=%s", createStatus, err, string(createdJSON)))
		return
	}

	// 5. Return created user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(createdJSON)
}
