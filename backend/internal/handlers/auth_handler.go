package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"mini-app-backend/internal/telegram"
	"mini-app-backend/internal/user"
)

type AuthHandler struct {
	userService *user.UserService
	botToken    string
	db          *sql.DB
}

func NewAuthHandler(userService *user.UserService, botToken string, db *sql.DB) *AuthHandler {
	return &AuthHandler{
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TelegramAuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Invalid request body: %v", err)
		SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.User == nil || req.InitData == "" {
		SendErrorResponse(w, "User and initData are required", http.StatusBadRequest)
		return
	}

	isValid, _, err := telegram.ValidateAndParseInitData(req.InitData, h.botToken)
	if err != nil {
		log.Printf("Error validating initData: %v", err)
		SendErrorResponse(w, "Error validating initData", http.StatusInternalServerError)
		return
	}

	if !isValid {
		SendErrorResponse(w, "Invalid initData", http.StatusUnauthorized)
		return
	}

	createdUser, err := h.userService.CreateOrUpdateUser(req.User)
	if err != nil {
		log.Printf("Error creating/updating user: %v", err)
		SendErrorResponse(w, "Error creating/updating user", http.StatusInternalServerError)
		return
	}

	token := "token_" + strconv.FormatInt(createdUser.ID, 10)

	response := TelegramAuthResponse{
		Success: true,
		User:    createdUser,
		Token:   token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		SendErrorResponse(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		SendErrorResponse(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		SendErrorResponse(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		SendErrorResponse(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) GetUserData(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		SendErrorResponse(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		SendErrorResponse(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	userData, err := h.userService.GetUserDataByUserID(userID)
	if err != nil {
		log.Printf("Error getting user data: %v", err)
		SendErrorResponse(w, "Error getting user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userData)
}

func (h *AuthHandler) SaveUserData(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		SendErrorResponse(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		SendErrorResponse(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	var data map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		SendErrorResponse(w, "Error marshaling data", http.StatusInternalServerError)
		return
	}

	_, err = h.userService.SaveUserData(userID, string(dataJSON))
	if err != nil {
		log.Printf("Error saving user data: %v", err)
		SendErrorResponse(w, "Error saving user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func SendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
