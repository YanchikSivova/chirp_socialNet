package service

import (
	"accountService/models"
	"accountService/repository"
	"context"
	"errors"
	"github.com/google/uuid"
)

type UsersService struct {
	repo *repository.UsersRepository
}

func NewUsersService(r *repository.UsersRepository) *UsersService {
	return &UsersService{repo: r}
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
