package handlers

import (
	"encoding/json"
	"mini-app-backend/internal/errors"
	"mini-app-backend/internal/logger"
	"mini-app-backend/internal/utils"
	"net/http"
	"strconv"
)

type contextKey string

const RequestIDKey contextKey = "requestID"

type BaseHandler struct {
	logger logger.Logger
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{
		logger: logger.GetLogger(),
	}
}

func (h *BaseHandler) getLogger(r *http.Request) logger.Logger {
	requestID := r.Context().Value(RequestIDKey)
	if requestID != nil {
		return h.logger.WithRequestID(requestID.(string))
	}
	return h.logger
}

func (h *BaseHandler) SendJSON(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.getLogger(r).Errorf("Error encoding response: %v", err)
	}
}

func (h *BaseHandler) SendError(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	errors.SendErrorResponse(w, err)
}

func (h *BaseHandler) GetUserIDFromHeader(r *http.Request) (int64, error) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		return 0, errors.NewAppError(http.StatusBadRequest, "X-User-ID is required")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, errors.NewAppError(http.StatusBadRequest, "Invalid User ID")
	}

	return userID, nil
}

func (h *BaseHandler) GetUserIDFromCookie(r *http.Request) (int64, error) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		return 0, errors.NewAppError(http.StatusUnauthorized, "User ID not found in context")
	}

	return userID.(int64), nil
}

func (h *BaseHandler) SetUserCookie(w http.ResponseWriter, userID int64) {
	encryptionUtil := utils.NewEncryptionUtil()
	
	userIDStr := strconv.FormatInt(userID, 10)
	h.logger.Infof("Encrypting user ID: %s", userIDStr)
	
	encryptedUserID, err := encryptionUtil.Encrypt(userIDStr)
	if err != nil {
		h.logger.Errorf("Error encrypting user ID: %v", err)
		encryptedUserID = userIDStr
	}
	
	h.logger.Infof("Encrypted user ID: %s", encryptedUserID)
	
	cookie := &http.Cookie{
		Name:     "user_id",
		Value:    encryptedUserID,
		Path:     "/",
		MaxAge:   86400 * 30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	h.logger.Infof("Set encrypted cookie for user ID: %d", userID)
}

func (h *BaseHandler) DecodeJSONBody(r *http.Request, target interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		h.getLogger(r).Errorf("Error decoding request body: %v", err)
		return errors.NewAppError(http.StatusBadRequest, "Invalid request body")
	}

	return nil
}

func (h *BaseHandler) LogRequest(r *http.Request, message string) {
	h.getLogger(r).Infof("%s: %s %s", message, r.Method, r.URL.Path)
}

func (h *BaseHandler) LogError(r *http.Request, err error, message string) {
	if err != nil {
		h.getLogger(r).Errorf("%s: %s %s - %v", message, r.Method, r.URL.Path, err)
	} else {
		h.getLogger(r).Errorf("%s: %s %s", message, r.Method, r.URL.Path)
	}
}

func (h *BaseHandler) LogInfo(r *http.Request, message string) {
	h.getLogger(r).Infof("%s: %s %s", message, r.Method, r.URL.Path)
}

func (h *BaseHandler) LogDebug(r *http.Request, message string) {
	h.getLogger(r).Debugf("%s: %s %s", message, r.Method, r.URL.Path)
}