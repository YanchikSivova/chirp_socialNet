package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"messageService/models"
	"net/http"
)

type PostClient struct {
	baseURL string
	client  *http.Client
}

func NewPostClient() *PostClient {
	return &PostClient{
		baseURL: "http://post-service:8181",
		client:  &http.Client{},
	}
}

func (c *PostClient) BatchPosts(posts models.PostIDs) (*models.PostPreviews, error) {
	url := fmt.Sprintf("%v/internal/posts", c.baseURL)
	body, err := json.Marshal(&posts)
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
		return nil, fmt.Errorf("unexpected status code: %d", resp.Status)
	}
	var previews models.PostPreviews
	err = json.NewDecoder(resp.Body).Decode(&previews)
	if err != nil {
		return nil, err
	}
	return &previews, nil
}
