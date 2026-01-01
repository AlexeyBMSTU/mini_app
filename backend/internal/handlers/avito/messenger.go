package avito

import (
	"fmt"
	"mini-app-backend/internal/errors"
	"mini-app-backend/utils"
	"net/http"
)

func GetMesseges(w http.ResponseWriter, r *http.Request) {
	if avitoHandler == nil {
		avitoHandler = NewAvitoHandler()
	}
	avitoHandler.GetMesseges(w, r)
}

func SendMessege(w http.ResponseWriter, r *http.Request) {
	if avitoHandler == nil {
		avitoHandler = NewAvitoHandler()
	}
	avitoHandler.SendMessege(w, r)
}

func (h *AvitoHandler) GetMesseges(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "GetMesseges request")

	user_id, err := utils.GetUserInfoAvito(r.Context())
	if err != nil {
		h.LogError(r, err, "Failed to get user info")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Failed to get user info", err.Error()), http.StatusInternalServerError)
		return
	}

	reqURL := fmt.Sprintf("https://api.avito.ru/messenger/v2/accounts/%d/chats", user_id.Id)

	token, err := utils.GetToken(r.Context())
	if err != nil {
		h.LogError(r, err, "Failed to get token")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Failed to get token", err.Error()), http.StatusInternalServerError)
		return
	}

	h.LogDebug(r, "Token obtained")

	authHeader := token.TokenType + " " + token.AccessToken

	h.LogDebug(r, fmt.Sprintf(("authToken: %s"), authHeader))

	resp, err := h.httpClient.Get(r.Context(), reqURL, map[string]string{
		"Authorization": authHeader,
	})
	if err != nil {
		h.LogError(r, err, "Error making request to external API")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error making request to external API", err.Error()), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		h.LogError(r, nil, "External API returned non-OK status")
		h.SendError(w, r, errors.NewAppErrorWithDetails(resp.StatusCode, "External API returned non-OK status", string(resp.Body)), resp.StatusCode)
		return
	}

	h.SendJSON(w, r, resp.Body, http.StatusOK)
}

func (h *AvitoHandler) SendMessege(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "SendMessege request")

	queryParams := r.URL.Query()
	userId := queryParams.Get("user_id")
	chatId := queryParams.Get("chat_id")

	if userId == "" || chatId == "" {
		h.LogError(r, nil, "Missing required parameters")
		h.SendError(w, r, errors.NewAppError(http.StatusBadRequest, "Missing required parameters"), http.StatusBadRequest)
		return
	}

	reqURL := fmt.Sprintf("https://api.avito.ru/core/v1/messenger/chats/%s/messages/%s", userId, chatId)

	token, err := utils.GetToken(r.Context())
	if err != nil {
		h.LogError(r, err, "Failed to get token")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Failed to get token", err.Error()), http.StatusInternalServerError)
		return
	}

	h.LogDebug(r, "Token obtained")

	authHeader := token.TokenType + " " + token.AccessToken

	resp, err := h.httpClient.Get(r.Context(), reqURL, map[string]string{
		"Authorization": authHeader,
	})
	if err != nil {
		h.LogError(r, err, "Error making request to external API")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error making request to external API", err.Error()), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		h.LogError(r, nil, "External API returned non-OK status")
		h.SendError(w, r, errors.NewAppErrorWithDetails(resp.StatusCode, "External API returned non-OK status", string(resp.Body)), resp.StatusCode)
		return
	}

	h.SendJSON(w, r, resp.Body, http.StatusOK)
}
