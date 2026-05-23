package client

import (
	"encoding/json"
	"feedService/models"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type Client struct {
	baseUrl string
	client  *http.Client
}

func NewAccountClient() *Client {
	return &Client{
		baseUrl: "http://account-service:8080",
		client:  &http.Client{},
	}
}

func (c *Client) GetFollowers(profileID uuid.UUID) (*models.Followers, error) {
	url := fmt.Sprintf("%s/internal/users/%s/followers", c.baseUrl, profileID.String())
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var followers models.Followers
	err = json.NewDecoder(resp.Body).Decode(&followers)
	if err != nil {
		return nil, err
	}
	return &followers, nil
}
