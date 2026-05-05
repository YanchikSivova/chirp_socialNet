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

func (s *AuthService) ResetVerification(email string) (string, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	var cred *models.Credentials
	cred, err = s.repo.GetCredentialsByEmail(ctx, tx, email)
	if err != nil {
		return "", err
	}
	if cred.Status != "pending" {
		return "", errors.New("Email already active")
	}

	err = s.repo.DeleteVerification(ctx, tx, cred.CredentialsID)
	if err != nil {
		return "", err
	}

	newVerif := models.Verification{
		VerificationID: uuid.New(),
		CredentialsID:  cred.CredentialsID,
		Code:           utils.GenerateCode(),
		ExpiresAt:      time.Now().Add(5 * time.Minute),
	}

	err = s.repo.CreateVerification(ctx, tx, newVerif)
	if err != nil {
		return "", err
	}
	err = mail.SendVerificationEmail(email, newVerif.Code)
	if err != nil {
		return "", err
	}
	token, err := utils.GenerateRegisterToken(email)
	if err != nil {
		return "", err
	}
	return token, tx.Commit(ctx)
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

func (s *AuthService) Login(email, password string) (string, string, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return "", "", err
	}
	defer tx.Rollback(ctx)

	cred, err := s.repo.GetCredentialsByEmail(ctx, tx, email)
	if err != nil {
		return "", "", err
	}
	if ok := utils.CheckPasswordHash(password, cred.HashedPassword); !ok {
		return "", "", errors.New("Invalid password")
	}
	if cred.Status != "active" {
		return "", "", errors.New("Email not verified")
	}
	if cred.ProfileID == nil {
		id := uuid.New()
		err = s.repo.CreateProfile(ctx, tx, id)
		if err != nil {
			return "", "", err
		}
		err = s.repo.AttachProfile(ctx, tx, cred.CredentialsID, id)
		if err != nil {
			return "", "", err
		}
		cred.ProfileID = &id
	}
	access, err := utils.GenerateAccessToken(*cred.ProfileID)
	if err != nil {
		return "", "", err
	}
	refreshT := models.RefreshToken{
		RefreshTokenID: uuid.New(),
		ProfileID:      *cred.ProfileID,
	}
	expT, rT, err := utils.GenerateRefreshToken(*cred.ProfileID, refreshT.RefreshTokenID)
	if err != nil {
		return "", "", err
	}
	refreshT.RefreshToken = rT
	refreshT.ExpiresAt = expT
	err = s.repo.SaveRefreshToken(ctx, tx, refreshT.ProfileID, refreshT.RefreshTokenID, refreshT.RefreshToken, refreshT.ExpiresAt)
	if err != nil {
		return "", "", err
	}
	return access, refreshT.RefreshToken, nil
}
