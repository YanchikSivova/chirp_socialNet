package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"notificationService/service"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler(service *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func parseProfileIdHeader(c *gin.Context) (*uuid.UUID, error) {
	profileIdStr := c.GetHeader("X-User-Id")
	if profileIdStr == "" {
		return nil, errors.New("no profile id")
	}
	profileId, err := uuid.Parse(profileIdStr)
	if err != nil {
		return nil, errors.New("invalid profile id")
	}
	return &profileId, nil
}

func parseIdParam(c *gin.Context) (*uuid.UUID, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return nil, errors.New("no id")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	return &id, nil
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	profileId, err := parseProfileIdHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	notifications, err := h.service.GetNotifications(*profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	profileId, err := parseProfileIdHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	notifID, err := parseIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.DeleteNotification(*profileId, *notifID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Success: true})
}

func (h *NotificationHandler) DeleteAllNotifications(c *gin.Context) {
	profileId, err := parseProfileIdHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	err = h.service.DeleteAllNotifications(*profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Success: true})
}
