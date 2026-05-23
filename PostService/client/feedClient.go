package client

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type FeedClient struct {
	baseURL string
	client  *http.Client
}

func NewFeedClient() *FeedClient {
	return &FeedClient{
		baseURL: "http://feed-service:8282",
		client:  &http.Client{},
	}
}

type FeedResponse struct {
	Posts []uuid.UUID `json:"posts"`
}

func (c *FeedClient) GetFeedPosts(profileID uuid.UUID) ([]uuid.UUID, error) {
	log.Println("feed client started")
	url := fmt.Sprintf("%s/internal/feed/%s",
		c.baseURL,
		profileID.String(),
	)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var feedResp FeedResponse
	err = json.NewDecoder(resp.Body).Decode(&feedResp)
	if err != nil {
		return nil, err
	}
	return feedResp.Posts, nil
}
