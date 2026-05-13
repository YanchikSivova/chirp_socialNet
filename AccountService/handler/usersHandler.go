package handler

import (
	"accountService/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type UsersHandler struct {
	service *service.UsersService
}

func NewUsersHandler(s *service.UsersService) *UsersHandler {
	return &UsersHandler{service: s}
}

type UsernameRequest struct {
	Username string `json:"username"`
}

type ExistsResponse struct {
	Exists bool `json:"exists"`
}

func (h *UsersHandler) CheckUsername(c *gin.Context) {
	var req UsernameRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exists, err := h.service.CheckUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ExistsResponse{
		Exists: exists})
}

type FillProfileRequest struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

func (h *UsersHandler) FillProfile(c *gin.Context) {
	profileIdStr := c.GetHeader("X-User-Id")
	if profileIdStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile_id found"})
		return
	}
	profileId, err := uuid.Parse(profileIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req FillProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Username == "" || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or name is empty"})
		return
	}
	err = h.service.FillProfile(profileId, req.Name, req.Username, req.Avatar, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
	})
}

func (h *UsersHandler) UpdateProfile(c *gin.Context) {
	profileIdStr := c.GetHeader("X-User-Id")
	if profileIdStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile_id found"})
		return
	}
	profileId, err := uuid.Parse(profileIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req FillProfileRequest
	if err = c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.UpdateProfile(profileId, req.Name, req.Username, req.Avatar, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}

type ProfileResponse struct {
	ProfileID         uuid.UUID `json:"profile_id"`
	Name              string    `json:"name"`
	Username          string    `json:"username"`
	Avatar            string    `json:"avatar"`
	Description       string    `json:"description"`
	SubscribersAmount int       `json:"subscribers_amount"`
	SubscribedAmount  int       `json:"subscribed_amount"`
	PostsAmount       int       `json:"posts_amount"`
}

func (h *UsersHandler) GetProfileMe(c *gin.Context) {
	profileIdStr := c.GetHeader("X-User-Id")
	if profileIdStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile_id found"})
		return
	}
	profileId, err := uuid.Parse(profileIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profile, err := h.service.GetProfile(profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ProfileResponse{
		ProfileID:         profile.ProfileID,
		Name:              *profile.Name,
		Username:          *profile.Username,
		Avatar:            *profile.Avatar,
		Description:       *profile.Description,
		SubscribersAmount: profile.SubscribersAmount,
		SubscribedAmount:  profile.SubscribedAmount,
		PostsAmount:       profile.PostsAmount,
	})
}
