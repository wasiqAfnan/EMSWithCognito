package handler

import (
	"encoding/json"
	"net/http"
	"strings"

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
	// 1. Extract Access Token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Error(w, http.StatusUnauthorized, "Missing or invalid authorization header")
		return
	}
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	// 2. Extract ID Token
	idTokenStr := r.Header.Get("X-Id-Token")
	if idTokenStr == "" {
		utils.Error(w, http.StatusBadRequest, "Missing X-Id-Token header")
		return
	}

	// 3. Verify Access Token to get sub
	parsedAccToken, err := h.jwtVerifier.VerifyAccessToken(accessToken)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Invalid access token")
		return
	}

	claims, ok := parsedAccToken.Claims.(jwt.MapClaims)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "Invalid token claims")
		return
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		utils.Error(w, http.StatusUnauthorized, "sub missing from access token")
		return
	}

	// 4. Call API Gateway to get user by sub
	userJSON, statusCode, err := h.apiGatewayClient.Get("/users/" + sub)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to call backend services")
		return
	}

	if statusCode == 200 {
		// User exists!
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

	// 5. User not found (404), provision new user
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

	if email == "" || name == "" {
		utils.Error(w, http.StatusBadRequest, "ID token missing required claims (email, name)")
		return
	}

	if phoneNumber == "" {
		phoneNumber = "" // Ensure phone number can be safely empty if not provided by Cognito depending on config
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

	// 6. Call API Gateway to create user
	createdJSON, createStatus, err := h.apiGatewayClient.Post("/users", payloadBytes)
	if err != nil || createStatus != 201 {
		utils.Error(w, http.StatusInternalServerError, "Failed to create user in backend")
		return
	}

	// 7. Return created user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(createdJSON)
}
