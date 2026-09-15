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

	CognitoClientID        string
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
		CognitoClientID:        os.Getenv("COGNITO_CLIENT_ID"),
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

	config.CognitoClientID, err = getParameter("/ems/prod/cognito-client-id")
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
		"COGNITO_CLIENT_ID":         config.CognitoClientID,
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
