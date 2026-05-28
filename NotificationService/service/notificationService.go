package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"notificationService/client"
	"notificationService/models"
	"notificationService/repository"
)

type NotificationService struct {
	repo          *repository.NotificationRepository
	accountClient *client.AccountClient
}

func NewNotificationService(repo *repository.NotificationRepository, ac *client.AccountClient) *NotificationService {
	return &NotificationService{repo, ac}
}

func (s *NotificationService) GetNotifications(profileID uuid.UUID) ([]models.NotificationsList, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	not, err := s.repo.GetNotifications(ctx, tx, profileID)
	if err != nil {
		return nil, err
	}
	if len(not) == 0 {
		return []models.NotificationsList{}, nil
	}
	actorIDs := make(map[uuid.UUID]struct{})
	for i := range not {
		actorIDs[not[i].ActorID] = struct{}{}
	}
	profiles, err := s.accountClient.BatchProfiles(actorIDs)
	if err != nil {
		return nil, err
	}
	actorProfiles := make(map[uuid.UUID]models.Profile)
	for i := range profiles {
		actorProfiles[profiles[i].ProfileID] = profiles[i]
	}
	successCount := 0
	var notes []models.NotificationsList
	for i := range not {
		actor, ok := actorProfiles[not[i].ActorID]
		if !ok {
			continue
		}
		note := models.NotificationsList{
			Notification: not[i],
			Actor:        actor,
		}
		notes = append(notes, note)
		successCount++
	}
	if successCount == 0 {
		return nil, errors.New("no notifications found")
	}
	return notes, tx.Commit(ctx)
}

func (s *NotificationService) DeleteNotification(profileID, notifID uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	ownerID, err := s.repo.GetProfileID(ctx, tx, notifID)
	if err != nil {
		return err
	}
	if ownerID != profileID {
		return errors.New("forbidden")
	}
	err = s.repo.DeleteNotificationByID(ctx, tx, notifID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *NotificationService) DeleteAllNotifications(profileID uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = s.repo.DeleteAllNotifications(ctx, tx, profileID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
