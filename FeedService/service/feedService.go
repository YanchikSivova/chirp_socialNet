package service

import (
	"context"
	"errors"
	"feedService/client"
	"feedService/kafkaEvents"
	"feedService/repository"
	"github.com/google/uuid"
	"log"
)

type FeedService struct {
	repo          *repository.FeedRepository
	accountClient *client.Client
}

func NewFeedService(repo *repository.FeedRepository, ac *client.Client) *FeedService {
	return &FeedService{
		repo:          repo,
		accountClient: ac,
	}
}

func (s *FeedService) ProcessNewPost(event kafkaEvents.PostEvent) error {
	followers, err := s.accountClient.GetFollowers(event.ProfileID)
	if err != nil {
		return err
	}
	if len(followers.Followers) == 0 {
		return nil
	}
	ctx := context.Background()
	successCount := 0
	for _, follower := range followers.Followers {
		err = s.repo.AddPostToFeed(ctx, follower.String(), event.PostID.String(), float64(event.CreatedAt.Unix()))
		if err != nil {
			log.Println(err)
			continue
		}
		successCount++
	}
	if successCount == 0 {
		return errors.New("failed to add post to all feeds")
	}
	return nil
}

func (s *FeedService) GetFeed(profileID uuid.UUID) ([]uuid.UUID, error) {
	ctx := context.Background()
	postsIDStr, err := s.repo.GetFeed(ctx, profileID.String())
	if err != nil {
		return []uuid.UUID{}, err
	}
	var posts []uuid.UUID
	for _, postIDStr := range postsIDStr {
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			err = s.repo.RemovePostFromFeed(ctx, profileID.String(), postIDStr)
			log.Printf("failed to remove post from feed: %v", err)
		}
		posts = append(posts, postID)
	}
	return posts, nil
}
