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
