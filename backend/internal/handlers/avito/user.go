package avito

import (
	"fmt"
	"io"
	"log"
	"mini-app-backend/internal/handlers"
	"mini-app-backend/utils"
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	reqURL := "https://api.avito.ru/core/v1/accounts/self"

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Printf("Error creating request to external API: %v", err)
		handlers.SendErrorResponse(w, "Error creating request to external API", http.StatusInternalServerError)
		return
	}

	token, _ := utils.GetToken()
	log.Printf("Token: %s", string(token.AccessToken))

	authHeader := fmt.Sprintf("%s %s", token.TokenType, token.AccessToken)

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to external API: %v", err)
		handlers.SendErrorResponse(w, "Error making request to external API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response from external API: %v", err)
		handlers.SendErrorResponse(w, "Error reading response from external API", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("External API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
		handlers.SendErrorResponse(w, fmt.Sprintf("External API returned status: %d", resp.StatusCode), resp.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
