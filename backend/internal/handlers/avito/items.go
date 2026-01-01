package avito

import (
	"fmt"
	"mini-app-backend/internal/errors"
	"mini-app-backend/internal/logger"
	"mini-app-backend/utils"
	"net/http"
)

var avitoHandler *AvitoHandler

func GetItems(w http.ResponseWriter, r *http.Request) {
	if avitoHandler == nil {
		avitoHandler = NewAvitoHandler()
	}
	avitoHandler.GetItems(w, r)
}

func (h *AvitoHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "GetItems request")

	queryParams := r.URL.Query()
	category := queryParams.Get("category")

	if category == "" {
		category = "338"
	}

	baseURL := "https://api.avito.ru/core/v1/items"
	reqURL := fmt.Sprintf("%s?category=%s", baseURL, category)

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

	requestID := r.Header.Get("X-Request-ID")
	logger.WithRequestID(requestID).Info("Successfully2 retrieved items")
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp.Body)
}
