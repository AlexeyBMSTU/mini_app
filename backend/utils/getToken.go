package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mini-app-backend/internal/config"
	"mini-app-backend/internal/logger"
	"mini-app-backend/internal/middleware"
	"net/http"
	"strings"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetToken(ctx context.Context) (TokenResponse, error) {
	tokenURL := "https://api.avito.ru/token/"

	// Пытаемся получить client_id и client_secret из контекста
	clientID, clientIDOk := ctx.Value("avito_client_id").(string)
	clientSecret, clientSecretOk := ctx.Value("avito_client_secret").(string)
	
	// Если в контексте нет данных, используем значения из конфигурации
	var cfg *config.Config
	if !clientIDOk || !clientSecretOk {
		cfg = config.Load()
		if !clientIDOk {
			clientID = cfg.AvitoClientId
		}
		if !clientSecretOk {
			clientSecret = cfg.AvitoClientSecret
		}
	}

	formData := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", clientID, clientSecret)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData))
	if err != nil {
		logger.Errorf("Error creating request to external API: %v", err)
		return TokenResponse{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	requestLogger := logger.GetLogger()
	if requestID := ctx.Value(middleware.RequestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok {
			requestLogger = logger.WithRequestID(id)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		requestLogger.Errorf("Error making request to external API: %v", err)
		return TokenResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		requestLogger.Errorf("Error reading response from external API: %v", err)
		return TokenResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		requestLogger.Errorf("External API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
		return TokenResponse{}, err
	}

	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return TokenResponse{}, fmt.Errorf("error unmarshaling token response: %v", err)
	}

	if tokenResponse.AccessToken == "" {
		return TokenResponse{}, fmt.Errorf("empty access token in response")
	}

	return tokenResponse, nil
}