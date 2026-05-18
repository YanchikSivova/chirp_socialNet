package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"postService/models"
	"postService/service"
	"time"
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

func parseIdParam(c *gin.Context) (*uuid.UUID, error) {
	postIdStr := c.Param("id")
	if postIdStr == "" {
		return nil, errors.New("no id found")
	}
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		return nil, err
	}
	return &postId, nil
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

func (h *PostHandler) UpdatePost(c *gin.Context) {
	profileID, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	postId, err := parseIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req models.UpdatePostRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.UpdatePost(*profileID, *postId, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Success: true})
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	profileID, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	postId, err := parseIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.DeletePost(*profileID, *postId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Success: true})
}

type PostResponse struct {
	Author         models.Profile       `json:"author"`
	PostID         uuid.UUID            `json:"post_id"`
	Content        string               `json:"content"`
	Images         []models.ImageList   `json:"images"`
	Hashtags       []models.HashtagList `json:"hashtags"`
	LikesAmount    int                  `json:"likes_amount"`
	CommentsAmount int                  `json:"comments_amount"`
	RepostsAmount  int                  `json:"reposts_amount"`
	PublishedAt    time.Time            `json:"published_at"`
	IsLiked        bool                 `json:"is_liked"`
	IsReposted     bool                 `json:"is_reposted"`
}

func (h *PostHandler) GetPost(c *gin.Context) {
	profileID, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	postId, err := parseIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post, profile, images, hashtags, liked, reposted, err := h.service.GetPost(*postId, *profileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, PostResponse{
		Author:         *profile,
		PostID:         post.PostID,
		Content:        post.Content,
		Images:         images,
		Hashtags:       hashtags,
		LikesAmount:    post.LikesAmount,
		CommentsAmount: post.CommentsAmount,
		RepostsAmount:  post.RepostsAmount,
		PublishedAt:    *post.PublishedAt,
		IsLiked:        liked,
		IsReposted:     reposted,
	})
}

func (h *PostHandler) CreateLike(c *gin.Context) {
	profileID, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	postId, err := parseIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.CreateLike(*postId, *profileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, SuccessResponse{Success: true})
}

func (h *PostHandler) DeleteLike(c *gin.Context) {
	profileID, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	postId, err := parseIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.DeleteLike(*postId, *profileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Success: true})
}

func (h *PostHandler) CreateRepost(c *gin.Context) {
	profileID, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	postId, err := parseIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.CreateRepost(*postId, *profileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, SuccessResponse{Success: true})
}

func (h *PostHandler) DeleteRepost(c *gin.Context) {
	profileID, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	postId, err := parseIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.DeleteRepost(*postId, *profileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Success: true})
}

type ReportRequest struct {
	Reason string `json:"reason"`
}

func (h *PostHandler) CreateReport(c *gin.Context) {
	profileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	postId, err := parseIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req ReportRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.CreateReport(*postId, *profileId, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, SuccessResponse{Success: true})
}

type CommentRequest struct {
	Content         string    `json:"content"`
	ParentCommentID uuid.UUID `json:"parent_comment_id"`
}

func (h *PostHandler) CreateComment(c *gin.Context) {
	profileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	postId, err := parseIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req CommentRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.CreateComment(*postId, *profileId, req.Content, &req.ParentCommentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, SuccessResponse{Success: true})
}

func (h *PostHandler) CreateCommentLike(c *gin.Context) {
	profileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	commentId, err := parseIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.CreateCommentLike(*commentId, *profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, SuccessResponse{Success: true})
}

func (h *PostHandler) DeleteCommentLike(c *gin.Context) {
	profileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	commentId, err := parseIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.DeleteCommentLike(*commentId, *profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Success: true})
}
