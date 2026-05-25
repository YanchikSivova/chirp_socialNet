package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"postService/models"
	"postService/service"
	"strconv"
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

func parseOffsetQuery(c *gin.Context) (*int, error) {
	offsetStr := c.Query("offset")
	if offsetStr == "" {
		offset := 0
		return &offset, nil
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return nil, err
	}
	return &offset, nil
}

func parseLimitQuery(c *gin.Context) (*int, error) {
	limitStr := c.Query("limit")
	if limitStr == "" {
		limit := 20
		return &limit, nil
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, err
	}
	return &limit, nil
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
	if post.Status == "published" {
		c.JSON(http.StatusOK, models.PostResponse{
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
	} else {
		c.JSON(http.StatusOK, models.PostResponse{
			Author:         *profile,
			PostID:         post.PostID,
			Content:        post.Content,
			Images:         images,
			Hashtags:       hashtags,
			LikesAmount:    post.LikesAmount,
			CommentsAmount: post.CommentsAmount,
			RepostsAmount:  post.RepostsAmount,
			LastEditedAt:   post.LastEditedAt,
			IsLiked:        liked,
			IsReposted:     reposted,
		})
	}
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
	Content         string `json:"content"`
	ParentCommentID string `json:"parent_comment_id"`
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
	err = h.service.CreateComment(*postId, *profileId, req.Content, req.ParentCommentID)
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

func (h *PostHandler) GetComment(c *gin.Context) {
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
	comment, author, isLiked, err := h.service.GetComment(*commentId, *profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.CommentResponse{
		Author:      *author,
		CommentID:   comment.CommentID,
		Content:     comment.Content,
		LikesAmount: comment.LikesAmount,
		CreatedAt:   comment.CreatedAt,
		IsLiked:     isLiked,
	})
}

func (h *PostHandler) DeleteComment(c *gin.Context) {
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
	err = h.service.DeleteComment(*commentId, *profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Success: true})
}

func (h *PostHandler) PublishPost(c *gin.Context) {
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
	err = h.service.PublishPost(*postId, *profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Success: true})
}

func (h *PostHandler) GetPostsMe(c *gin.Context) {
	profileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	status := c.Param("status")
	if status != "published" && status != "draft" && status != "banned" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	offset, err := parseOffsetQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	limit, err := parseLimitQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	posts, err := h.service.GetPosts(*profileId, *profileId, status, *offset, *limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Posts{Posts: posts})
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	profileId, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	authorId, err := parseIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offset, err := parseOffsetQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	limit, err := parseLimitQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	posts, err := h.service.GetPosts(*authorId, *profileId, "published", *offset, *limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Posts{Posts: posts})
}

type Comments struct {
	Comments []models.CommentResponse `json:"comments"`
}

func (h *PostHandler) GetComments(c *gin.Context) {
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
	offset, err := parseOffsetQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	limit, err := parseLimitQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comments, err := h.service.GetComments(*profileId, *postId, *limit, *offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, Comments{Comments: comments})
}

func (h *PostHandler) GetCommentAnswers(c *gin.Context) {
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
	offset, err := parseOffsetQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	limit, err := parseLimitQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comments, err := h.service.GetCommentAnswers(*profileId, *commentId, *limit, *offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, Comments{Comments: comments})
}

type FeedRequest struct {
	ProfileID uuid.UUID `json:"profileId"`
	PostsID   []string  `json:"posts_id"`
}

func (h *PostHandler) GetFeed(c *gin.Context) {
	profileID, err := parseProfileHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	posts, err := h.service.GetFeedPosts(*profileID)
	var postsResp models.Posts
	successCount := 0
	for _, postID := range posts {
		post, author, images, hashtags, liked, reposted, err := h.service.GetPost(postID, *profileID)
		if err != nil {
			log.Printf("Failed to get post %v", postID)
			continue
		}
		if author.ProfileID == *profileID {
			continue
		}
		if err = h.service.GetRelationship(*profileID, author.ProfileID); err != nil {
			log.Printf("Invalid relationship: %v", err)
			continue
		}
		postResp := models.PostResponse{
			Author:         *author,
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
		}
		postsResp.Posts = append(postsResp.Posts, postResp)
		successCount++
	}
	if successCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no posts found"})
		return
	}
	c.JSON(http.StatusOK, models.Posts{
		Posts: postsResp.Posts,
	})
}

func (h *PostHandler) GetPostPreviews(c *gin.Context) {
	var req models.PostIDs
	if c.ShouldBind(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	previews, err := h.service.GetPostPreviews(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.PostPreviews{
		Posts: previews,
	})
}
