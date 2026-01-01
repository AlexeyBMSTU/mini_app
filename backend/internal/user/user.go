package user

import (
	"mini-app-backend/internal/errors"
	"net/http"
	"time"
)

type User struct {
	ID           int64     `json:"id" db:"id"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Username     string    `json:"username" db:"username"`
	LanguageCode string    `json:"language_code" db:"language_code"`
	IsPremium    bool      `json:"is_premium" db:"is_premium"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserData struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Data      string    `json:"data" db:"data"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Client struct {
	ID           int64     `json:"id" db:"id"`
	ClientID     string    `json:"client_id" db:"client_id"`
	ClientSecret string    `json:"client_secret" db:"client_secret"`
	UserID       int64     `json:"user_id" db:"user_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByID(id int64) (*User, error)
	UpdateUser(user *User) error
	GetUserByTelegramID(telegramID int64) (*User, error)

	CreateUserData(userData *UserData) error
	GetUserDataByUserID(userID int64) (*UserData, error)
	UpdateUserData(userData *UserData) error
	
	CreateClient(client *Client) error
	GetClientByUserID(userID int64) (*Client, error)
	GetClientByCredentials(clientID, clientSecret string) (*Client, error)
	GetClientByID(clientID string) (*Client, error)
	GetClientBySecret(clientSecret string) (*Client, error)
	GetClientsWithPagination(limit, offset int) ([]*Client, error)
	GetClientsCount() (int, error)
	GetClientsByUserIDWithPagination(userID int64, limit, offset int) ([]*Client, error)
	GetClientsCountByUserID(userID int64) (int, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateOrUpdateUser(telegramUser *User) (*User, error) {
	existingUser, err := s.repo.GetUserByTelegramID(telegramUser.ID)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		existingUser.FirstName = telegramUser.FirstName
		existingUser.LastName = telegramUser.LastName
		existingUser.Username = telegramUser.Username
		existingUser.LanguageCode = telegramUser.LanguageCode
		existingUser.IsPremium = telegramUser.IsPremium
		existingUser.UpdatedAt = time.Now()

		err := s.repo.UpdateUser(existingUser)
		if err != nil {
			return nil, err
		}

		return existingUser, nil
	}

	telegramUser.CreatedAt = time.Now()
	telegramUser.UpdatedAt = time.Now()

	err = s.repo.CreateUser(telegramUser)
	if err != nil {
		return nil, err
	}

	return telegramUser, nil
}

func (s *UserService) GetUserByID(id int64) (*User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) GetUserDataByUserID(userID int64) (*UserData, error) {
	return s.repo.GetUserDataByUserID(userID)
}

func (s *UserService) SaveUserData(userID int64, data string) (*UserData, error) {
	existingData, err := s.repo.GetUserDataByUserID(userID)
	if err != nil {
		return nil, err
	}

	if existingData != nil {
		existingData.Data = data
		existingData.UpdatedAt = time.Now()

		err := s.repo.UpdateUserData(existingData)
		if err != nil {
			return nil, err
		}

		return existingData, nil
	}

	userData := &UserData{
		UserID:    userID,
		Data:      data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.CreateUserData(userData)
	if err != nil {
		return nil, err
	}

	return userData, nil
}

func (s *UserService) CreateClient(userID int64, clientID, clientSecret string) (*Client, error) {
	existingClientByID, err := s.repo.GetClientByID(clientID)
	if err != nil {
		return nil, err
	}
	
	if existingClientByID != nil {
		return nil, errors.NewAppError(http.StatusConflict, "client with the same client_id already exists")
	}
	
	existingClientBySecret, err := s.repo.GetClientBySecret(clientSecret)
	if err != nil {
		return nil, err
	}
	
	if existingClientBySecret != nil {
		return nil, errors.NewAppError(http.StatusConflict, "client with the same client_secret already exists")
	}
	
	client := &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		UserID:       userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	
	err = s.repo.CreateClient(client)
	if err != nil {
		return nil, err
	}
	
	return client, nil
}

func (s *UserService) GetClientsWithPagination(limit, offset int) ([]*Client, error) {
	return s.repo.GetClientsWithPagination(limit, offset)
}

func (s *UserService) GetClientsCount() (int, error) {
	return s.repo.GetClientsCount()
}

func (s *UserService) GetClientsByUserIDWithPagination(userID int64, limit, offset int) ([]*Client, error) {
	return s.repo.GetClientsByUserIDWithPagination(userID, limit, offset)
}

func (s *UserService) GetClientsCountByUserID(userID int64) (int, error) {
	return s.repo.GetClientsCountByUserID(userID)
}
