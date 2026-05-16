package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"postService/models"
	"postService/service"
)

type PostHandler struct {
	service *service.PostService
}

func NewPostHandler(service *service.PostService) *PostHandler { return &PostHandler{service: service} }

func parseProfileHeader(c *gin.Context) (*uuid.UUID, error) {
	profileIdStr := c.GetHeader("X-User-Id")
	if profileIdStr == "" {
		return nil, errors.New("no profile_id found")
	}
	profileId, err := uuid.Parse(profileIdStr)
	if err != nil {
		return nil, err
	}
	return &profileId, nil
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	profileID, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.CreatePost(*profileID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, SuccessResponse{Success: true})
}
