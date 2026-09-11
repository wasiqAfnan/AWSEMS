package routes

import (
	"net/http"

	"awsems/internal/cognito"
	"awsems/internal/handler"
)

func Setup(
	mux *http.ServeMux,
	cognitoClient *cognito.Client,
	jwtVerifier *cognito.JWTVerifier,
) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("AWSEMS server is running"))
	})

	// Authentication routes
	authHandler := handler.NewAuthHandler(cognitoClient)

	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/callback", authHandler.Callback)

	// Authentication middleware for protected routes
	// authMiddleware := middleware.Authentication(jwtVerifier)

	// Employee routes will be protected with authMiddleware
}
