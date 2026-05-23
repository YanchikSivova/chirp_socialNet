package handler

import (
	"errors"
	"feedService/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type FeedHandler struct {
	service *service.FeedService
}

func NewFeedHandler(service *service.FeedService) *FeedHandler {
	return &FeedHandler{service: service}
}

func parseIdParam(c *gin.Context) (*uuid.UUID, error) {
	isStr := c.Param("id")
	if isStr == "" {
		return nil, errors.New("id is required")
	}
	id, err := uuid.Parse(isStr)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (h *FeedHandler) GetFeed(c *gin.Context) {
	profileID, err := parseIdParam(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	posts, err := h.service.GetFeed(*profileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}
