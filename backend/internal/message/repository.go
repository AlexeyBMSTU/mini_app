package message

import (
	"database/sql"
	"log"
	"time"
)

type SQLMessageRepository struct {
	db *sql.DB
}

func NewSQLMessageRepository(db *sql.DB) *SQLMessageRepository {
	return &SQLMessageRepository{
		db: db,
	}
}

func (r *SQLMessageRepository) CreateMessage(message *Message) error {
	query := `
		INSERT INTO messages (id, client_id, client_secret, message, name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(query,
		message.ID,
		message.ClientID,
		message.ClientSecret,
		message.Message,
		message.Name,
		message.IsActive,
		message.CreatedAt,
		message.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating message: %v", err)
		return err
	}

	return nil
}

func (r *SQLMessageRepository) GetMessageByID(id string) (*Message, error) {
	query := `
		SELECT id, client_id, client_secret, message, name, is_active, created_at, updated_at
		FROM messages
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	msg := &Message{}
	err := row.Scan(
		&msg.ID,
		&msg.ClientID,
		&msg.ClientSecret,
		&msg.Message,
		&msg.Name,
		&msg.IsActive,
		&msg.CreatedAt,
		&msg.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error getting message by ID: %v", err)
		return nil, err
	}

	return msg, nil
}

func (r *SQLMessageRepository) GetMessagesByClientID(clientID string) ([]*Message, error) {
	query := `
		SELECT id, client_id, client_secret, message, name, is_active, created_at, updated_at
		FROM messages
		WHERE client_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, clientID)
	if err != nil {
		log.Printf("Error getting messages by client ID: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		err := rows.Scan(
			&msg.ID,
			&msg.ClientID,
			&msg.ClientSecret,
			&msg.Message,
			&msg.Name,
			&msg.IsActive,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning message: %v", err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *SQLMessageRepository) GetMessageByClientCredentials(clientID, clientSecret string) (*Message, error) {
	query := `
		SELECT id, client_id, client_secret, message, name, is_active, created_at, updated_at
		FROM messages
		WHERE client_id = $1 AND client_secret = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	row := r.db.QueryRow(query, clientID, clientSecret)

	msg := &Message{}
	err := row.Scan(
		&msg.ID,
		&msg.ClientID,
		&msg.ClientSecret,
		&msg.Message,
		&msg.Name,
		&msg.IsActive,
		&msg.CreatedAt,
		&msg.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error getting message by client credentials: %v", err)
		return nil, err
	}

	return msg, nil
}

func (r *SQLMessageRepository) UpdateMessage(message *Message) error {
	query := `
		UPDATE messages
		SET client_id = $2, client_secret = $3, message = $4, name = $5, is_active = $6, updated_at = $7
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		message.ID,
		message.ClientID,
		message.ClientSecret,
		message.Message,
		message.Name,
		message.IsActive,
		message.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error updating message: %v", err)
		return err
	}

	return nil
}

func (r *SQLMessageRepository) DeleteMessage(id string) error {
	query := `DELETE FROM messages WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting message: %v", err)
		return err
	}

	return nil
}

func (r *SQLMessageRepository) CountMessagesByClientID(clientID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM messages
		WHERE client_id = $1
	`

	var count int
	err := r.db.QueryRow(query, clientID).Scan(&count)
	if err != nil {
		log.Printf("Error counting messages by client ID: %v", err)
		return 0, err
	}

	return count, nil
}

func (r *SQLMessageRepository) CountMessagesByClientIDInTimeRange(clientID string, duration time.Duration) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM messages
		WHERE client_id = $1 AND created_at > $2
	`

	since := time.Now().Add(-duration)
	var count int
	err := r.db.QueryRow(query, clientID, since).Scan(&count)
	if err != nil {
		log.Printf("Error counting messages by client ID in time range: %v", err)
		return 0, err
	}

	return count, nil
}

func (r *SQLMessageRepository) CreateMessagesTable() error {
	messagesTable := `
		CREATE TABLE IF NOT EXISTS messages (
			id UUID PRIMARY KEY,
			client_id VARCHAR(255) NOT NULL,
			client_secret VARCHAR(255) NOT NULL,
			message TEXT NOT NULL,
			name VARCHAR(255) NOT NULL,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`

	_, err := r.db.Exec(messagesTable)
	if err != nil {
		log.Printf("Error creating messages table: %v", err)
		return err
	}

	indexQuery := `
		CREATE INDEX IF NOT EXISTS idx_messages_client_id ON messages(client_id);
	`

	_, err = r.db.Exec(indexQuery)
	if err != nil {
		log.Printf("Error creating index: %v", err)
		return err
	}

	return nil
}