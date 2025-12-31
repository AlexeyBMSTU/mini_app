package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mini-app-backend/internal/config"
	"net/http"
	"strings"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetToken() (TokenResponse, error) {
	tokenURL := "https://api.avito.ru/token/"

	config := config.Load()
	formData := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", config.AvitoClientId, config.AvitoClientSecret)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData))
	if err != nil {
		log.Printf("Error creating request to external API: %v", err)
		return TokenResponse{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to external API: %v", err)
		return TokenResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response from external API: %v", err)
		return TokenResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("External API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
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
