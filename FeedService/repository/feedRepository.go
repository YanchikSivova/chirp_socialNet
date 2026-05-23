package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type FeedRepository struct {
	Redis *redis.Client
}

func NewFeedRepository(client *redis.Client) *FeedRepository {
	return &FeedRepository{Redis: client}
}

func (r *FeedRepository) AddPostToFeed(ctx context.Context, profileID string, postID string, score float64) error {
	key := fmt.Sprintf("feed:%s", profileID)

	err := r.Redis.ZAdd(ctx, key, redis.Z{Score: score, Member: postID}).Err()
	if err != nil {
		return err
	}
	return r.Redis.ZRemRangeByRank(ctx, key, 0, -21).Err()
}

func (r *FeedRepository) GetFeed(ctx context.Context, profileID string) ([]string, error) {
	key := fmt.Sprintf("feed:%s", profileID)
	return r.Redis.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   key,
		Start: 0,
		Stop:  19,
		Rev:   true,
	}).Result()
}

func (r *FeedRepository) RemovePostFromFeed(ctx context.Context, profileID string, postID string) error {
	key := fmt.Sprintf("feed:%s", profileID)
	return r.Redis.ZRem(ctx, key, postID).Err()
}
