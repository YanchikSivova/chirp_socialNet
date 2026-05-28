package service

import (
	"accountService/kafka"
	"accountService/kafkaEvents"
	"accountService/models"
	"accountService/repository"
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
	"time"
)

type UsersService struct {
	repo     *repository.UsersRepository
	producer *kafka.NotificationProducer
}

func NewUsersService(r *repository.UsersRepository, pr *kafka.NotificationProducer) *UsersService {
	return &UsersService{repo: r, producer: pr}
}

func (s *UsersService) CheckUsername(username string) (bool, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)
	exists, err := s.repo.UsernameExists(ctx, tx, username)
	if err != nil {
		return false, err
	}
	return exists, tx.Commit(ctx)
}

func (s *UsersService) FillProfile(profileId uuid.UUID, name, username, avatar, description string) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	isCompleted, err := s.repo.CheckProfileCompleted(ctx, tx, profileId)
	if err != nil {
		return err
	}
	if isCompleted {
		return errors.New("profile already completed")
	}
	usernameExists, err := s.repo.UsernameExists(ctx, tx, username)
	if err != nil {
		return err
	}
	if usernameExists {
		return errors.New("username already exists")
	}

	//сам проставляет is_completed=true
	err = s.repo.FillProfile(ctx, tx, profileId, name, username, avatar, description)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *UsersService) UpdateProfile(profileId uuid.UUID, name, username, avatar, description string) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if name != "" {
		err := s.repo.UpdateName(ctx, tx, profileId, name)
		if err != nil {
			return err
		}
	}
	if username != "" {
		usernameExists, err := s.repo.UsernameExists(ctx, tx, username)
		if err != nil {
			return err
		}
		if usernameExists {
			return errors.New("username already exists")
		}
		err = s.repo.UpdateUsername(ctx, tx, profileId, username)
		if err != nil {
			return err
		}
	}
	if avatar != "" {
		err := s.repo.UpdateAvatar(ctx, tx, profileId, avatar)
		if err != nil {
			return err
		}
	}
	if description != "" {
		err := s.repo.UpdateDescription(ctx, tx, profileId, description)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (s *UsersService) GetProfile(profileID uuid.UUID) (*models.Profile, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	profile, err := s.repo.GetProfileByProfileID(ctx, tx, profileID)
	if err != nil {
		return nil, err
	}
	if !profile.IsCompleted {
		return nil, errors.New("profile is not completed")
	}
	return profile, tx.Commit(ctx)
}

func (s *UsersService) GetRelationship(meProfileId, profileId uuid.UUID) (string, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	isMeBlocked, err := s.repo.CheckIsBlocked(ctx, tx, meProfileId, profileId)
	if err != nil {
		return "", err
	}
	isAnotherBlocked, err := s.repo.CheckIsBlocked(ctx, tx, profileId, meProfileId)
	if err != nil {
		return "", err
	}
	isSubscribed, err := s.repo.CheckIsSubscribed(ctx, tx, meProfileId, profileId)
	if err != nil {
		return "", err
	}
	if isMeBlocked || isAnotherBlocked {
		return "blocked", tx.Commit(ctx)
	}
	if isSubscribed {
		return "subscribed", tx.Commit(ctx)
	}
	return "", tx.Commit(ctx)
}

func (s *UsersService) Follow(subscriberID, subscribedID uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = s.repo.Follow(ctx, tx, subscriberID, subscribedID)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	log.Println("Started to form kafka event")
	event := kafkaEvents.NotificationEvent{
		EventID:   uuid.New(),
		ProfileID: subscribedID,
		ActorID:   subscriberID,
		Type:      "subscription",
		EntityID:  subscriberID,
		CreatedAt: time.Now(),
	}
	err = s.producer.SendSubscriptionCreated(ctx, event)
	if err != nil {
		log.Println("Failed to send subscription created event", err)
	}
	return nil
}

func (s *UsersService) Unfollow(subscriberID, subscribedID uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = s.repo.Unfollow(ctx, tx, subscriberID, subscribedID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *UsersService) Block(bannedProfileId, profileId uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = s.repo.Block(ctx, tx, bannedProfileId, profileId)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *UsersService) Unblock(bannedProfileId, profileId uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = s.repo.Unblock(ctx, tx, bannedProfileId, profileId)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *UsersService) GetFollowers(profileId uuid.UUID, limit, offset int) ([]models.ProfileMinimum, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	followers, err := s.repo.GetFollowers(ctx, tx, profileId, limit, offset)
	if err != nil {
		return nil, err
	}
	return followers, tx.Commit(ctx)
}

func (s *UsersService) GetFollowings(profileId uuid.UUID, limit, offset int) ([]models.ProfileMinimum, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	followings, err := s.repo.GetFollowings(ctx, tx, profileId, limit, offset)
	if err != nil {
		return nil, err
	}
	return followings, tx.Commit(ctx)
}

func (s *UsersService) GetBlacklist(profileId uuid.UUID, limit, offset int) ([]models.ProfileMinimum, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	blacklist, err := s.repo.GetBlacklist(ctx, tx, profileId, limit, offset)
	if err != nil {
		return nil, err
	}
	return blacklist, tx.Commit(ctx)
}

func (s *UsersService) SearchByName(searchName string, limit, offset int) ([]models.ProfileMinimum, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	profiles, err := s.repo.SearchByName(ctx, tx, searchName, limit, offset)
	if err != nil {
		return nil, err
	}
	return profiles, tx.Commit(ctx)
}

func (s *UsersService) SearchByUsername(searchUsername string, limit, offset int) ([]models.ProfileMinimum, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	profiles, err := s.repo.SearchByUsername(ctx, tx, searchUsername, limit, offset)
	if err != nil {
		return nil, err
	}
	return profiles, tx.Commit(ctx)
}

func (s *UsersService) ProfileExists(profileID uuid.UUID) (bool, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)
	exists, err := s.repo.ProfileExists(ctx, tx, profileID)
	if err != nil {
		return false, err
	}
	return exists, tx.Commit(ctx)
}

func (s *UsersService) GetProfileMinimum(profileID uuid.UUID) (*models.ProfileMinimum, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	profileMinimum, err := s.repo.GetProfileMinimumById(ctx, tx, profileID)
	if err != nil {
		return nil, err
	}
	return profileMinimum, tx.Commit(ctx)
}

func (s *UsersService) GetFollowersID(profileID uuid.UUID) ([]uuid.UUID, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	followers, err := s.repo.GetFollowersID(ctx, tx, profileID)
	if err != nil {
		return nil, err
	}
	return followers, tx.Commit(ctx)
}
