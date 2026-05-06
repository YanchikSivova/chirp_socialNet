package handler

import (
	"accountService/service"
	"accountService/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type AuthRequest struct {
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

type SuccessResponse struct {
	Success bool `json:"success"`
}

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req AuthRequest

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
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}

type ResendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required"`
}

type ResendVerificationEmailResponse struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
}

func (h *AuthHandler) ResendVerificationEmail(c *gin.Context) {
	var req ResendVerificationEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.service.ResetVerification(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResendVerificationEmailResponse{
		Token:   token,
		Success: true,
	})
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	access, refresh, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		if err == errors.New("Email not verified") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("refresh_token", refresh, 7*24*60*60, "/auth/refresh", "", false, true)
	c.JSON(http.StatusOK, LoginResponse{
		AccessToken: access,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	oldRefreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		if err == http.ErrNoCookie {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	access, refresh, err := h.service.AuthRefresh(oldRefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("refresh_token", refresh, 7*24*60*60, "/auth/refresh", "", false, true)
	c.JSON(http.StatusOK, LoginResponse{
		AccessToken: access,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		if err == http.ErrNoCookie {
			c.JSON(http.StatusOK, SuccessResponse{
				Success: true,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.Logout(refreshToken)
	if err != nil {
		if err == errors.New("failed to logout") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}

func (h *AuthHandler) LogoutAll(c *gin.Context) {
	profileIdStr := c.GetHeader("X-User-Id")
	if profileIdStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile_id provided"})
		return
	}
	profileID, err := uuid.Parse(profileIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.LogoutAll(profileID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}
