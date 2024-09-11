package postgresql

import (
	"MESSAGEAPI/src/internal/messages/models"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

type IMessageRepository interface {
	GetAllSentMessages() ([]models.Message, error)
	GetUnSentMessages() ([]models.Message, error)
	BeginTransaction() (pgx.Tx, error)
	MarkMessageAsSent(messageID string, tx pgx.Tx) error
}

type MessageRepository struct {
	dbPool *pgxpool.Pool
}

func NewMessageRepository(dbPool *pgxpool.Pool) IMessageRepository {
	return &MessageRepository{dbPool: dbPool}
}

func (mr *MessageRepository) GetAllSentMessages() ([]models.Message, error) {
	ctx := context.Background()
	rows, err := mr.dbPool.Query(ctx, "SELECT * FROM messages WHERE sent = true")
	if err != nil {
		log.Error("Error while getting sent messages: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.Content, &m.RecipientPhoneNumber, &m.Sent); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func (mr *MessageRepository) GetUnSentMessages() ([]models.Message, error) {
	ctx := context.Background()
	rows, err := mr.dbPool.Query(ctx, "SELECT * FROM messages WHERE sent = false")
	if err != nil {
		log.Error("Error while getting unsent messages: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.Content, &m.RecipientPhoneNumber, &m.Sent); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func (mr *MessageRepository) BeginTransaction() (pgx.Tx, error) {
	ctx := context.Background()
	tx, err := mr.dbPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (mr *MessageRepository) MarkMessageAsSent(messageID string, tx pgx.Tx) error {
	_, err := tx.Exec(context.Background(), "UPDATE messages SET sent = true WHERE id = $1", messageID)
	if err != nil {
		log.Errorf("Error while marking message as sent: %w", err)
		return err
	}
	return err
}
