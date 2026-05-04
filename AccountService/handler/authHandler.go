package handler

import (
	"accountService/service"
	"accountService/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type VerifyEmailResponse struct {
	Success bool `json:"success"`
}

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.service.Register(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, RegisterResponse{
		Token:   token,
		Success: true,
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	claims, err := utils.ParseRegisterToken(req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email, ok := claims["email"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
	}
	err = h.service.VerifyEmail(email, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, VerifyEmailResponse{
		Success: true,
	})
}
