package repository

import (
	"accountService/models"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository struct {
	DB *pgxpool.Pool
}

func NewUsersRepository(DB *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{DB: DB}
}

func (r *UsersRepository) UsernameExists(ctx context.Context, tx pgx.Tx, username string) (bool, error) {
	var exists bool
	row := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM profile WHERE username=$1)`, username)
	err := row.Scan(&exists)
	return exists, err
}

func (r *UsersRepository) FillProfile(ctx context.Context, tx pgx.Tx, profileID uuid.UUID, name, username, description, avatar string) error {
	_, err := tx.Exec(ctx, `UPDATE profile SET name=$1, username=$2, avatar=$3, description=$4, is_completed=true WHERE profile_id=$5`, name, username, avatar, description, profileID)
	return err
}

func (r *UsersRepository) CheckProfileCompleted(ctx context.Context, tx pgx.Tx, profileID uuid.UUID) (bool, error) {
	var isCompleted bool
	row := tx.QueryRow(ctx, `SELECT is_completed FROM profile WHERE profile_id=$1`, profileID)
	err := row.Scan(&isCompleted)
	return isCompleted, err
}

func (r *UsersRepository) UpdateName(ctx context.Context, tx pgx.Tx, profileID uuid.UUID, newName string) error {
	_, err := tx.Exec(ctx, `UPDATE profile SET name=$1 WHERE profile_id=$2`, newName, profileID)
	return err
}

func (r *UsersRepository) UpdateUsername(ctx context.Context, tx pgx.Tx, profileID uuid.UUID, newUsername string) error {
	_, err := tx.Exec(ctx, `UPDATE profile SET username=$1 WHERE profile_id=$2`, newUsername, profileID)
	return err
}

func (r *UsersRepository) UpdateDescription(ctx context.Context, tx pgx.Tx, profileID uuid.UUID, newDescription string) error {
	_, err := tx.Exec(ctx, `UPDATE profile SET description=$1 WHERE profile_id=$2`, newDescription, profileID)
	return err
}

func (r *UsersRepository) UpdateAvatar(ctx context.Context, tx pgx.Tx, profileID uuid.UUID, newAvatar string) error {
	_, err := tx.Exec(ctx, `UPDATE profile SET avatar=$1 WHERE profile_id=$2`, newAvatar, profileID)
	return err
}

func (r *UsersRepository) GetProfileByProfileID(ctx context.Context, tx pgx.Tx, profileID uuid.UUID) (*models.Profile, error) {
	var profile models.Profile
	row := tx.QueryRow(ctx, `SELECT * FROM profile WHERE profile_id=$1`, profileID)
	err := row.Scan(
		&profile.ProfileID,
		&profile.Name,
		&profile.Username,
		&profile.Avatar,
		&profile.Description,
		&profile.SubscribersAmount,
		&profile.SubscribedAmount,
		&profile.PostsAmount,
		&profile.IsCompleted)
	return &profile, err
}
