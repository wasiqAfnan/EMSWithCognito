package routes

import (
	"net/http"

	"awsems/internal/apigateway"
	"awsems/internal/cognito"
	"awsems/internal/config"
	"awsems/internal/handler"
	"awsems/internal/middleware"
)

func Setup(
	mux *http.ServeMux,
	cfg *config.Config,
	jwtVerifier *cognito.JWTVerifier,
	apiGatewayClient *apigateway.Client,
) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("AWSEMS server is running"))
	})

	// Authentication middleware for protected routes
	authMiddleware := middleware.Authentication(jwtVerifier)

	// Auth routes (protected)
	authHandler := handler.NewAuthHandler(apiGatewayClient, jwtVerifier)
	mux.Handle(
		"GET /api/me",
		authMiddleware(http.HandlerFunc(authHandler.Me)),
	)

	// Employee routes (protected)
	employeeHandler := handler.NewEmployeeHandler(apiGatewayClient)

	mux.Handle(
		"GET /api/employees",
		authMiddleware(http.HandlerFunc(employeeHandler.GetEmployees)),
	)
	mux.Handle(
		"POST /api/employees",
		authMiddleware(http.HandlerFunc(employeeHandler.CreateEmployee)),
	)
	mux.Handle(
		"GET /api/employees/search",
		authMiddleware(http.HandlerFunc(employeeHandler.SearchEmployees)),
	)
	mux.Handle(
		"PATCH /api/employees/{empId}",
		authMiddleware(http.HandlerFunc(employeeHandler.UpdateEmployee)),
	)
	mux.Handle(
		"DELETE /api/employees/{empId}",
		authMiddleware(http.HandlerFunc(employeeHandler.DeleteEmployee)),
	)
}
