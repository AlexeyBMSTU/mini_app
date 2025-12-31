package avito

import (
	"fmt"
	"io"
	"log"
	"mini-app-backend/utils"
	"mini-app-backend/internal/handlers"
	"net/http"
)

func GetMesseges(w http.ResponseWriter, r *http.Request) {	
	user_id, _:= utils.GetUserInfoAvito()

	reqURL := fmt.Sprintf("https://api.avito.ru/messenger/v2/accounts/%d/chats", user_id.Id)
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

func SendMessege(w http.ResponseWriter, r *http.Request) {	
	queryParams := r.URL.Query()
	userId := queryParams.Get("user_id")
	chatId := queryParams.Get("chat_id")

	reqURL := fmt.Sprintf("https://api.avito.ru/messenger/v1/accounts/%s/chats/%s/messages", userId, chatId)

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
