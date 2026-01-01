package handlers

import (
	"database/sql"
	"encoding/json"
	"mini-app-backend/internal/config"
	"mini-app-backend/internal/errors"
	"mini-app-backend/internal/message"
	"mini-app-backend/internal/telegram"
	"mini-app-backend/internal/user"
	"net/http"
	"strconv"
)

type AuthHandler struct {
	*BaseHandler
	userService    *user.UserService
	messageService *message.MessageService
	botToken       string
	db             *sql.DB
	config         *config.Config
}

func NewAuthHandler(userService *user.UserService, messageService *message.MessageService, botToken string, db *sql.DB, config *config.Config) *AuthHandler {
	return &AuthHandler{
		BaseHandler:    NewBaseHandler(),
		userService:    userService,
		messageService: messageService,
		botToken:       botToken,
		db:             db,
		config:         config,
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

	existingUser, err := h.userService.GetUserByID(req.User.ID)
	if err != nil {
		h.LogError(r, err, "Error checking user existence")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error checking user existence", err.Error()), http.StatusInternalServerError)
		return
	}
	
	if existingUser != nil {
		h.LogError(r, nil, "User already exists")
		h.SendError(w, r, errors.NewAppError(http.StatusBadRequest, "User already exists"), http.StatusBadRequest)
		return
	}
	
	createdUser, err := h.userService.CreateOrUpdateUser(req.User)
	if err != nil {
		h.LogError(r, err, "Error creating user")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error creating user", err.Error()), http.StatusInternalServerError)
		return
	}

	_, err = h.messageService.CreateMessage(
		h.config.AvitoClientId,
		h.config.AvitoClientSecret,
		"Ваше сообщение по умолчанию",
		"Автоответчик",
	)
	if err != nil {
		h.LogError(r, err, "Error creating default message")
	}

	token := "token_" + strconv.FormatInt(createdUser.ID, 10)

	h.SetUserCookie(w, createdUser.ID)

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

	userID, err := h.GetUserIDFromCookie(r)
	if err != nil {
		h.LogError(r, err, "Failed to get user ID from cookie")
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

	userID, err := h.GetUserIDFromCookie(r)
	if err != nil {
		h.LogError(r, err, "Failed to get user ID from cookie")
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

	userID, err := h.GetUserIDFromCookie(r)
	if err != nil {
		h.LogError(r, err, "Failed to get user ID from cookie")
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

type CreateClientRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type CreateClientResponse struct {
	Success bool        `json:"success"`
	Client  *user.Client `json:"client,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (h *AuthHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "CreateClient request")

	if r.Method != http.MethodPost {
		h.LogError(r, nil, "Method not allowed")
		h.SendError(w, r, errors.NewAppError(http.StatusMethodNotAllowed, "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.GetUserIDFromCookie(r)
	if err != nil {
		h.LogError(r, err, "Failed to get user ID from cookie")
		h.SendError(w, r, err, http.StatusBadRequest)
		return
	}

	var req CreateClientRequest
	err = h.DecodeJSONBody(r, &req)
	if err != nil {
		h.LogError(r, err, "Invalid request body")
		h.SendError(w, r, err, http.StatusBadRequest)
		return
	}

	if req.ClientID == "" || req.ClientSecret == "" {
		h.LogError(r, nil, "Client ID and Client Secret are required")
		h.SendError(w, r, errors.NewAppError(http.StatusBadRequest, "Client ID and Client Secret are required"), http.StatusBadRequest)
		return
	}

	client, err := h.userService.CreateClient(userID, req.ClientID, req.ClientSecret)
	if err != nil {
		h.LogError(r, err, "Error creating client")
		
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == http.StatusConflict {
			h.SendError(w, r, err, appErr.Code)
			return
		}
		
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error creating client", err.Error()), http.StatusInternalServerError)
		return
	}

	response := CreateClientResponse{
		Success: true,
		Client:  client,
	}

	h.LogInfo(r, "Successfully created client")
	h.SendJSON(w, r, response, http.StatusCreated)
}

type GetClientsResponse struct {
	Success   bool        `json:"success"`
	Clients   []*user.Client `json:"clients,omitempty"`
	TotalCount int        `json:"total_count,omitempty"`
	Error     string      `json:"error,omitempty"`
}

func (h *AuthHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "GetClients request")

	if r.Method != http.MethodGet {
		h.LogError(r, nil, "Method not allowed")
		h.SendError(w, r, errors.NewAppError(http.StatusMethodNotAllowed, "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.GetUserIDFromCookie(r)
	if err != nil {
		h.LogError(r, err, "Failed to get user ID from cookie")
		h.SendError(w, r, err, http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 
	offset := 0 

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			h.LogError(r, err, "Invalid limit parameter")
			h.SendError(w, r, errors.NewAppError(http.StatusBadRequest, "Invalid limit parameter"), http.StatusBadRequest)
			return
		}
		limit = parsedLimit
	}

	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			h.LogError(r, err, "Invalid offset parameter")
			h.SendError(w, r, errors.NewAppError(http.StatusBadRequest, "Invalid offset parameter"), http.StatusBadRequest)
			return
		}
		offset = parsedOffset
	}

	clients, err := h.userService.GetClientsByUserIDWithPagination(userID, limit, offset)
	if err != nil {
		h.LogError(r, err, "Error getting clients")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error getting clients", err.Error()), http.StatusInternalServerError)
		return
	}

	totalCount, err := h.userService.GetClientsCountByUserID(userID)
	if err != nil {
		h.LogError(r, err, "Error getting clients count")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error getting clients count", err.Error()), http.StatusInternalServerError)
		return
	}

	response := GetClientsResponse{
		Success:    true,
		Clients:    clients,
		TotalCount: totalCount,
	}

	h.LogInfo(r, "Successfully retrieved clients")
	h.SendJSON(w, r, response, http.StatusOK)
}
