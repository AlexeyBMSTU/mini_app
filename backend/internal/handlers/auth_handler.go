package handlers

import (
	"database/sql"
	"encoding/json"
	"mini-app-backend/internal/errors"
	"mini-app-backend/internal/telegram"
	"mini-app-backend/internal/user"
	"net/http"
	"strconv"
)

type AuthHandler struct {
	*BaseHandler
	userService *user.UserService
	botToken    string
	db          *sql.DB
}

func NewAuthHandler(userService *user.UserService, botToken string, db *sql.DB) *AuthHandler {
	return &AuthHandler{
		BaseHandler: NewBaseHandler(),
		userService: userService,
		botToken:    botToken,
		db:          db,
	}
}

type TelegramAuthRequest struct {
	User     *user.User `json:"user"`
	InitData string     `json:"initData"`
}

type TelegramAuthResponse struct {
	Success bool       `json:"success"`
	User    *user.User `json:"user,omitempty"`
	Error   string     `json:"error,omitempty"`
	Token   string     `json:"token,omitempty"`
}

func (h *AuthHandler) TelegramAuth(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "TelegramAuth request")

	if r.Method != http.MethodPost {
		h.LogError(r, nil, "Method not allowed")
		h.SendError(w, r, errors.NewAppError(http.StatusMethodNotAllowed, "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var req TelegramAuthRequest
	err := h.DecodeJSONBody(r, &req)
	if err != nil {
		h.LogError(r, err, "Invalid request body")
		h.SendError(w, r, err, http.StatusBadRequest)
		return
	}

	if req.User == nil || req.InitData == "" {
		h.LogError(r, nil, "User and initData are required")
		h.SendError(w, r, errors.NewAppError(http.StatusBadRequest, "User and initData are required"), http.StatusBadRequest)
		return
	}

	isValid, _, err := telegram.ValidateAndParseInitData(req.InitData, h.botToken)
	if err != nil {
		h.LogError(r, err, "Error validating initData")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error validating initData", err.Error()), http.StatusInternalServerError)
		return
	}

	if !isValid {
		h.LogError(r, nil, "Invalid initData")
		h.SendError(w, r, errors.NewAppError(http.StatusUnauthorized, "Invalid initData"), http.StatusUnauthorized)
		return
	}

	createdUser, err := h.userService.CreateOrUpdateUser(req.User)
	if err != nil {
		h.LogError(r, err, "Error creating/updating user")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error creating/updating user", err.Error()), http.StatusInternalServerError)
		return
	}

	token := "token_" + strconv.FormatInt(createdUser.ID, 10)

	response := TelegramAuthResponse{
		Success: true,
		User:    createdUser,
		Token:   token,
	}

	h.LogInfo(r, "Successfully authenticated user")
	h.SendJSON(w, r, response, http.StatusOK)
}

func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "GetUser request")

	if r.Method != http.MethodGet {
		h.LogError(r, nil, "Method not allowed")
		h.SendError(w, r, errors.NewAppError(http.StatusMethodNotAllowed, "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.GetUserIDFromHeader(r)
	if err != nil {
		h.LogError(r, err, "Failed to get user ID from header")
		h.SendError(w, r, err, http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		h.LogError(r, err, "Error getting user")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error getting user", err.Error()), http.StatusInternalServerError)
		return
	}

	if user == nil {
		h.LogError(r, nil, "User not found")
		h.SendError(w, r, errors.NewAppError(http.StatusNotFound, "User not found"), http.StatusNotFound)
		return
	}

	h.LogInfo(r, "Successfully retrieved user")
	h.SendJSON(w, r, user, http.StatusOK)
}

func (h *AuthHandler) GetUserData(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "GetUserData request")

	if r.Method != http.MethodGet {
		h.LogError(r, nil, "Method not allowed")
		h.SendError(w, r, errors.NewAppError(http.StatusMethodNotAllowed, "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.GetUserIDFromHeader(r)
	if err != nil {
		h.LogError(r, err, "Failed to get user ID from header")
		h.SendError(w, r, err, http.StatusBadRequest)
		return
	}

	userData, err := h.userService.GetUserDataByUserID(userID)
	if err != nil {
		h.LogError(r, err, "Error getting user data")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error getting user data", err.Error()), http.StatusInternalServerError)
		return
	}

	h.LogInfo(r, "Successfully retrieved user data")
	h.SendJSON(w, r, userData, http.StatusOK)
}

func (h *AuthHandler) SaveUserData(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "SaveUserData request")

	if r.Method != http.MethodPost {
		h.LogError(r, nil, "Method not allowed")
		h.SendError(w, r, errors.NewAppError(http.StatusMethodNotAllowed, "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.GetUserIDFromHeader(r)
	if err != nil {
		h.LogError(r, err, "Failed to get user ID from header")
		h.SendError(w, r, err, http.StatusBadRequest)
		return
	}

	var data map[string]interface{}
	err = h.DecodeJSONBody(r, &data)
	if err != nil {
		h.LogError(r, err, "Invalid request body")
		h.SendError(w, r, err, http.StatusBadRequest)
		return
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		h.LogError(r, err, "Error marshaling data")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error marshaling data", err.Error()), http.StatusInternalServerError)
		return
	}

	_, err = h.userService.SaveUserData(userID, string(dataJSON))
	if err != nil {
		h.LogError(r, err, "Error saving user data")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error saving user data", err.Error()), http.StatusInternalServerError)
		return
	}

	h.LogInfo(r, "Successfully saved user data")
	h.SendJSON(w, r, map[string]bool{"success": true}, http.StatusOK)
}
