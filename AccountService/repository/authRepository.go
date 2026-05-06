package repository

import (
	"accountService/models"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type AuthRepository struct {
	DB *pgxpool.Pool
}

func NewAuthRepository(DB *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{DB: DB}
}

func (r *AuthRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM credentials WHERE email = $1)`, email).Scan(&exists)
	return exists, err
}

// Создать credentials
func (r *AuthRepository) CreateCredentials(ctx context.Context, tx pgx.Tx, c models.Credentials) error {
	_, err := tx.Exec(ctx, `INSERT INTO credentials (credentials_id, email, hashed_password, created_at, status) VALUES ($1, $2, $3, $4, $5)`, c.CredentialsID, c.Email, c.HashedPassword, c.CreatedAt, c.Status)
	return err
}

// Создать verification
func (r *AuthRepository) CreateVerification(ctx context.Context, tx pgx.Tx, v models.Verification) error {
	_, err := tx.Exec(ctx, `INSERT INTO verification (verification_id, credentials_id, code, expires_at) VALUES ($1, $2, $3, $4)`, v.VerificationID, v.CredentialsID, v.Code, v.ExpiresAt)
	return err
}

// Удалить старый verification
func (r *AuthRepository) DeleteVerification(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	_, err := tx.Exec(ctx, `DELETE FROM verification WHERE credentials_id = $1`, id)
	return err
}

// Транзакция создание credentials и verification
func (r *AuthRepository) RegisterTx(ctx context.Context, c models.Credentials, v models.Verification) error {
	tx, err := r.DB.Begin(ctx) //транзакция
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //откат если упало

	err = r.CreateCredentials(ctx, tx, c)
	if err != nil {
		return err
	}

	err = r.DeleteVerification(ctx, tx, c.CredentialsID)
	if err != nil {
		return err
	}

	err = r.CreateVerification(ctx, tx, v)
	if err != nil {
		return err
	}

	return tx.Commit(ctx) //если Commit то Rollback игнорируется
}

// Получение credentials по email
func (r *AuthRepository) GetCredentialsByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.Credentials, error) {
	row := tx.QueryRow(ctx, "SELECT credentials_id, profile_id, email, hashed_password, created_at, status FROM credentials WHERE email = $1", email)
	var cred models.Credentials
	err := row.Scan(
		&cred.CredentialsID,
		&cred.ProfileID,
		&cred.Email,
		&cred.HashedPassword,
		&cred.CreatedAt,
		&cred.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("credentials not found")
		}
		return nil, err
	}
	return &cred, nil
}

// Получение verification по credentials_id
func (r *AuthRepository) GetVerificationByCredentialsID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*models.Verification, error) {
	row := tx.QueryRow(ctx, "SELECT verification_id, credentials_id, code, expires_at FROM verification WHERE credentials_id = $1", id)
	var verif models.Verification
	err := row.Scan(
		&verif.VerificationID,
		&verif.CredentialsID,
		&verif.Code,
		&verif.ExpiresAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("verification not found")
		}
		return nil, err
	}
	return &verif, nil
}

// Создание профиля
func (r *AuthRepository) CreateProfile(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	_, err := tx.Exec(ctx, `INSERT INTO profile (profile_id)  VALUES ($1)`, id)
	return err
}

// Сопоставление профиля и credentials
func (r *AuthRepository) AttachProfile(ctx context.Context, tx pgx.Tx, credID, profID uuid.UUID) error {
	_, err := tx.Exec(ctx, `UPDATE credentials SET profile_id=$1 WHERE credentials_id=$2`, profID, credID)
	return err
}

// Смена статуса credentials на active
func (r *AuthRepository) ActivateCredentials(ctx context.Context, tx pgx.Tx, credID uuid.UUID) error {
	_, err := tx.Exec(ctx, `UPDATE credentials SET status='active' WHERE credentials_id=$1`, credID)
	return err
}

// Сохранение refresh токена
func (r *AuthRepository) SaveRefreshToken(ctx context.Context, tx pgx.Tx, profileID, refreshID uuid.UUID, refreshToken string, expTime time.Time) error {
	_, err := tx.Exec(ctx, `INSERT INTO refresh_token (refresh_token_id, profile_id, refresh_token, expires_at) VALUES ($1, $2, $3, $4)`, refreshID, profileID, refreshToken, expTime)
	return err
}

// Удаление refresh токена
func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, tx pgx.Tx, refreshTokenID uuid.UUID) error {
	_, err := tx.Exec(ctx, `DELETE FROM refresh_token WHERE refresh_token_id = $1`, refreshTokenID)
	return err
}
