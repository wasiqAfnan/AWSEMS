package routes

import (
	"net/http"

	"awsems/internal/apigateway"
	"awsems/internal/cognito"
	"awsems/internal/handler"
	"awsems/internal/middleware"
)

func Setup(
	mux *http.ServeMux,
	cognitoClient *cognito.Client,
	jwtVerifier *cognito.JWTVerifier,
	apiGatewayClient *apigateway.Client,
) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("AWSEMS server is running"))
	})

	// Authentication routes
	authHandler := handler.NewAuthHandler(cognitoClient)

	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/callback", authHandler.Callback)

	// Authentication middleware for protected routes
	authMiddleware := middleware.Authentication(jwtVerifier)

	// Authentication routes (protected)
	mux.Handle(
		"/api/auth/logout",
		authMiddleware(http.HandlerFunc(authHandler.Logout)),
	)

	// Employee routes (protected)
	employeeHandler := handler.NewEmployeeHandler(apiGatewayClient)

	mux.Handle(
		"/api/employees",
		authMiddleware(http.HandlerFunc(employeeHandler.GetEmployees)),
	)
	mux.Handle(
		"/api/employees/search",
		authMiddleware(http.HandlerFunc(employeeHandler.SearchEmployees)),
	)
}
