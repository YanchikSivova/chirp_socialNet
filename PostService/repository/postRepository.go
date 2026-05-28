package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"postService/models"
	"time"
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

func (r *PostRepository) UpdatePostContent(ctx context.Context, tx pgx.Tx, postID uuid.UUID, newContent string) error {
	_, err := tx.Exec(ctx, `UPDATE post SET content = $1 WHERE post_id = $2`, newContent, postID)
	return err
}

func (r *PostRepository) DeleteImages(ctx context.Context, tx pgx.Tx, postID uuid.UUID) error {
	_, err := tx.Exec(ctx, `DELETE FROM image WHERE post_id = $1`, postID)
	return err
}

func (r *PostRepository) DeletePostHashtags(ctx context.Context, tx pgx.Tx, postID uuid.UUID) error {
	_, err := tx.Exec(ctx, `DELETE FROM post_hashtag WHERE post_id = $1`, postID)
	return err
}

func (r *PostRepository) GetAuthorId(ctx context.Context, tx pgx.Tx, postId uuid.UUID) (uuid.UUID, error) {
	var authorId uuid.UUID
	err := tx.QueryRow(ctx, `SELECT profile_id FROM post WHERE post_id = $1`, postId).Scan(&authorId)
	return authorId, err
}

func (r *PostRepository) GetPostStatus(ctx context.Context, tx pgx.Tx, postId uuid.UUID) (string, error) {
	var postStatus string
	err := tx.QueryRow(ctx, `SELECT status FROM post WHERE post_id = $1`, postId).Scan(&postStatus)
	return postStatus, err
}

func (r *PostRepository) UpdateLastEditedTime(ctx context.Context, tx pgx.Tx, postId uuid.UUID) error {
	_, err := tx.Exec(ctx, `UPDATE post SET last_edited_at = $1 WHERE post_id = $2`, time.Now(), postId)
	return err
}

func (r *PostRepository) DeletePost(ctx context.Context, tx pgx.Tx, postId uuid.UUID) error {
	_, err := tx.Exec(ctx, `DELETE FROM post WHERE post_id = $1`, postId)
	return err
}

func (r *PostRepository) GetPost(ctx context.Context, tx pgx.Tx, postId uuid.UUID) (*models.Post, error) {
	var post models.Post
	err := tx.QueryRow(ctx, `SELECT post_id, profile_id, content, last_edited_at, published_at, likes_amount, comments_amount, reposts_amount, status FROM post WHERE post_id=$1`, postId).Scan(
		&post.PostID,
		&post.ProfileID,
		&post.Content,
		&post.LastEditedAt,
		&post.PublishedAt,
		&post.LikesAmount,
		&post.CommentsAmount,
		&post.RepostsAmount,
		&post.Status,
	)
	return &post, err
}

func (r *PostRepository) GetImages(ctx context.Context, tx pgx.Tx, postId uuid.UUID) ([]models.ImageList, error) {
	var images []models.ImageList
	rows, err := tx.Query(ctx, `SELECT url, order_index FROM image WHERE post_id=$1`, postId)
	defer rows.Close()
	if err != nil {
		if err == pgx.ErrNoRows {
			return []models.ImageList{}, nil
		}
		return nil, err
	}
	for rows.Next() {
		var image models.ImageList
		err = rows.Scan(&image.Url, &image.OrderIndex)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}

func (r *PostRepository) GetPostHashtags(ctx context.Context, tx pgx.Tx, postID uuid.UUID) ([]models.HashtagList, error) {
	var hashtags []models.HashtagList
	rows, err := tx.Query(ctx, `SELECT h.hashtag_name, ph.order_index FROM post_hashtag ph JOIN hashtag h ON ph.hashtag_id = h.hashtag_id WHERE ph.post_id = $1`, postID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []models.HashtagList{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var hashtag models.HashtagList
		err = rows.Scan(&hashtag.HashtagName, &hashtag.OrderIndex)
		if err != nil {
			return nil, err
		}
		hashtags = append(hashtags, hashtag)
	}
	return hashtags, nil
}

func (r *PostRepository) CheckLike(ctx context.Context, tx pgx.Tx, postId, profileId uuid.UUID) (bool, error) {
	var liked bool
	err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = $1 AND profile_id = $2)`, postId, profileId).Scan(&liked)
	return liked, err

}

func (r *PostRepository) CheckRepost(ctx context.Context, tx pgx.Tx, postId, profileId uuid.UUID) (bool, error) {
	var reposted bool
	err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM repost WHERE post_id = $1 AND profile_id = $2)`, postId, profileId).Scan(&reposted)
	return reposted, err
}

func (r *PostRepository) CreateLike(ctx context.Context, tx pgx.Tx, postID, profileID uuid.UUID) error {
	_, err := tx.Exec(ctx, `INSERT INTO likes (likes_id, post_id, profile_id) VALUES ($1, $2, $3)`, uuid.New(), postID, profileID)
	return err
}

func (r *PostRepository) DeleteLike(ctx context.Context, tx pgx.Tx, postID, profileID uuid.UUID) error {
	_, err := tx.Exec(ctx, `DELETE FROM likes WHERE post_id=$1 AND profile_id=$2`, postID, profileID)
	return err
}

func (r *PostRepository) CreateRepost(ctx context.Context, tx pgx.Tx, postID, profileID uuid.UUID) error {
	_, err := tx.Exec(ctx, `INSERT INTO repost (repost_id, post_id, profile_id) VALUES ($1, $2, $3)`, uuid.New(), postID, profileID)
	return err
}

func (r *PostRepository) DeleteRepost(ctx context.Context, tx pgx.Tx, postID, profileID uuid.UUID) error {
	_, err := tx.Exec(ctx, `DELETE FROM repost WHERE post_id=$1 AND profile_id=$2`, postID, profileID)
	return err
}

func (r *PostRepository) CreateReport(ctx context.Context, tx pgx.Tx, postID, profileID uuid.UUID, reason string) error {
	_, err := tx.Exec(ctx, `INSERT INTO report (report_id, post_id, profile_id, reason) VALUES ($1, $2, $3, $4)`, uuid.New(), postID, profileID, reason)
	return err
}

func (r *PostRepository) CreateComment(ctx context.Context, tx pgx.Tx, commentID, postID, profileID uuid.UUID, content string) error {
	_, err := tx.Exec(ctx, `INSERT INTO comment (comment_id, post_id, profile_id, content) VALUES ($1, $2, $3, $4)`, commentID, postID, profileID, content)
	return err
}

func (r *PostRepository) CreateCommentWithParent(ctx context.Context, tx pgx.Tx, commentID, postID, profileID, parentCommentID uuid.UUID, content string) error {
	_, err := tx.Exec(ctx, `INSERT INTO comment (comment_id, post_id, profile_id, parent_comment_id, content) VALUES ($1, $2, $3, $4, $5)`, commentID, postID, profileID, parentCommentID, content)
	return err
}

func (r *PostRepository) CreateCommentLike(ctx context.Context, tx pgx.Tx, commentID, profileID uuid.UUID) error {
	_, err := tx.Exec(ctx, `INSERT INTO comment_like (comment_like_id, comment_id, profile_id) VALUES ($1, $2, $3)`, uuid.New(), commentID, profileID)
	return err
}

func (r *PostRepository) DeleteCommentLike(ctx context.Context, tx pgx.Tx, commentID, profileID uuid.UUID) error {
	_, err := tx.Exec(ctx, `DELETE FROM comment_like WHERE comment_id=$1 AND profile_id=$2`, commentID, profileID)
	return err
}

func (r *PostRepository) GetPostIDFromComment(ctx context.Context, tx pgx.Tx, commentID uuid.UUID) (uuid.UUID, error) {
	var postID uuid.UUID
	err := tx.QueryRow(ctx, `SELECT post_id FROM comment WHERE comment_id=$1`, commentID).Scan(&postID)
	return postID, err
}

func (r *PostRepository) GetComment(ctx context.Context, tx pgx.Tx, commentID uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	err := tx.QueryRow(ctx, `SELECT * FROM comment WHERE comment_id=$1`, commentID).Scan(
		&comment.CommentID,
		&comment.PostID,
		&comment.ProfileID,
		&comment.ParentCommentID,
		&comment.Content,
		&comment.LikesAmount,
		&comment.CreatedAt,
	)
	return &comment, err
}

func (r *PostRepository) CheckCommentLike(ctx context.Context, tx pgx.Tx, commentID, profileID uuid.UUID) (bool, error) {
	var isLiked bool
	err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM comment_like WHERE comment_id=$1 AND profile_id=$2)`, commentID, profileID).Scan(&isLiked)
	return isLiked, err
}

func (r *PostRepository) GetProfileIDFromComment(ctx context.Context, tx pgx.Tx, commentID uuid.UUID) (uuid.UUID, error) {
	var profileID uuid.UUID
	err := tx.QueryRow(ctx, `SELECT profile_id FROM comment WHERE comment_id=$1`, commentID).Scan(&profileID)
	return profileID, err
}

func (r *PostRepository) DeleteComment(ctx context.Context, tx pgx.Tx, commentID uuid.UUID) error {
	_, err := tx.Exec(ctx, `DELETE FROM comment WHERE comment_id=$1`, commentID)
	return err
}

func (r *PostRepository) PublishPost(ctx context.Context, tx pgx.Tx, postID uuid.UUID) error {
	_, err := tx.Exec(ctx, `UPDATE post SET status='published' WHERE post_id=$1`, postID)
	return err
}

func (r *PostRepository) GetPosts(ctx context.Context, tx pgx.Tx, profileID uuid.UUID, status string, offset, limit int) ([]models.Post, error) {
	var posts []models.Post
	rows, err := tx.Query(ctx, `SELECT post_id, content, published_at, last_edited_at, likes_amount, comments_amount, reposts_amount FROM post WHERE status=$1 AND profile_id=$2 ORDER BY last_edited_at DESC  LIMIT $3 OFFSET $4 `, status, profileID, limit, offset)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []models.Post{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		err = rows.Scan(
			&post.PostID,
			&post.Content,
			&post.PublishedAt,
			&post.LastEditedAt,
			&post.LikesAmount,
			&post.CommentsAmount,
			&post.RepostsAmount,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) GetComments(ctx context.Context, tx pgx.Tx, postID uuid.UUID, limit, offset int) ([]models.Comment, error) {
	var comments []models.Comment
	rows, err := tx.Query(ctx, `SELECT comment_id, profile_id, content, likes_amount, created_at FROM comment WHERE post_id=$1 AND parent_comment_id is NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3 `, postID, limit, offset)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []models.Comment{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		err = rows.Scan(
			&comment.CommentID,
			&comment.ProfileID,
			&comment.Content,
			&comment.LikesAmount,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *PostRepository) GetCommentAnswers(ctx context.Context, tx pgx.Tx, commentID uuid.UUID, limit, offset int) ([]models.Comment, error) {
	var comments []models.Comment
	rows, err := tx.Query(ctx, `SELECT comment_id, profile_id, content, likes_amount, created_at FROM comment WHERE parent_comment_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3 `, commentID, limit, offset)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []models.Comment{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		err = rows.Scan(
			&comment.CommentID,
			&comment.ProfileID,
			&comment.Content,
			&comment.LikesAmount,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *PostRepository) GetPostsContent(ctx context.Context, tx pgx.Tx, postID uuid.UUID) (string, error) {
	var content string
	err := tx.QueryRow(ctx, `SELECT content FROM post WHERE post_id=$1 and status='published'`, postID).Scan(&content)
	return content, err
}

func (r *PostRepository) GetCommentAuthor(ctx context.Context, tx pgx.Tx, commentID uuid.UUID) (uuid.UUID, error) {
	var authorID uuid.UUID
	err := tx.QueryRow(ctx, `SELECT profile_id FROM comment WHERE comment_id = $1`, commentID).Scan(&authorID)
	return authorID, err
}
