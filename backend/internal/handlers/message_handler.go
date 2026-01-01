package handlers

import (
	"database/sql"
	"mini-app-backend/internal/errors"
	"mini-app-backend/internal/message"
	"net/http"
	"strings"
)

type MessageHandler struct {
	*BaseHandler
	messageService *message.MessageService
	db             *sql.DB
}

func NewMessageHandler(messageService *message.MessageService, db *sql.DB) *MessageHandler {
	return &MessageHandler{
		BaseHandler:    NewBaseHandler(),
		messageService: messageService,
		db:             db,
	}
}

func (h *MessageHandler) GetIDParam(r *http.Request) (string, error) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		idStr := parts[len(parts)-1]
		if idStr == "" {
			return "", errors.NewAppError(http.StatusBadRequest, "ID not found in URL")
		}
		return idStr, nil
	}
	return "", errors.NewAppError(http.StatusBadRequest, "ID not found in URL")
}

type CreateMessageRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Message      string `json:"message"`
	Name         string `json:"name"`
}

type UpdateMessageRequest struct {
	Message  string `json:"message"`
	Name     string `json:"name"`
	IsActive *bool  `json:"is_active"`
}

type MessageResponse struct {
	Success  bool              `json:"success"`
	Message  *message.Message  `json:"message,omitempty"`
	Messages []*message.Message `json:"messages,omitempty"`
	Error    string            `json:"error,omitempty"`
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "CreateMessage request")

	if r.Method != http.MethodPost {
		h.LogError(r, nil, "Method not allowed")
		h.SendError(w, r, errors.NewAppError(http.StatusMethodNotAllowed, "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var req CreateMessageRequest
	err := h.DecodeJSONBody(r, &req)
	if err != nil {
		h.LogError(r, err, "Invalid request body")
		h.SendError(w, r, err, http.StatusBadRequest)
		return
	}

	if req.ClientID == "" || req.ClientSecret == "" || req.Message == "" || req.Name == "" {
		h.LogError(r, nil, "client_id, client_secret, message and name are required")
		h.SendError(w, r, errors.NewAppError(http.StatusBadRequest, "client_id, client_secret, message and name are required"), http.StatusBadRequest)
		return
	}

	message, err := h.messageService.CreateMessage(req.ClientID, req.ClientSecret, req.Message, req.Name)
	if err != nil {
		h.LogError(r, err, "error creating message")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "error creating message", err.Error()), http.StatusInternalServerError)
		return
	}

	response := MessageResponse{
		Success: true,
		Message: message,
	}

	h.LogInfo(r, "successfully created message")
	h.SendJSON(w, r, response, http.StatusCreated)
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "GetMessages request")

	if r.Method != http.MethodGet {
		h.LogError(r, nil, "method not allowed")
		h.SendError(w, r, errors.NewAppError(http.StatusMethodNotAllowed, "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		h.LogError(r, nil, "client_id is required")
		h.SendError(w, r, errors.NewAppError(http.StatusBadRequest, "client_id is required"), http.StatusBadRequest)
		return
	}

	messages, err := h.messageService.GetMessagesByClientID(clientID)
	if err != nil {
		h.LogError(r, err, "Error getting messages")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error getting messages", err.Error()), http.StatusInternalServerError)
		return
	}

	response := MessageResponse{
		Success:  true,
		Messages: messages,
	}

	h.LogInfo(r, "Successfully retrieved messages")
	h.SendJSON(w, r, response, http.StatusOK)
}

func (h *MessageHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "GetMessage request")

	if r.Method != http.MethodGet {
		h.LogError(r, nil, "Method not allowed")
		h.SendError(w, r, errors.NewAppError(http.StatusMethodNotAllowed, "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	messageID := r.URL.Query().Get("message_id")
	if messageID == "" {
		h.LogError(r, nil, "invalid message_id")
		h.SendError(w, r, errors.NewAppError(http.StatusBadRequest, "invalid message_id"), http.StatusBadRequest)
		return
	}

	message, err := h.messageService.GetMessageByID(messageID)
	if err != nil {
		h.LogError(r, err, "Error getting message")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error getting message", err.Error()), http.StatusInternalServerError)
		return
	}

	if message == nil {
		h.LogError(r, nil, "Message not found")
		h.SendError(w, r, errors.NewAppError(http.StatusNotFound, "Message not found"), http.StatusNotFound)
		return
	}

	response := MessageResponse{
		Success: true,
		Message: message,
	}

	h.LogInfo(r, "Successfully retrieved message")
	h.SendJSON(w, r, response, http.StatusOK)
}

func (h *MessageHandler) GetMessageByCredentials(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "GetMessageByCredentials request")

	if r.Method != http.MethodGet {
		h.LogError(r, nil, "Method not allowed")
		h.SendError(w, r, errors.NewAppError(http.StatusMethodNotAllowed, "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}

	err := h.DecodeJSONBody(r, &req)
	if err != nil {
		h.LogError(r, err, "Invalid request body")
		h.SendError(w, r, err, http.StatusBadRequest)
		return
	}

	if req.ClientID == "" || req.ClientSecret == "" {
		h.LogError(r, nil, "client_id and client_secret are required")
		h.SendError(w, r, errors.NewAppError(http.StatusBadRequest, "ClientID and ClientSecret are required"), http.StatusBadRequest)
		return
	}

	message, err := h.messageService.GetMessageByClientCredentials(req.ClientID, req.ClientSecret)
	if err != nil {
		h.LogError(r, err, "Error getting message")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error getting message", err.Error()), http.StatusInternalServerError)
		return
	}

	if message == nil {
		h.LogError(r, nil, "Message not found")
		h.SendError(w, r, errors.NewAppError(http.StatusNotFound, "Message not found"), http.StatusNotFound)
		return
	}

	response := MessageResponse{
		Success: true,
		Message: message,
	}

	h.LogInfo(r, "Successfully retrieved message by credentials")
	h.SendJSON(w, r, response, http.StatusOK)
}

func (h *MessageHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "UpdateMessage request")

	if r.Method != http.MethodPut {
		h.LogError(r, nil, "Method not allowed")
		h.SendError(w, r, errors.NewAppError(http.StatusMethodNotAllowed, "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	messageID := r.URL.Query().Get("message_id")
	if messageID == "" {
		h.LogError(r, nil, "invalid message_id")
		h.SendError(w, r, errors.NewAppError(http.StatusBadRequest, "invalid message_id"), http.StatusBadRequest)
		return
	}

	var req UpdateMessageRequest
	err := h.DecodeJSONBody(r, &req)
	if err != nil {
		h.LogError(r, err, "invalid request body")
		h.SendError(w, r, err, http.StatusBadRequest)
		return
	}

	existingMessage, err := h.messageService.GetMessageByID(messageID)
	if err != nil {
		h.LogError(r, err, "Error getting message")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error getting message", err.Error()), http.StatusInternalServerError)
		return
	}

	if existingMessage == nil {
		h.LogError(r, nil, "Message not found")
		h.SendError(w, r, errors.NewAppError(http.StatusNotFound, "Message not found"), http.StatusNotFound)
		return
	}

	messageText := existingMessage.Message
	name := existingMessage.Name
	isActive := existingMessage.IsActive

	if req.Message != "" {
		messageText = req.Message
	}

	if req.Name != "" {
		name = req.Name
	}

	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	updatedMessage, err := h.messageService.UpdateMessage(messageID, messageText, name, isActive)
	if err != nil {
		h.LogError(r, err, "Error updating message")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error updating message", err.Error()), http.StatusInternalServerError)
		return
	}

	response := MessageResponse{
		Success: true,
		Message: updatedMessage,
	}

	h.LogInfo(r, "Successfully updated message")
	h.SendJSON(w, r, response, http.StatusOK)
}

func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(r, "DeleteMessage request")

	if r.Method != http.MethodDelete {
		h.LogError(r, nil, "Method not allowed")
		h.SendError(w, r, errors.NewAppError(http.StatusMethodNotAllowed, "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	messageID := r.URL.Query().Get("message_id")
	if messageID == "" {
		h.LogError(r, nil, "invalid message_id")
		h.SendError(w, r, errors.NewAppError(http.StatusBadRequest, "invalid message_id"), http.StatusBadRequest)
		return
	}

	existingMessage, err := h.messageService.GetMessageByID(messageID)
	if err != nil {
		h.LogError(r, err, "Error getting message")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error getting message", err.Error()), http.StatusInternalServerError)
		return
	}

	if existingMessage == nil {
		h.LogError(r, nil, "Message not found")
		h.SendError(w, r, errors.NewAppError(http.StatusNotFound, "Message not found"), http.StatusNotFound)
		return
	}

	err = h.messageService.DeleteMessage(messageID)
	if err != nil {
		h.LogError(r, err, "Error deleting message")
		h.SendError(w, r, errors.NewAppErrorWithDetails(http.StatusInternalServerError, "Error deleting message", err.Error()), http.StatusInternalServerError)
		return
	}

	h.LogInfo(r, "Successfully deleted message")
	h.SendJSON(w, r, map[string]bool{"success": true}, http.StatusOK)
}