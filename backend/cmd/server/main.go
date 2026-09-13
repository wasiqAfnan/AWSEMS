package main

import (
	"log"
	"net/http"

	"awsems/internal/apigateway"
	"awsems/internal/cognito"
	"awsems/internal/config"
	"awsems/internal/middleware"
	"awsems/internal/routes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	cognitoClient := cognito.NewClient(cfg)
	if cognitoClient.Provider == nil {
		log.Fatal("Failed to initialize Cognito")
	}

	jwtVerifier, err := cognito.NewJWTVerifier(cfg)
	if err != nil {
		log.Fatalf("failed to initialize JWT verifier: %v", err)
	}

	log.Println("Cognito OIDC provider initialized successfully")

	// Create API Gateway client
	apiGatewayClient := apigateway.NewClient(cfg)

	// Create HTTP multiplexer
	mux := http.NewServeMux()

	// Register routes
	routes.Setup(mux, cfg, cognitoClient, jwtVerifier, apiGatewayClient)

	// Add CORS and logging middleware
	handler := middleware.CORS(cfg, middleware.Logging(mux))

	// Create server
	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: handler,
	}

	log.Printf("Server running on port %s", cfg.AppPort)

	log.Fatal(server.ListenAndServe())
}
