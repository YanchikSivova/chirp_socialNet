package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"notificationService/models"
)

type NotificationRepository struct {
	DB *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{
		DB: db,
	}
}
func (r *NotificationRepository) CheckEventProcessed(ctx context.Context, tx pgx.Tx, eventID uuid.UUID) (bool, error) {
	var exists bool
	err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM processed_events WHERE event_id = $1)`, eventID).Scan(&exists)
	return exists, err
}

func (r *NotificationRepository) ProcessEvent(ctx context.Context, tx pgx.Tx, eventID uuid.UUID) error {
	_, err := tx.Exec(ctx, `INSERT INTO processed_events(event_id) VALUES ($1)`, eventID)
	return err
}

func (r *NotificationRepository) SaveNotification(ctx context.Context, tx pgx.Tx, not models.Notification) error {
	_, err := tx.Exec(ctx, `INSERT INTO notification (notification_id, profile_id, actor_id, type, entity_id, created_at) VALUES ($1, $2, $3, $4, $5, $6)`, not.NotificationID, not.ProfileID, not.ActorID, not.Type, not.EntityID, not.CreatedAt)
	return err
}

func (r *NotificationRepository) GetNotifications(ctx context.Context, tx pgx.Tx, profileID uuid.UUID) ([]models.Notification, error) {
	rows, err := tx.Query(ctx, `SELECT * FROM notification WHERE profile_id = $1 ORDER BY created_at DESC`, profileID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []models.Notification{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		err = rows.Scan(
			&notification.NotificationID,
			&notification.ProfileID,
			&notification.ActorID,
			&notification.Type,
			&notification.EntityID,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) DeleteNotificationByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	_, err := tx.Exec(ctx, `DELETE FROM notification WHERE notification_id = $1`, id)
	return err
}

func (r *NotificationRepository) DeleteAllNotifications(ctx context.Context, tx pgx.Tx, profileID uuid.UUID) error {
	_, err := tx.Exec(ctx, `DELETE FROM notification WHERE profile_id = $1`, profileID)
	return err
}

func (r *NotificationRepository) GetProfileID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (uuid.UUID, error) {
	var profileID uuid.UUID
	err := tx.QueryRow(ctx, `SELECT profile_id FROM notification WHERE notification_id = $1`, id).Scan(&profileID)
	return profileID, err
}
