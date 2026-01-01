package message

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID           string    `json:"id" db:"id"`
	ClientID     string    `json:"client_id" db:"client_id"`
	ClientSecret string    `json:"client_secret" db:"client_secret"`
	Message      string    `json:"message" db:"message"`
	Name         string    `json:"name" db:"name"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type MessageRepository interface {
	CreateMessage(message *Message) error
	GetMessageByID(id string) (*Message, error)
	GetMessagesByClientID(clientID string) ([]*Message, error)
	GetMessageByClientCredentials(clientID, clientSecret string) (*Message, error)
	UpdateMessage(message *Message) error
	DeleteMessage(id string) error
	CountMessagesByClientID(clientID string) (int, error)
	CountMessagesByClientIDInTimeRange(clientID string, duration time.Duration) (int, error)
}

type MessageService struct {
	repo MessageRepository
}

func NewMessageService(repo MessageRepository) *MessageService {
	return &MessageService{
		repo: repo,
	}
}

func (s *MessageService) CreateMessage(clientID, clientSecret, message, name string) (*Message, error) {
	msg := &Message{
		ID:           uuid.New().String(),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Message:      message,
		Name:         name,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := s.repo.CreateMessage(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *MessageService) GetMessageByID(id string) (*Message, error) {
	return s.repo.GetMessageByID(id)
}

func (s *MessageService) GetMessagesByClientID(clientID string) ([]*Message, error) {
	return s.repo.GetMessagesByClientID(clientID)
}

func (s *MessageService) GetMessageByClientCredentials(clientID, clientSecret string) (*Message, error) {
	return s.repo.GetMessageByClientCredentials(clientID, clientSecret)
}

func (s *MessageService) UpdateMessage(id string, message, name string, isActive bool) (*Message, error) {
	msg, err := s.repo.GetMessageByID(id)
	if err != nil {
		return nil, err
	}

	if msg == nil {
		return nil, nil
	}

	msg.Message = message
	msg.Name = name
	msg.IsActive = isActive
	msg.UpdatedAt = time.Now()

	err = s.repo.UpdateMessage(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *MessageService) DeleteMessage(id string) error {
	return s.repo.DeleteMessage(id)
}

func (s *MessageService) CountMessagesByClientID(clientID string) (int, error) {
	return s.repo.CountMessagesByClientID(clientID)
}

func (s *MessageService) CountMessagesByClientIDInTimeRange(clientID string, duration time.Duration) (int, error) {
	return s.repo.CountMessagesByClientIDInTimeRange(clientID, duration)
}