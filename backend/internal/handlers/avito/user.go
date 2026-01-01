package avito

import (
	"mini-app-backend/internal/errors"
	"mini-app-backend/internal/handlers"
	"mini-app-backend/internal/httpclient"
	"mini-app-backend/internal/logger"
	"mini-app-backend/utils"
	"net/http"
)

type AvitoHandler struct {
	*handlers.BaseHandler
	httpClient *httpclient.Client
}

func NewAvitoHandler() *AvitoHandler {
	return &AvitoHandler{
		BaseHandler: handlers.NewBaseHandler(),
		httpClient: httpclient.NewClient(
			httpclient.WithLogger(logger.GetLogger()),
		),
	}
}

func (h *AvitoHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "GetUser request")

	reqURL := "https://api.avito.ru/core/v1/accounts/self"

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

	h.LogInfo(r, "Successfully retrieved user data")
	h.SendJSON(w, r, resp.Body, http.StatusOK)
}
