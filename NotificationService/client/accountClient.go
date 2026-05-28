package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"notificationService/models"
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

func (c *AccountClient) BatchProfiles(profileIDs map[uuid.UUID]struct{}) ([]models.Profile, error) {
	url := fmt.Sprintf("%s/internal/batch-profiles", c.baseURL)
	var IDs models.Profiles
	for key, _ := range profileIDs {
		IDs.Profiles = append(IDs.Profiles, key)
	}
	body, err := json.Marshal(IDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var profiles []models.Profile
	err = json.NewDecoder(resp.Body).Decode(&profiles)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return profiles, nil
}
