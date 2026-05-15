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

func (r *UsersRepository) CheckIsSubscribed(ctx context.Context, tx pgx.Tx, subscriberId, subscribedId uuid.UUID) (bool, error) {
	var isSubscribed bool
	row := tx.QueryRow(ctx, `SELECT exists(SELECT 1 FROM subscription WHERE subscriber_id=$1 AND subscribed_id=$2)`, subscriberId, subscribedId)
	err := row.Scan(&isSubscribed)
	return isSubscribed, err
}

func (r *UsersRepository) Follow(ctx context.Context, tx pgx.Tx, subscriberID, subscribedID uuid.UUID) error {
	_, err := tx.Exec(ctx, `INSERT INTO subscription (subscription_id, subscriber_id, subscribed_id) VALUES ($1, $2, $3)`, uuid.New(), subscriberID, subscribedID)
	return err
}

func (r *UsersRepository) UpdateFollowersAmount(ctx context.Context, tx pgx.Tx, profileID uuid.UUID) error {
	_, err := tx.Exec(ctx, `UPDATE profile SET subscribers_amount = (SELECT COUNT(*) FROM subscription WHERE subscribed_id=$1) WHERE profile_id = $1`, profileID)
	return err
}

func (r *UsersRepository) UpdateFollowingsAmount(ctx context.Context, tx pgx.Tx, profileID uuid.UUID) error {
	_, err := tx.Exec(ctx, `UPDATE profile SET subscribed_amount = (SELECT COUNT(*) FROM subscription WHERE subscriber_id=$1) WHERE profile_id = $1`, profileID)
	return err
}

func (r *UsersRepository) Unfollow(ctx context.Context, tx pgx.Tx, subscriberID, subscribedID uuid.UUID) error {
	_, err := tx.Exec(ctx, `DELETE FROM subscription WHERE subscriber_id=$1 AND subscribed_id=$2`, subscriberID, subscribedID)
	return err
}

func (r *UsersRepository) CheckIsBlocked(ctx context.Context, tx pgx.Tx, bannedProfileId, profileId uuid.UUID) (bool, error) {
	var isBlocked bool
	row := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM blacklist WHERE banned_profile_id=$1 AND profile_id=$2)`, bannedProfileId, profileId)
	err := row.Scan(&isBlocked)
	return isBlocked, err
}

func (r *UsersRepository) Block(ctx context.Context, tx pgx.Tx, bannedProfileId, profileId uuid.UUID) error {
	_, err := tx.Exec(ctx, `INSERT INTO blacklist (blacklist_id, banned_profile_id, profile_id) VALUES ($1, $2, $3)`, uuid.New(), bannedProfileId, profileId)
	return err
}

func (r *UsersRepository) Unblock(ctx context.Context, tx pgx.Tx, bannedProfileId, profileId uuid.UUID) error {
	_, err := tx.Exec(ctx, `DELETE FROM blacklist WHERE banned_profile_id=$1 and profile_id=$2`, bannedProfileId, profileId)
	return err
}

func (r *UsersRepository) GetFollowers(ctx context.Context, tx pgx.Tx, profileId uuid.UUID, limit int, offset int) ([]models.ProfileMinimum, error) {
	var followers []models.ProfileMinimum
	rows, err := tx.Query(ctx, `SELECT p.profile_id, p.name, p.username, p.avatar FROM subscription s JOIN profile p ON p.profile_id = s.subscriber_id WHERE s.subscribed_id = $1 LIMIT $2 OFFSET $3`, profileId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var profile models.ProfileMinimum
		err = rows.Scan(
			&profile.ProfileID,
			&profile.Name,
			&profile.Username,
			&profile.Avatar,
		)
		if err != nil {
			return nil, err
		}
		followers = append(followers, profile)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return followers, nil
}

func (r *UsersRepository) GetFollowings(ctx context.Context, tx pgx.Tx, profileId uuid.UUID, limit int, offset int) ([]models.ProfileMinimum, error) {
	var followings []models.ProfileMinimum
	rows, err := tx.Query(ctx, `SELECT p.profile_id, p.name, p.username, p.avatar FROM subscription s JOIN profile p ON p.profile_id = s.subscribed_id WHERE subscriber_id=$1 LIMIT $2 OFFSET $3`, profileId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var profile models.ProfileMinimum
		err = rows.Scan(
			&profile.ProfileID,
			&profile.Name,
			&profile.Username,
			&profile.Avatar,
		)
		if err != nil {
			return nil, err
		}
		followings = append(followings, profile)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return followings, nil
}

func (r *UsersRepository) GetBlacklist(ctx context.Context, tx pgx.Tx, profileId uuid.UUID, limit int, offset int) ([]models.ProfileMinimum, error) {
	var blacklist []models.ProfileMinimum
	rows, err := tx.Query(ctx, `SELECT p.profile_id, p.name, p.username, p.avatar FROM blacklist b JOIN profile p ON p.profile_id = b.banned_profile_id WHERE b.profile_id = $1 LIMIT $2 OFFSET $3`, profileId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var profile models.ProfileMinimum
		err = rows.Scan(
			&profile.ProfileID,
			&profile.Name,
			&profile.Username,
			&profile.Avatar,
		)
		if err != nil {
			return nil, err
		}
		blacklist = append(blacklist, profile)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return blacklist, nil
}
