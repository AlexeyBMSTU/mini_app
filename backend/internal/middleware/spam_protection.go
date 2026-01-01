package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"mini-app-backend/internal/errors"
	"mini-app-backend/internal/message"
	"net/http"
	"time"
)

const (
	maxMessagesPerUser      = 10
	maxMessagesPerMinute    = 5
	rateLimitWindowDuration = time.Minute
)

type SpamProtectionMiddleware struct {
	messageService *message.MessageService
}

func NewSpamProtectionMiddleware(messageService *message.MessageService) *SpamProtectionMiddleware {
	return &SpamProtectionMiddleware{
		messageService: messageService,
	}
}

func (m *SpamProtectionMiddleware) Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/message/" {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			errors.SendErrorResponse(w, errors.NewAppError(http.StatusInternalServerError, "Error reading request body"))
			return
		}
		
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		var req struct {
			ClientID string `json:"client_id"`
		}
		
		if err := json.Unmarshal(body, &req); err != nil || req.ClientID == "" {
			errors.SendErrorResponse(w, errors.NewAppError(http.StatusBadRequest, "client_id is required"))
			return
		}
		
		clientID := req.ClientID

		totalMessages, err := m.messageService.CountMessagesByClientID(clientID)
		if err != nil {
			errors.SendErrorResponse(w, errors.NewAppError(http.StatusInternalServerError, "Error checking message count"))
			return
		}

		if totalMessages >= maxMessagesPerUser {
			errors.SendErrorResponse(w, errors.NewAppError(http.StatusTooManyRequests, "Message limit exceeded. Maximum 10 messages allowed."))
			return
		}

		recentMessages, err := m.messageService.CountMessagesByClientIDInTimeRange(clientID, rateLimitWindowDuration)
		if err != nil {
			errors.SendErrorResponse(w, errors.NewAppError(http.StatusInternalServerError, "Error checking recent message count"))
			return
		}

		if recentMessages >= maxMessagesPerMinute {
			errors.SendErrorResponse(w, errors.NewAppError(http.StatusTooManyRequests, "Rate limit exceeded. Please wait before sending another message."))
			return
		}

		next.ServeHTTP(w, r)
	})
}
