package user

import (
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

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByID(id int64) (*User, error)
	UpdateUser(user *User) error
	GetUserByTelegramID(telegramID int64) (*User, error)

	CreateUserData(userData *UserData) error
	GetUserDataByUserID(userID int64) (*UserData, error)
	UpdateUserData(userData *UserData) error
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
