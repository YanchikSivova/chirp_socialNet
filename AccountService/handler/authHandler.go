package handler

import (
	"accountService/service"
	"accountService/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

type AuthRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SuccessWithTokenResponse struct {
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
	c.JSON(http.StatusCreated, SuccessWithTokenResponse{
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
	claims, err := utils.ParseTemporaryToken(req.Token)
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

type EmailRequest struct {
	Email string `json:"email" binding:"required"`
}

type ResendVerificationEmailResponse struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
}

func (h *AuthHandler) ResendVerificationWithToken(c *gin.Context) {
	var req EmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.service.ResendVerificationWithToken(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResendVerificationEmailResponse{
		Token:   token,
		Success: true,
	})
}

type NewEmailRequest struct {
	NewEmail string `json:"new_email" binding:"required"`
}

func (h *AuthHandler) ResendVerificationForChangeEmail(c *gin.Context) {
	var req NewEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profileIdStr := c.GetHeader("X-User-Id")
	if profileIdStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile id found"})
		return
	}
	profileId, err := uuid.Parse(profileIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.ResendVerificationForChangeEmail(profileId, req.NewEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})

}

func (h *AuthHandler) ResendVerificationForChangePassword(c *gin.Context) {
	profileIdStr := c.GetHeader("X-User-Id")
	if profileIdStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile id found"})
		return
	}
	profileId, err := uuid.Parse(profileIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.ResendVerificationForChangePassword(profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
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
	c.SetCookie("refresh_token", refresh, 7*24*60*60, "/", "", false, true)
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

type ChangeEmailRequest struct {
	NewEmail string `json:"new_email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) ChangeEmail(c *gin.Context) {
	var req ChangeEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	err = h.service.ChangeEmail(profileID, req.NewEmail, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}

type VerifyCodeRequest struct {
	Code string `json:"code" binding:"required"`
}

func (h *AuthHandler) ChangeEmailConfirm(c *gin.Context) {
	var req VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	err = h.service.ChangeEmailConfirm(profileID, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	err = h.service.ChangePassword(profileID, req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}

func (h *AuthHandler) ChangePasswordConfirm(c *gin.Context) {
	var req VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	err = h.service.ChangePasswordConfirm(profileID, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}

func (h *AuthHandler) ResetPasswordRequest(c *gin.Context) {
	var req EmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.service.ResetPasswordRequest(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessWithTokenResponse{
		Success: true,
		Token:   token,
	})
}

type ResetPasswordConfirmRequest struct {
	NewPassword string `json:"new_password" binding:"required"`
	Token       string `json:"token" binding:"required"`
	Code        string `json:"code" binding:"required"`
}

func (h *AuthHandler) ResetPasswordConfirm(c *gin.Context) {
	var req ResetPasswordConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.ResetPasswordConfirm(req.Token, req.Code, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}

type PasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) DeleteAccountRequest(c *gin.Context) {
	var req PasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	profileIdStr := c.GetHeader("X-User-Id")
	if profileIdStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile_id provided"})
		return
	}
	profileID, err := uuid.Parse(profileIdStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	err = h.service.DeleteAccountRequest(profileID, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}

func (h *AuthHandler) DeleteAccountConfirm(c *gin.Context) {
	var req VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profileIdStr := c.GetHeader("X-User-Id")
	if profileIdStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile_id provided"})
		return
	}
	profileID, err := uuid.Parse(profileIdStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	err = h.service.DeleteAccountConfirm(profileID, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}
