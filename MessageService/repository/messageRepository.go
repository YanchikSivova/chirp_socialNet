package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"messageService/models"
	"time"
)

type MessageRepository struct {
	DB *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{
		DB: db,
	}
}

func (r *MessageRepository) GetConversationID(ctx context.Context, tx pgx.Tx, meProfileID, targetProfileID uuid.UUID) (uuid.UUID, error) {
	var conversationID uuid.UUID
	err := tx.QueryRow(ctx, `SELECT conversation_id FROM conversation_member WHERE profile_id in ($1, $2) group by conversation_id having count(*) = 2`, meProfileID, targetProfileID).Scan(&conversationID)
	return conversationID, err
}

func (r *MessageRepository) CreateConversation(ctx context.Context, tx pgx.Tx, meProfileID, targetProfileID, conversationID uuid.UUID) error {
	_, err := tx.Exec(ctx, `INSERT INTO conversation (conversation_id) VALUES ($1)`, conversationID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO conversation_member (conversation_member_id, conversation_id, profile_id) values ($1, $2, $3)`, uuid.New(), conversationID, meProfileID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO conversation_member (conversation_member_id, conversation_id, profile_id) values ($1, $2, $3)`, uuid.New(), conversationID, targetProfileID)
	return err
}

func (r *MessageRepository) GetMessages(ctx context.Context, tx pgx.Tx, conversationID uuid.UUID, cursor *time.Time) ([]models.Message, error) {
	var rows pgx.Rows
	var err error
	if cursor == nil {
		rows, err = tx.Query(ctx, `SELECT message_id, sender_id, type, content, created_at FROM message WHERE conversation_id = $1 ORDER BY created_at DESC limit 20`, conversationID)
	} else {
		rows, err = tx.Query(ctx, `SELECT message_id, sender_id, type, content, created_at FROM message WHERE conversation_id = $1 AND created_at <$2 ORDER BY created_at DESC limit 20`, conversationID, *cursor)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return []models.Message{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	var messages []models.Message
	for rows.Next() {
		var message models.Message
		err := rows.Scan(
			&message.MessageID,
			&message.SenderID,
			&message.Type,
			&message.Content,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageRepository) GetTargetProfileID(ctx context.Context, tx pgx.Tx, meProfileID, conversationID uuid.UUID) (*uuid.UUID, error) {
	var targetProfileID uuid.UUID
	err := tx.QueryRow(ctx, `SELECT profile_id from conversation_member where conversation_id=$1 and profile_id != $2`, conversationID, meProfileID).Scan(&targetProfileID)
	log.Println(conversationID, meProfileID)
	return &targetProfileID, err
}

func (r *MessageRepository) GetConversationList(ctx context.Context, tx pgx.Tx, meProfileID uuid.UUID) ([]models.ConversationPreview, error) {
	var conversations []models.ConversationPreview
	rows, err := tx.Query(ctx, `SELECT conversation_id, unread_count FROM conversation_member WHERE profile_id = $1`, meProfileID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []models.ConversationPreview{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	successCount := 0
	for rows.Next() {
		var conversation models.ConversationPreview
		err := rows.Scan(&conversation.ConversationID, &conversation.UnreadCount)
		if err != nil {
			continue
		}
		conversations = append(conversations, conversation)
		successCount++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if successCount == 0 {
		return nil, errors.New("no conversation found")
	}
	return conversations, nil
}
