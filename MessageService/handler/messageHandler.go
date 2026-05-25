package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"messageService/hub"
	"messageService/models"
	"messageService/service"
	"net/http"
	"time"
)

type MessageHandler struct {
	service *service.MessageService
	hub     *hub.Hub
}

func NewMessageHandler(service *service.MessageService, h *hub.Hub) *MessageHandler {
	return &MessageHandler{service: service, hub: h}
}

func parseProfileIDHeader(c *gin.Context) (*uuid.UUID, error) {
	profileIdStr := c.GetHeader("X-User-Id")
	if profileIdStr == "" {
		return nil, errors.New("no profile id")
	}
	profileId, err := uuid.Parse(profileIdStr)
	if err != nil {
		return nil, err
	}
	return &profileId, nil
}

func parseIdParam(c *gin.Context) (*uuid.UUID, error) {
	idParam := c.Param("id")
	if idParam == "" {
		return nil, errors.New("no id")
	}
	id, err := uuid.Parse(idParam)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func parseCursorQuery(c *gin.Context) (*time.Time, error) {
	cursorQuery := c.Query("cursor")
	if cursorQuery == "" {
		return nil, nil
	}
	cursorTime, err := time.Parse(time.RFC3339, cursorQuery)
	if err != nil {
		return nil, err
	}
	return &cursorTime, nil
}

type GetConversationIDRequest struct {
	TargetProfileID uuid.UUID `json:"target_profile_id"`
}

type GetConversationIDResponse struct {
	ConversationID uuid.UUID `json:"conversation_id"`
}

func (h *MessageHandler) GetConversationID(c *gin.Context) {
	profileID, err := parseProfileIDHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var req GetConversationIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	conversationID, err := h.service.GetConversationByProfile(*profileID, req.TargetProfileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, GetConversationIDResponse{
		ConversationID: *conversationID,
	})
}

func (h *MessageHandler) GetConversation(c *gin.Context) {
	profileID, err := parseProfileIDHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	conversationID, err := parseIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cursor, err := parseCursorQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	targetProfile, err := h.service.GetProfileByID(*profileID, *conversationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	messages, newCursor, err := h.service.GetConversation(*profileID, *conversationID, cursor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.ConversationResponse{
		ConversationID: *conversationID,
		Profile:        *targetProfile,
		Messages:       messages,
		NextCursor:     newCursor,
	})
}

func (h *MessageHandler) GetConversationList(c *gin.Context) {
	profileID, err := parseProfileIDHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	conversations, err := h.service.GetConversationList(*profileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(conversations) == 0 {
		c.JSON(http.StatusOK, models.ConversationPreviewList{})
		return
	}
	successCount := 0
	for i := range conversations {
		targetProfile, err := h.service.GetProfileByID(*profileID, conversations[i].ConversationID)
		if err != nil {
			continue
		}
		conversations[i].Profile = *targetProfile
		successCount++
	}
	if successCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no conversation found"})
		return
	}
	c.JSON(http.StatusOK, models.ConversationPreviewList{
		Conversations: conversations,
	})
}
