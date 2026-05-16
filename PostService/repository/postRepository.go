package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"postService/models"
)

type PostRepository struct {
	DB *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *PostRepository { return &PostRepository{db} }

func (r *PostRepository) CreatePost(ctx context.Context, tx pgx.Tx, post models.Post) error {
	_, err := tx.Exec(ctx, `INSERT INTO post (post_id, profile_id, content, status) VALUES ($1, $2, $3, $4)`, post.PostID, post.ProfileID, post.Content, post.Status)
	return err
}

func (r *PostRepository) CreateImage(ctx context.Context, tx pgx.Tx, image models.Image) error {
	_, err := tx.Exec(ctx, `INSERT INTO image (image_id, post_id, url, order_index) VALUES ($1, $2, $3, $4)`, image.ImageID, image.PostID, image.Url, image.OrderIndex)
	return err
}

func (r *PostRepository) GetHashtagId(ctx context.Context, tx pgx.Tx, hashtagName string) (uuid.UUID, error) {
	var hashtagId uuid.UUID
	err := tx.QueryRow(ctx, `SELECT hashtag_id FROM hashtag WHERE hashtag_name = $1`, hashtagName).Scan(&hashtagId)
	return hashtagId, err
}

func (r *PostRepository) CreateHashtag(ctx context.Context, tx pgx.Tx, hashtag models.Hashtag) error {
	_, err := tx.Exec(ctx, `INSERT INTO hashtag (hashtag_id, hashtag_name) VALUES ($1, $2)`, hashtag.HashtagID, hashtag.HashtagName)
	return err
}

func (r *PostRepository) CreatePostHashtag(ctx context.Context, tx pgx.Tx, postHashtag models.PostHashtag) error {
	_, err := tx.Exec(ctx, `INSERT INTO post_hashtag (post_hashtag_id, post_id, hashtag_id, order_index) VALUES ($1, $2, $3, $4)`, postHashtag.PostHashtagID, postHashtag.PostID, postHashtag.HashtagID, postHashtag.OrderIndex)
	return err
}
