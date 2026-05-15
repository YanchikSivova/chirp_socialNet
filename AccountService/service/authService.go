package service

import (
	"accountService/mail"
	"accountService/models"
	"accountService/repository"
	"accountService/utils"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	exists, err := s.repo.EmailExists(ctx, tx, email)
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

	err = s.repo.CreateCredentials(ctx, tx, cred)
	if err != nil {
		return "", err
	}

	err = s.repo.DeleteVerification(ctx, tx, cred.CredentialsID)
	if err != nil {
		return "", err
	}

	err = s.repo.CreateVerification(ctx, tx, verif)
	if err != nil {
		return "", err
	}

	err = mail.SendVerificationEmail(email, verif.Code)
	if err != nil {
		return "", err
	}
	token, err := utils.GenerateTemporaryToken(email)
	if err != nil {
		return "", err
	}
	return token, tx.Commit(ctx)
}

func resendVerification(s *AuthService, credentialsID uuid.UUID, email string, ctx context.Context, tx pgx.Tx) error {
	err := s.repo.DeleteVerification(ctx, tx, credentialsID)
	if err != nil {
		return err
	}
	newVerif := models.Verification{
		VerificationID: uuid.New(),
		CredentialsID:  credentialsID,
		Code:           utils.GenerateCode(),
		ExpiresAt:      time.Now().Add(5 * time.Minute),
	}
	err = s.repo.CreateVerification(ctx, tx, newVerif)
	if err != nil {
		return err
	}
	err = mail.SendVerificationEmail(email, newVerif.Code)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) ResendVerificationWithToken(email string) (string, error) {
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
	err = resendVerification(s, cred.CredentialsID, cred.Email, ctx, tx)
	if err != nil {
		return "", err
	}
	token, err := utils.GenerateTemporaryToken(cred.Email)
	if err != nil {
		return "", err
	}
	return token, tx.Commit(ctx)
}

func (s *AuthService) ResendVerificationForChangeEmail(profileID uuid.UUID, newEmail string) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	cred, err := s.repo.GetCredentialsByProfile(ctx, tx, profileID)
	if err != nil {
		return err
	}
	changeEmail, err := s.repo.GetEmailChange(ctx, tx, cred.CredentialsID)
	if err != nil {
		return err
	}
	if newEmail != changeEmail.NewEmail {
		return errors.New("New email does not match")
	}
	err = resendVerification(s, cred.CredentialsID, newEmail, ctx, tx)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *AuthService) ResendVerificationForChangePassword(profileID uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	cred, err := s.repo.GetCredentialsByProfile(ctx, tx, profileID)
	if err != nil {
		return err
	}
	err = resendVerification(s, cred.CredentialsID, cred.Email, ctx, tx)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
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
	return access, refreshT.RefreshToken, tx.Commit(ctx)
}

func (s *AuthService) AuthRefresh(refresh string) (string, string, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return "", "", err
	}
	defer tx.Rollback(ctx)
	claims, err := utils.ParseRefreshToken(refresh)
	if err != nil {
		return "", "", err
	}

	tokenId, err := uuid.Parse(claims["token_id"].(string))
	if err != nil {
		return "", "", errors.New("Invalid refresh token")
	}
	exists, err := s.repo.RefreshExists(ctx, tx, tokenId)
	if err != nil {
		return "", "", err
	}
	if !exists {
		return "", "", errors.New("Invalid refresh token")
	}
	profileId, err := uuid.Parse(claims["profile_id"].(string))
	if err != nil {
		return "", "", errors.New("Invalid refresh token")
	}

	access, err := utils.GenerateAccessToken(profileId)
	if err != nil {
		return "", "", err
	}
	newRefreshID := uuid.New()
	expTime, tokenStr, err := utils.GenerateRefreshToken(profileId, newRefreshID)
	if err != nil {
		return "", "", err
	}
	newRefreshToken := models.RefreshToken{
		RefreshTokenID: newRefreshID,
		ProfileID:      profileId,
		RefreshToken:   tokenStr,
		ExpiresAt:      expTime,
	}

	err = s.repo.DeleteRefreshToken(ctx, tx, tokenId)
	if err != nil {
		return "", "", err
	}

	err = s.repo.SaveRefreshToken(ctx, tx, newRefreshToken.ProfileID, newRefreshToken.RefreshTokenID, newRefreshToken.RefreshToken, expTime)
	if err != nil {
		return "", "", err
	}
	return access, newRefreshToken.RefreshToken, tx.Commit(ctx)
}

func (s *AuthService) Logout(refresh string) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	claims, err := utils.ParseRefreshToken(refresh)
	if err != nil {
		if !errors.Is(err, jwt.ErrTokenExpired) {
			return errors.New("failed to logout")
		}
	}
	tokenId, err := uuid.Parse(claims["token_id"].(string))
	if err != nil {
		return errors.New("Invalid refresh token")
	}
	err = s.repo.DeleteRefreshToken(ctx, tx, tokenId)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *AuthService) LogoutAll(profile_id uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = s.repo.DeleteRefreshByProfile(ctx, tx, profile_id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *AuthService) ChangeEmail(profileID uuid.UUID, newEmail, password string) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	cred, err := s.repo.GetCredentialsByProfile(ctx, tx, profileID)
	if err != nil {
		return err
	}
	if ok := utils.CheckPasswordHash(password, cred.HashedPassword); !ok {
		return errors.New("Invalid password")
	}
	exists, err := s.repo.EmailExists(ctx, tx, newEmail)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Email already exists")
	}
	err = s.repo.DeleteEmailChangeByCreds(ctx, tx, cred.CredentialsID)
	if err != nil {
		return err
	}
	emailChangeID := uuid.New()
	err = s.repo.SaveEmailChange(ctx, tx, emailChangeID, cred.CredentialsID, newEmail)
	if err != nil {
		return err
	}
	err = s.repo.DeleteVerification(ctx, tx, cred.CredentialsID)
	if err != nil {
		return err
	}

	newVerif := models.Verification{
		VerificationID: uuid.New(),
		CredentialsID:  cred.CredentialsID,
		Code:           utils.GenerateCode(),
		ExpiresAt:      time.Now().Add(5 * time.Minute),
	}

	err = s.repo.CreateVerification(ctx, tx, newVerif)
	if err != nil {
		return err
	}
	err = mail.SendVerificationEmail(newEmail, newVerif.Code)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *AuthService) ChangeEmailConfirm(profileID uuid.UUID, code string) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	cred, err := s.repo.GetCredentialsByProfile(ctx, tx, profileID)
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
	emailChange, err := s.repo.GetEmailChange(ctx, tx, cred.CredentialsID)
	if err != nil {
		return err
	}
	exists, err := s.repo.EmailExists(ctx, tx, emailChange.NewEmail)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Email already exists")
	}
	err = s.repo.ChangeEmail(ctx, tx, cred.CredentialsID, emailChange.NewEmail)
	if err != nil {
		return err
	}
	err = s.repo.DeleteEmailChangeByCreds(ctx, tx, cred.CredentialsID)
	if err != nil {
		return err
	}
	err = s.repo.DeleteVerification(ctx, tx, cred.CredentialsID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *AuthService) ChangePassword(profileID uuid.UUID, oldPassword, newPassword string) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	cred, err := s.repo.GetCredentialsByProfile(ctx, tx, profileID)
	if err != nil {
		return err
	}
	if ok := utils.CheckPasswordHash(oldPassword, cred.HashedPassword); !ok {
		return errors.New("Invalid old password")
	}
	newPasswordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	err = s.repo.DeletePasswordChangeByCreds(ctx, tx, cred.CredentialsID)
	if err != nil {
		return err
	}
	passwordChangeID := uuid.New()
	err = s.repo.SavePasswordChange(ctx, tx, passwordChangeID, cred.CredentialsID, newPasswordHash)
	if err != nil {
		return err
	}
	err = s.repo.DeleteVerification(ctx, tx, cred.CredentialsID)
	if err != nil {
		return err
	}

	newVerif := models.Verification{
		VerificationID: uuid.New(),
		CredentialsID:  cred.CredentialsID,
		Code:           utils.GenerateCode(),
		ExpiresAt:      time.Now().Add(5 * time.Minute),
	}

	err = s.repo.CreateVerification(ctx, tx, newVerif)
	if err != nil {
		return err
	}
	err = mail.SendVerificationEmail(cred.Email, newVerif.Code)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *AuthService) ChangePasswordConfirm(profileID uuid.UUID, code string) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	cred, err := s.repo.GetCredentialsByProfile(ctx, tx, profileID)
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
	changePassword, err := s.repo.GetPasswordChange(ctx, tx, cred.CredentialsID)
	if err != nil {
		return err
	}
	err = s.repo.ChangePassword(ctx, tx, cred.CredentialsID, changePassword.NewPasswordHash)
	if err != nil {
		return err
	}
	err = s.repo.DeletePasswordChangeByCreds(ctx, tx, cred.CredentialsID)
	if err != nil {
		return err
	}
	err = s.repo.DeleteVerification(ctx, tx, cred.CredentialsID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *AuthService) ResetPasswordRequest(email string) (string, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	cred, err := s.repo.GetCredentialsByEmail(ctx, tx, email)
	if err != nil {
		return "", err
	}
	verif := models.Verification{
		VerificationID: uuid.New(),
		CredentialsID:  cred.CredentialsID,
		Code:           utils.GenerateCode(),
		ExpiresAt:      time.Now().Add(5 * time.Minute),
	}
	err = s.repo.DeleteVerification(ctx, tx, cred.CredentialsID)
	if err != nil {
		return "", err
	}
	err = s.repo.CreateVerification(ctx, tx, verif)
	if err != nil {
		return "", err
	}
	err = mail.SendVerificationEmail(email, verif.Code)
	if err != nil {
		return "", err
	}
	token, err := utils.GenerateTemporaryToken(email)
	if err != nil {
		return "", err
	}
	return token, tx.Commit(ctx)
}

func (s *AuthService) ResetPasswordConfirm(token, code, newPassword string) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	claims, err := utils.ParseTemporaryToken(token)
	if err != nil {
		return err
	}
	email := claims["email"].(string)
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
	newPasswordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	err = s.repo.ChangePassword(ctx, tx, cred.CredentialsID, newPasswordHash)
	if err != nil {
		return err
	}
	err = s.repo.DeleteVerification(ctx, tx, cred.CredentialsID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
