package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"postService/client"
	"postService/models"
	"postService/repository"
	"unicode/utf8"
)

type PostService struct {
	repo          *repository.PostRepository
	accountClient *client.AccountClient
}

func NewPostService(repo *repository.PostRepository, ac *client.AccountClient) *PostService {
	return &PostService{repo, ac}
}

func (s *PostService) CreatePost(profileID uuid.UUID, postRequest models.CreatePostRequest) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	exists, err := s.accountClient.ProfileExists(profileID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("profile does not exist")
	}
	contentLen := utf8.RuneCountInString(postRequest.Content)
	if contentLen == 0 || contentLen > 250 {
		return errors.New("invalid post content")
	}
	if len(postRequest.Hashtags) > 3 {
		return errors.New("too many hashtags")
	}
	if len(postRequest.Images) > 3 {
		return errors.New("too many images")
	}
	if postRequest.Status != "published" && postRequest.Status != "draft" {
		return errors.New("invalid post status")
	}
	post := models.Post{
		PostID:    uuid.New(),
		ProfileID: profileID,
		Content:   postRequest.Content,
		Status:    postRequest.Status,
	}
	err = s.repo.CreatePost(ctx, tx, post)
	if err != nil {
		return err
	}
	for i, image := range postRequest.Images {
		newImage := models.Image{
			ImageID:    uuid.New(),
			PostID:     post.PostID,
			Url:        image.Url,
			OrderIndex: i + 1,
		}
		err = s.repo.CreateImage(ctx, tx, newImage)
		if err != nil {
			return err
		}
	}
	for i, hashtag := range postRequest.Hashtags {
		id, err := s.repo.GetHashtagId(ctx, tx, hashtag.HashtagName)
		if err != nil && err != pgx.ErrNoRows {
			return err
		}
		if err == pgx.ErrNoRows {
			id = uuid.New()
			newHashtag := models.Hashtag{
				HashtagID:   id,
				HashtagName: hashtag.HashtagName,
			}
			err = s.repo.CreateHashtag(ctx, tx, newHashtag)
			if err != nil {
				return err
			}
		}
		newPostHashtag := models.PostHashtag{
			PostHashtagID: uuid.New(),
			PostID:        post.PostID,
			HashtagID:     id,
			OrderIndex:    i + 1,
		}
		err = s.repo.CreatePostHashtag(ctx, tx, newPostHashtag)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
