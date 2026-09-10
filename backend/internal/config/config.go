package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	CognitoRegion       string
	CognitoUserPoolID   string
	CognitoClientID     string
	CognitoClientSecret string
	CognitoDomain       string

	APIGatewayBaseURL string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load .env: %w", err)
	}

	config := &Config{
		AppEnv:              os.Getenv("APP_ENV"),
		AppPort:             os.Getenv("APP_PORT"),
		CognitoRegion:       os.Getenv("COGNITO_REGION"),
		CognitoUserPoolID:   os.Getenv("COGNITO_USER_POOL_ID"),
		CognitoClientID:     os.Getenv("COGNITO_CLIENT_ID"),
		CognitoClientSecret: os.Getenv("COGNITO_CLIENT_SECRET"),
		CognitoDomain:       os.Getenv("COGNITO_DOMAIN"),
		APIGatewayBaseURL:   os.Getenv("API_GATEWAY_BASE_URL"),
	}

	if err := validate(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validate(config *Config) error {
	required := map[string]string{
		"APP_ENV":               config.AppEnv,
		"APP_PORT":              config.AppPort,
		"COGNITO_REGION":        config.CognitoRegion,
		"COGNITO_USER_POOL_ID":  config.CognitoUserPoolID,
		"COGNITO_CLIENT_ID":     config.CognitoClientID,
		"COGNITO_CLIENT_SECRET": config.CognitoClientSecret,
		"COGNITO_DOMAIN":        config.CognitoDomain,
		"API_GATEWAY_BASE_URL":  config.APIGatewayBaseURL,
	}

	for name, value := range required {
		if value == "" {
			return fmt.Errorf("required environment variable %s is not set", name)
		}
	}

	return nil
}
