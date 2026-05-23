package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"postService/models"
)

type AccountClient struct {
	baseURL string
	client  *http.Client
}

func NewAccountClient() *AccountClient {
	return &AccountClient{
		baseURL: "http://account-service:8080",
		client:  &http.Client{},
	}
}

func (c *AccountClient) ProfileExists(profileID uuid.UUID) (bool, error) {
	url := fmt.Sprintf(
		"%s/internal/users/%s/exists",
		c.baseURL,
		profileID.String(),
	)
	resp, err := c.client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var result struct {
		Exists bool `json:"exists"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return false, err
	}
	return result.Exists, nil
}

func (c *AccountClient) ProfileMinimum(profileID uuid.UUID) (*models.Profile, error) {
	url := fmt.Sprintf(
		"%s/internal/users/%s/profile",
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
	var profile models.Profile
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

type RelationshipRequest struct {
	ProfileID uuid.UUID `json:"profile_id"`
	AuthorID  uuid.UUID `json:"author_id"`
}

type RelationshipResponse struct {
	Relationship string `json:"relationship"`
}

func (c *AccountClient) GetRelationship(profileID, authorID uuid.UUID) (string, error) {
	url := fmt.Sprintf("%v/internal/relationship", c.baseURL)
	req := RelationshipRequest{
		ProfileID: profileID,
		AuthorID:  authorID,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var relationship RelationshipResponse
	err = json.NewDecoder(resp.Body).Decode(&relationship)
	if err != nil {
		return "", err
	}
	log.Println(relationship.Relationship)
	return relationship.Relationship, nil
}
