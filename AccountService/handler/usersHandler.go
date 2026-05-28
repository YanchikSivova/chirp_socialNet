package handler

import (
	"accountService/models"
	"accountService/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type UsersHandler struct {
	service *service.UsersService
}

func NewUsersHandler(s *service.UsersService) *UsersHandler {
	return &UsersHandler{service: s}
}

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

func parseProfileParam(c *gin.Context) (*uuid.UUID, error) {
	profileIdStr := c.Param("id")
	if profileIdStr == "" {
		return nil, errors.New("no profile_id found")
	}
	profileId, err := uuid.Parse(profileIdStr)
	if err != nil {
		return nil, err
	}
	return &profileId, nil
}

func parseOffsetParam(c *gin.Context) (*int, error) {
	offsetStr := c.Query("offset")
	if offsetStr == "" {
		offset := 0 //Дефолтное значение
		return &offset, nil
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return nil, err
	}
	return &offset, nil
}

func parseLimitParam(c *gin.Context) (*int, error) {
	limitStr := c.Query("limit")
	if limitStr == "" {
		limit := 20 //Дефолтное значение
		return &limit, nil
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, err
	}
	return &limit, nil
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
	profileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if profileId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile_id found"})
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
	err = h.service.FillProfile(*profileId, req.Name, req.Username, req.Avatar, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
	})
}

func (h *UsersHandler) UpdateProfile(c *gin.Context) {
	profileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if profileId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile_id found"})
	}
	var req FillProfileRequest
	if err = c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.UpdateProfile(*profileId, req.Name, req.Username, req.Avatar, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}

type ProfileMeResponse struct {
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
	profileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if profileId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile_id found"})
		return
	}
	profile, err := h.service.GetProfile(*profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ProfileMeResponse{
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

type ProfileResponse struct {
	Profile      ProfileMeResponse `json:"profile"`
	Relationship string            `json:"relationship"`
}

func (h *UsersHandler) GetProfileById(c *gin.Context) {
	meProfileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if meProfileId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile_id found"})
	}
	profileId, err := parseProfileParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if profileId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no profile_id found"})
	}
	profile, err := h.service.GetProfile(*profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	relationship, err := h.service.GetRelationship(*meProfileId, *profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ProfileResponse{
		Profile: ProfileMeResponse{
			ProfileID:         profile.ProfileID,
			Name:              *profile.Name,
			Username:          *profile.Username,
			Avatar:            *profile.Avatar,
			Description:       *profile.Description,
			SubscribersAmount: profile.SubscribersAmount,
			SubscribedAmount:  profile.SubscribedAmount,
			PostsAmount:       profile.PostsAmount,
		},
		Relationship: relationship,
	})
}

func (h *UsersHandler) Follow(c *gin.Context) {
	meProfileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if meProfileId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile_id found"})
	}
	profileId, err := parseProfileParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if profileId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no profile_id found"})
	}
	relationship, err := h.service.GetRelationship(*meProfileId, *profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if relationship == "blocked" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "blocked"})
		return
	}
	err = h.service.Follow(*meProfileId, *profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
	})
}

func (h *UsersHandler) Unfollow(c *gin.Context) {
	meProfileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if meProfileId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile_id found"})
	}
	profileId, err := parseProfileParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if profileId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no profile_id found"})
	}
	err = h.service.Unfollow(*meProfileId, *profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}

func (h *UsersHandler) Block(c *gin.Context) {
	meProfileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if meProfileId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile_id found"})
		return
	}
	profileId, err := parseProfileParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if profileId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no profile_id found"})
		return
	}
	err = h.service.Block(*profileId, *meProfileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}

func (h *UsersHandler) Unblock(c *gin.Context) {
	meProfileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if meProfileId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no profile_id found"})
	}
	bannedProfileId, err := parseProfileParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if bannedProfileId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no profile_id found"})
	}
	err = h.service.Unblock(*bannedProfileId, *meProfileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
	})
}

type ProfileSliceResponse struct {
	Profiles []models.ProfileMinimum `json:"profiles"`
	Limit    int                     `json:"limit"`
	Offset   int                     `json:"offset"`
}

func (h *UsersHandler) GetFollowers(c *gin.Context) {
	profileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	offset, err := parseOffsetParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, err := parseLimitParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	followers, err := h.service.GetFollowers(*profileId, *limit, *offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ProfileSliceResponse{
		Profiles: followers,
		Limit:    *limit,
		Offset:   *offset,
	})
}

func (h *UsersHandler) GetFollowersById(c *gin.Context) {
	profileId, err := parseProfileParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offset, err := parseOffsetParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, err := parseLimitParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	followers, err := h.service.GetFollowers(*profileId, *limit, *offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ProfileSliceResponse{
		Profiles: followers,
		Limit:    *limit,
		Offset:   *offset,
	})
}

func (h *UsersHandler) GetFollowings(c *gin.Context) {
	profileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	offset, err := parseOffsetParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, err := parseLimitParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	followings, err := h.service.GetFollowings(*profileId, *limit, *offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ProfileSliceResponse{
		Profiles: followings,
		Limit:    *limit,
		Offset:   *offset,
	})
}

func (h *UsersHandler) GetFollowingsById(c *gin.Context) {
	profileId, err := parseProfileParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offset, err := parseOffsetParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, err := parseLimitParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	followings, err := h.service.GetFollowings(*profileId, *limit, *offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ProfileSliceResponse{
		Profiles: followings,
		Limit:    *limit,
		Offset:   *offset,
	})
}

func (h *UsersHandler) GetBlacklist(c *gin.Context) {
	profileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	offset, err := parseOffsetParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, err := parseLimitParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	blacklist, err := h.service.GetBlacklist(*profileId, *limit, *offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ProfileSliceResponse{
		Profiles: blacklist,
		Limit:    *limit,
		Offset:   *offset,
	})
}

func (h *UsersHandler) SearchByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	offset, err := parseOffsetParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, err := parseLimitParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profiles, err := h.service.SearchByName(name, *limit, *offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ProfileSliceResponse{
		Profiles: profiles,
		Limit:    *limit,
		Offset:   *offset,
	})
}

func (h *UsersHandler) SearchByUsername(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	offset, err := parseOffsetParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, err := parseLimitParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profiles, err := h.service.SearchByUsername(username, *limit, *offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ProfileSliceResponse{
		Profiles: profiles,
		Limit:    *limit,
		Offset:   *offset,
	})
}

func (h *UsersHandler) UserExists(c *gin.Context) {
	profileId, err := parseProfileParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exists, err := h.service.ProfileExists(*profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ExistsResponse{
		Exists: exists,
	})
}

func (h *UsersHandler) GetProfileMinimum(c *gin.Context) {
	profileId, err := parseProfileParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profile, err := h.service.GetProfileMinimum(*profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.ProfileMinimum{
		ProfileID: profile.ProfileID,
		Name:      profile.Name,
		Username:  profile.Username,
		Avatar:    profile.Avatar,
	})
}

func (h *UsersHandler) GetFollowersID(c *gin.Context) {
	profileID, err := parseProfileParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	followers, err := h.service.GetFollowersID(*profileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Followers{
		Followers: followers,
	})
}

type RelationshipRequest struct {
	ProfileID uuid.UUID `json:"profile_id"`
	AuthorID  uuid.UUID `json:"author_id"`
}

type RelationshipResponse struct {
	Relationship string `json:"relationship"`
}

func (h *UsersHandler) GetRelationships(c *gin.Context) {
	var req RelationshipRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	relationship, err := h.service.GetRelationship(req.ProfileID, req.AuthorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, RelationshipResponse{
		Relationship: relationship,
	})
}

func (h *UsersHandler) BatchProfiles(c *gin.Context) {
	var req models.ProfileIDs
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.Profiles) == 0 {
		c.JSON(http.StatusOK, []models.Profile{})
		return
	}
	var profiles []models.ProfileMinimum
	successCount := 0
	for _, profileID := range req.Profiles {
		profile, err := h.service.GetProfileMinimum(profileID)
		if err != nil {
			continue
		}
		profiles = append(profiles, *profile)
		successCount++
	}
	if successCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no profiles found"})
		return
	}
	c.JSON(http.StatusOK, profiles)
}
