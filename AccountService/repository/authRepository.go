package repository

import (
	"accountService/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM credentials WHERE email = $1)`, email).Scan(&exists)
	return exists, err
}

/*func (r *AuthRepository) CreateCredentials(ctx context.Context, c models.Credentials) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO credentials (credentials_id, email, hashed_password, created_at, status) VALUES ($1, $2, $3, $4, $5)`, c.CredentialsID, c.Email, c.HashedPassword, c.CreatedAt, c.Status)
	return err
}

func (r *AuthRepository) CreateVerification(ctx context.Context, v models.Verification) error {
	_, err := r.db.Exec(ctx, `INSERT INTO verification (verification_id, credentials_id, code, expires_at) VALUES ($1, $2, $3, $4)`, v.VerificationID, v.CredentialsID, v.Code, v.ExpiresAt)
	return err
}*/

func (r *AuthRepository) RegisterTx(ctx context.Context, c models.Credentials, v models.Verification) error {
	tx, err := r.db.Begin(ctx) //транзакция
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //откат если упало

	_, err = tx.Exec(ctx, `INSERT INTO credentials (credentials_id, email, hashed_password, created_at, status) VALUES ($1, $2, $3, $4, $5)`, c.CredentialsID, c.Email, c.HashedPassword, c.CreatedAt, c.Status)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO verification (verification_id, credentials_id, code, expires_at) VALUES ($1, $2, $3, $4)`, v.VerificationID, v.CredentialsID, v.Code, v.ExpiresAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx) //если Commit то Rollback игнорируется
}
