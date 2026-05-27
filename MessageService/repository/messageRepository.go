package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	return &targetProfileID, err
}

func (r *MessageRepository) GetConversationList(ctx context.Context, tx pgx.Tx, meProfileID uuid.UUID) ([]models.ConversationPreview, error) {
	var conversations []models.ConversationPreview
	rows, err := tx.Query(ctx, `SELECT cm.conversation_id, cm.last_read_message_id FROM conversation_member cm join conversation c on c.conversation_id = cm.conversation_id   WHERE profile_id = $1 ORDER BY c.last_message_created_at desc`, meProfileID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []models.ConversationPreview{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var conversation models.ConversationPreview
		err := rows.Scan(&conversation.ConversationID, &conversation.LastReadMessageID)
		if err != nil {
			continue
		}
		conversations = append(conversations, conversation)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return conversations, nil
}

func (r *MessageRepository) GetUnreadCount(ctx context.Context, tx pgx.Tx, conversationID, profileID uuid.UUID, lastReadMessageID *uuid.UUID) (int, error) {
	var count int
	if lastReadMessageID == nil {
		err := tx.QueryRow(ctx, `
            SELECT COUNT(*)
            FROM message
            WHERE conversation_id = $1 AND sender_id != $2
        `, conversationID, profileID).Scan(&count)
		return count, err
	}
	err := tx.QueryRow(ctx, `
        SELECT COUNT(*)
        FROM message m
        WHERE m.conversation_id = $1 AND m.sender_id != $2
        AND m.created_at > (
            SELECT created_at FROM message WHERE message_id = $3
        )
    `, conversationID, profileID, *lastReadMessageID).Scan(&count)

	return count, err
}

func (r *MessageRepository) CreateMessage(ctx context.Context, tx pgx.Tx, messageID, conversationID, senderID uuid.UUID, content, messageType string) error {
	_, err := tx.Exec(ctx, `INSERT INTO message (message_id, conversation_id, sender_id, type, content) VALUES ($1, $2, $3, $4, $5)`, messageID, conversationID, senderID, messageType, content)
	return err
}

func (r *MessageRepository) UpdateLastMessageCreatedAt(ctx context.Context, tx pgx.Tx, conversationID uuid.UUID, createdAt time.Time) error {
	_, err := tx.Exec(ctx, `UPDATE conversation SET last_message_created_at = $1 WHERE conversation_id = $2`, createdAt, conversationID)
	return err
}
func (r *MessageRepository) MarkConversationRead(ctx context.Context, tx pgx.Tx, conversationID, profileID, lastMessageID uuid.UUID) error {
	_, err := tx.Exec(ctx, `UPDATE conversation_member SET last_read_message_id=$1 WHERE conversation_id = $2 AND profile_id=$3`, lastMessageID, conversationID, profileID)
	return err
}

func (r *MessageRepository) GetLastMessageID(ctx context.Context, tx pgx.Tx, conversationID uuid.UUID) (uuid.UUID, error) {
	var lastMessageID uuid.UUID
	err := tx.QueryRow(ctx, `SELECT message_id FROM message WHERE conversation_id = $1 ORDER BY created_at DESC LIMIT 1`, conversationID).Scan(&lastMessageID)
	return lastMessageID, err
}
