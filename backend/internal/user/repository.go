package user

import (
	"database/sql"
	"encoding/json"
	"log"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{
		db: db,
	}
}

func (r *SQLRepository) CreateUser(user *User) error {
	query := `
		INSERT INTO users (id, first_name, last_name, username, language_code, is_premium, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := r.db.Exec(query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Username,
		user.LanguageCode,
		user.IsPremium,
		user.CreatedAt,
		user.UpdatedAt,
	)
	
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}
	
	return nil
}

func (r *SQLRepository) GetUserByID(id int64) (*User, error) {
	query := `
		SELECT id, first_name, last_name, username, language_code, is_premium, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	
	row := r.db.QueryRow(query, id)
	
	user := &User{}
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.LanguageCode,
		&user.IsPremium,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error getting user by ID: %v", err)
		return nil, err
	}
	
	return user, nil
}

func (r *SQLRepository) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET first_name = $2, last_name = $3, username = $4, language_code = $5, is_premium = $6, updated_at = $7
		WHERE id = $1
	`
	
	_, err := r.db.Exec(query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Username,
		user.LanguageCode,
		user.IsPremium,
		user.UpdatedAt,
	)
	
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return err
	}
	
	return nil
}

func (r *SQLRepository) GetUserByTelegramID(telegramID int64) (*User, error) {
	return r.GetUserByID(telegramID)
}

func (r *SQLRepository) CreateUserData(userData *UserData) error {
	query := `
		INSERT INTO user_data (user_id, data, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	
	var id int64
	err := r.db.QueryRow(query,
		userData.UserID,
		userData.Data,
		userData.CreatedAt,
		userData.UpdatedAt,
	).Scan(&id)
	
	if err != nil {
		log.Printf("Error creating user data: %v", err)
		return err
	}
	
	userData.ID = id
	return nil
}

func (r *SQLRepository) GetUserDataByUserID(userID int64) (*UserData, error) {
	query := `
		SELECT id, user_id, data, created_at, updated_at
		FROM user_data
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	
	row := r.db.QueryRow(query, userID)
	
	userData := &UserData{}
	err := row.Scan(
		&userData.ID,
		&userData.UserID,
		&userData.Data,
		&userData.CreatedAt,
		&userData.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error getting user data by user ID: %v", err)
		return nil, err
	}
	
	return userData, nil
}

func (r *SQLRepository) UpdateUserData(userData *UserData) error {
	query := `
		UPDATE user_data
		SET data = $2, updated_at = $3
		WHERE id = $1
	`
	
	_, err := r.db.Exec(query,
		userData.ID,
		userData.Data,
		userData.UpdatedAt,
	)
	
	if err != nil {
		log.Printf("Error updating user data: %v", err)
		return err
	}
	
	return nil
}

func (r *SQLRepository) CreateTables() error {
	usersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255),
			username VARCHAR(255),
			language_code VARCHAR(10),
			is_premium BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`
	
	_, err := r.db.Exec(usersTable)
	if err != nil {
		log.Printf("Error creating users table: %v", err)
		return err
	}
	
	userDataTable := `
		CREATE TABLE IF NOT EXISTS user_data (
			id SERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			data TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`
	
	_, err = r.db.Exec(userDataTable)
	if err != nil {
		log.Printf("Error creating user_data table: %v", err)
		return err
	}
	
	indexQuery := `
		CREATE INDEX IF NOT EXISTS idx_user_data_user_id ON user_data(user_id);
	`
	
	_, err = r.db.Exec(indexQuery)
	if err != nil {
		log.Printf("Error creating index: %v", err)
		return err
	}
	
	return nil
}

func (r *SQLRepository) GetUsersWithPagination(limit, offset int) ([]*User, error) {
	query := `
		SELECT id, first_name, last_name, username, language_code, is_premium, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		log.Printf("Error getting users with pagination: %v", err)
		return nil, err
	}
	defer rows.Close()
	
	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.LanguageCode,
			&user.IsPremium,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			return nil, err
		}
		users = append(users, user)
	}
	
	return users, nil
}

func (r *SQLRepository) GetUsersCount() (int, error) {
	query := `SELECT COUNT(*) FROM users`
	
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Error getting users count: %v", err)
		return 0, err
	}
	
	return count, nil
}

func (r *SQLRepository) GetUserDataJSON(userID int64) (map[string]interface{}, error) {
	userData, err := r.GetUserDataByUserID(userID)
	if err != nil {
		return nil, err
	}
	
	if userData == nil {
		return make(map[string]interface{}), nil
	}
	
	var data map[string]interface{}
	err = json.Unmarshal([]byte(userData.Data), &data)
	if err != nil {
		log.Printf("Error unmarshaling user data: %v", err)
		return nil, err
	}
	
	return data, nil
}