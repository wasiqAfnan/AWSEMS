package config

import (
	"fmt"
	"os"

	"awsems/internal/utils"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	CognitoRegion          string
	CognitoUserPoolID      string
	CognitoClientID        string
	CognitoClientSecret    string
	CognitoDomain          string
	CognitoRedirectURI     string
	CognitoLogoutURI       string
	CognitoOpenIDConfigURL string

	APIGatewayBaseURL string
	FrontendURL       string
}

func Load() (*Config, error) {
	appEnv := os.Getenv("APP_ENV")

	if appEnv == "" || appEnv == "development" {
		if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load .env: %w", err)
		}

		appEnv = os.Getenv("APP_ENV")
	}

	if appEnv == "production" {
		return loadProductionConfig()
	}

	config := &Config{
		AppEnv:                 appEnv,
		AppPort:                os.Getenv("APP_PORT"),
		CognitoRegion:          os.Getenv("COGNITO_REGION"),
		CognitoUserPoolID:      os.Getenv("COGNITO_USER_POOL_ID"),
		CognitoClientID:        os.Getenv("COGNITO_CLIENT_ID"),
		CognitoClientSecret:    os.Getenv("COGNITO_CLIENT_SECRET"),
		CognitoDomain:          os.Getenv("COGNITO_DOMAIN"),
		CognitoRedirectURI:     os.Getenv("COGNITO_REDIRECT_URI"),
		CognitoLogoutURI:       os.Getenv("COGNITO_LOGOUT_URI"),
		CognitoOpenIDConfigURL: os.Getenv("COGNITO_OPENID_CONFIG_URL"),
		APIGatewayBaseURL:      os.Getenv("API_GATEWAY_BASE_URL"),
		FrontendURL:            os.Getenv("FRONTEND_URL"),
	}

	if err := validate(config); err != nil {
		return nil, err
	}

	return config, nil
}

func loadProductionConfig() (*Config, error) {
	config := &Config{
		AppEnv:  "production",
		AppPort: os.Getenv("APP_PORT"),
	}

	var err error

	config.CognitoRegion, err = getParameter("/ems/prod/cognito-region")
	if err != nil {
		return nil, err
	}

	config.CognitoUserPoolID, err = getParameter("/ems/prod/cognito-user-pool-id")
	if err != nil {
		return nil, err
	}

	config.CognitoClientID, err = getParameter("/ems/prod/cognito-client-id")
	if err != nil {
		return nil, err
	}

	config.CognitoClientSecret, err = getParameter("/ems/prod/cognito-client-secret")
	if err != nil {
		return nil, err
	}

	config.CognitoDomain, err = getParameter("/ems/prod/cognito-domain")
	if err != nil {
		return nil, err
	}

	config.CognitoRedirectURI, err = getParameter("/ems/prod/cognito-redirect-uri")
	if err != nil {
		return nil, err
	}

	config.CognitoLogoutURI, err = getParameter("/ems/prod/cognito-logout-uri")
	if err != nil {
		return nil, err
	}

	config.CognitoOpenIDConfigURL, err = getParameter("/ems/prod/cognito-openid-config-url")
	if err != nil {
		return nil, err
	}

	config.APIGatewayBaseURL, err = getParameter("/ems/prod/api-gateway-base-url")
	if err != nil {
		return nil, err
	}

	config.FrontendURL, err = getParameter("/ems/prod/frontend-url")
	if err != nil {
		return nil, err
	}

	if err := validate(config); err != nil {
		return nil, err
	}

	return config, nil
}

func getParameter(name string) (string, error) {
	return utils.GetParameter(name, true)
}

func validate(config *Config) error {
	required := map[string]string{
		"APP_ENV":                   config.AppEnv,
		"APP_PORT":                  config.AppPort,
		"COGNITO_REGION":            config.CognitoRegion,
		"COGNITO_USER_POOL_ID":      config.CognitoUserPoolID,
		"COGNITO_CLIENT_ID":         config.CognitoClientID,
		"COGNITO_CLIENT_SECRET":     config.CognitoClientSecret,
		"COGNITO_DOMAIN":            config.CognitoDomain,
		"COGNITO_REDIRECT_URI":      config.CognitoRedirectURI,
		"COGNITO_LOGOUT_URI":        config.CognitoLogoutURI,
		"COGNITO_OPENID_CONFIG_URL": config.CognitoOpenIDConfigURL,
		"API_GATEWAY_BASE_URL":      config.APIGatewayBaseURL,
		"FRONTEND_URL":              config.FrontendURL,
	}

	for name, value := range required {
		if value == "" {
			return fmt.Errorf("required configuration %s is not set", name)
		}
	}

	return nil
}
