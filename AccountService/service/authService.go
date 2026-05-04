package service

import (
	"accountService/mail"
	"accountService/models"
	"accountService/repository"
	"accountService/utils"
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(r *repository.AuthRepository) *AuthService {
	return &AuthService{repo: r}
}

func (s *AuthService) Register(email, password string) (string, error) {
	ctx := context.Background()

	exists, err := s.repo.EmailExists(ctx, email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.New("Email already exists")
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}
	credentialsID := uuid.New()

	cred := models.Credentials{
		CredentialsID:  credentialsID,
		Email:          email,
		HashedPassword: hashed,
		Status:         "pending",
		CreatedAt:      time.Now(),
	}

	verif := models.Verification{
		VerificationID: uuid.New(),
		CredentialsID:  credentialsID,
		Code:           utils.GenerateCode(),
		ExpiresAt:      time.Now().Add(5 * time.Minute),
	}

	err = s.repo.RegisterTx(ctx, cred, verif)
	if err != nil {
		return "", err
	}

	err = mail.SendVerificationEmail(email, verif.Code)
	if err != nil {
		return "", err
	}
	token, err := utils.GenerateRegisterToken(email)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) VerifyEmail(email, code string) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	cred, err := s.repo.GetCredentialsByEmail(ctx, tx, email)
	if err != nil {
		return err
	}
	verif, err := s.repo.GetVerificationByCredentialsID(ctx, tx, cred.CredentialsID)
	if err != nil {
		return err
	}
	if verif.Code != code {
		return errors.New("Verification code does not match")
	}
	if time.Now().After(verif.ExpiresAt) {
		return errors.New("Code expired")
	}
	var profileID uuid.UUID
	if cred.ProfileID == nil {
		profileID = uuid.New()
		err = s.repo.CreateProfile(ctx, tx, profileID)
		if err != nil {
			return err
		}
		err = s.repo.AttachProfile(ctx, tx, cred.CredentialsID, profileID)
		if err != nil {
			return err
		}
	}

	err = s.repo.ActivateCredentials(ctx, tx, cred.CredentialsID)
	if err != nil {
		return err
	}
	err = s.repo.DeleteVerification(ctx, tx, verif.CredentialsID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
