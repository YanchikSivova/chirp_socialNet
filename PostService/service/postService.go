package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log"
	"postService/client"
	"postService/kafka"
	"postService/kafkaEvents"
	"postService/models"
	"postService/repository"
	"time"
	"unicode/utf8"
)

type PostService struct {
	repo          *repository.PostRepository
	accountClient *client.AccountClient
	feedClient    *client.FeedClient
	producer      *kafka.Producer
}

func NewPostService(repo *repository.PostRepository, ac *client.AccountClient, fc *client.FeedClient, pr *kafka.Producer) *PostService {
	return &PostService{repo, ac, fc, pr}
}

func (s *PostService) CreatePost(profileID uuid.UUID, postRequest models.CreatePostRequest) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	exists, err := s.accountClient.ProfileExists(profileID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("profile does not exist")
	}
	contentLen := utf8.RuneCountInString(postRequest.Content)
	if contentLen == 0 || contentLen > 250 {
		return errors.New("invalid post content")
	}
	if len(postRequest.Hashtags) > 3 {
		return errors.New("too many hashtags")
	}
	if len(postRequest.Images) > 3 {
		return errors.New("too many images")
	}
	if postRequest.Status != "published" && postRequest.Status != "draft" {
		return errors.New("invalid post status")
	}
	post := models.Post{
		PostID:    uuid.New(),
		ProfileID: profileID,
		Content:   postRequest.Content,
		Status:    postRequest.Status,
	}
	err = s.repo.CreatePost(ctx, tx, post)
	if err != nil {
		return err
	}
	for i, image := range postRequest.Images {
		newImage := models.Image{
			ImageID:    uuid.New(),
			PostID:     post.PostID,
			Url:        image.Url,
			OrderIndex: i + 1,
		}
		err = s.repo.CreateImage(ctx, tx, newImage)
		if err != nil {
			return err
		}
	}
	for i, hashtag := range postRequest.Hashtags {
		id, err := s.repo.GetHashtagId(ctx, tx, hashtag.HashtagName)
		if err != nil && err != pgx.ErrNoRows {
			return err
		}
		if err == pgx.ErrNoRows {
			id = uuid.New()
			newHashtag := models.Hashtag{
				HashtagID:   id,
				HashtagName: hashtag.HashtagName,
			}
			err = s.repo.CreateHashtag(ctx, tx, newHashtag)
			if err != nil {
				return err
			}
		}
		newPostHashtag := models.PostHashtag{
			PostHashtagID: uuid.New(),
			PostID:        post.PostID,
			HashtagID:     id,
			OrderIndex:    i + 1,
		}
		err = s.repo.CreatePostHashtag(ctx, tx, newPostHashtag)
		if err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	if post.Status == "published" {
		log.Println("started to form kafka event")
		event := kafkaEvents.PostEvent{
			EventID:   uuid.New(),
			ProfileID: profileID,
			PostID:    post.PostID,
			IsRepost:  false,
			CreatedAt: time.Now(),
		}
		err = s.producer.SendPostCreated(ctx, event)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func (s *PostService) UpdatePost(profileID, postID uuid.UUID, postReq models.UpdatePostRequest) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	authorID, err := s.repo.GetAuthorId(ctx, tx, postID)
	if err != nil {
		return err
	}
	if authorID != profileID {
		return errors.New("invalid author id")
	}
	status, err := s.repo.GetPostStatus(ctx, tx, postID)
	if err != nil {
		return err
	}
	if status != "draft" {
		return errors.New("invalid post status")
	}
	if postReq.Content == "" || utf8.RuneCountInString(postReq.Content) > 250 {
		return errors.New("invalid post content")
	}
	err = s.repo.UpdatePostContent(ctx, tx, postID, postReq.Content)
	if err != nil {
		return err
	}
	// Удаляем изображения чтобы перезаписать заново новые
	err = s.repo.DeleteImages(ctx, tx, postID)
	if err != nil {
		return err
	}
	for _, image := range postReq.Images {
		newImage := models.Image{
			ImageID:    uuid.New(),
			PostID:     postID,
			Url:        image.Url,
			OrderIndex: image.OrderIndex,
		}
		err = s.repo.CreateImage(ctx, tx, newImage)
		if err != nil {
			return err
		}
	}
	for _, hashtag := range postReq.Hashtags {
		id, err := s.repo.GetHashtagId(ctx, tx, hashtag.HashtagName)
		if err != nil && err != pgx.ErrNoRows {
			return err
		}
		if err == pgx.ErrNoRows {
			id = uuid.New()
			newHashtag := models.Hashtag{
				HashtagID:   id,
				HashtagName: hashtag.HashtagName,
			}
			err = s.repo.CreateHashtag(ctx, tx, newHashtag)
			if err != nil {
				return err
			}
		}
		newPostHashtag := models.PostHashtag{
			PostHashtagID: uuid.New(),
			PostID:        postID,
			HashtagID:     id,
			OrderIndex:    hashtag.OrderIndex,
		}
		err = s.repo.CreatePostHashtag(ctx, tx, newPostHashtag)
		if err != nil {
			return err
		}
	}
	err = s.repo.UpdateLastEditedTime(ctx, tx, postID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PostService) DeletePost(profileID, postID uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	authorID, err := s.repo.GetAuthorId(ctx, tx, postID)
	if err != nil {
		return err
	}
	if authorID != profileID {
		return errors.New("invalid author id")
	}
	status, err := s.repo.GetPostStatus(ctx, tx, postID)
	if err != nil {
		return err
	}
	err = s.repo.DeletePost(ctx, tx, postID)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	if status == "published" {
		log.Println("started to form kafka event")
		event := kafkaEvents.PostEvent{
			EventID:   uuid.New(),
			ProfileID: profileID,
			PostID:    postID,
			IsRepost:  false,
			CreatedAt: time.Now(),
		}
		err = s.producer.SendPostDeleted(ctx, event)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (s *PostService) GetPost(postID, profileID uuid.UUID) (*models.Post, *models.Profile, []models.ImageList, []models.HashtagList, bool, bool, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, nil, nil, nil, false, false, err
	}
	defer tx.Rollback(ctx)
	post, err := s.repo.GetPost(ctx, tx, postID)
	if err != nil {
		return nil, nil, nil, nil, false, false, err
	}
	if post.Status != "published" && post.ProfileID != profileID {
		return nil, nil, nil, nil, false, false, errors.New("not allowed")
	}
	author, err := s.accountClient.ProfileMinimum(post.ProfileID)
	if err != nil {
		return nil, nil, nil, nil, false, false, err
	}
	images, err := s.repo.GetImages(ctx, tx, postID)
	if err != nil {
		return nil, nil, nil, nil, false, false, err
	}
	hashtags, err := s.repo.GetPostHashtags(ctx, tx, postID)
	if err != nil {
		return nil, nil, nil, nil, false, false, err
	}
	liked, err := s.repo.CheckLike(ctx, tx, postID, profileID)
	if err != nil {
		return nil, nil, nil, nil, false, false, err
	}
	reposted, err := s.repo.CheckRepost(ctx, tx, postID, profileID)
	if err != nil {
		return nil, nil, nil, nil, false, false, err
	}
	return post, author, images, hashtags, liked, reposted, tx.Commit(ctx)
}

func (s *PostService) CreateLike(postID, profileID uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	status, err := s.repo.GetPostStatus(ctx, tx, postID)
	if err != nil {
		return err
	}
	if status != "published" {
		return errors.New("invalid post status")
	}
	err = s.repo.CreateLike(ctx, tx, postID, profileID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PostService) DeleteLike(postID, profileID uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = s.repo.DeleteLike(ctx, tx, postID, profileID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PostService) CreateRepost(postID, profileID uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	status, err := s.repo.GetPostStatus(ctx, tx, postID)
	if err != nil {
		return err
	}
	if status != "published" {
		return errors.New("invalid post status")
	}
	authorId, err := s.repo.GetAuthorId(ctx, tx, postID)
	if err != nil {
		return err
	}
	if authorId == profileID {
		return errors.New("cannot create repost")
	}
	err = s.repo.CreateRepost(ctx, tx, postID, profileID)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	log.Println("repost created event")
	event := kafkaEvents.PostEvent{
		EventID:   uuid.New(),
		ProfileID: profileID,
		PostID:    postID,
		IsRepost:  true,
		CreatedAt: time.Now(),
	}
	err = s.producer.SendPostCreated(ctx, event)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (s *PostService) DeleteRepost(postID, profileID uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = s.repo.DeleteRepost(ctx, tx, postID, profileID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PostService) CreateReport(postID, profileID uuid.UUID, reason string) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	authorId, err := s.repo.GetAuthorId(ctx, tx, postID)
	if err != nil {
		return err
	}
	if authorId == profileID {
		return errors.New("cannot create report")
	}
	err = s.repo.CreateReport(ctx, tx, postID, profileID, reason)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PostService) CreateComment(postID, profileID uuid.UUID, content string, parentCommentIDStr string) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	status, err := s.repo.GetPostStatus(ctx, tx, postID)
	if err != nil {
		return err
	}
	if status != "published" {
		return errors.New("invalid post status")
	}
	if content == "" || utf8.RuneCountInString(content) > 250 {
		return errors.New("invalid content length")
	}
	if parentCommentIDStr == "" {
		err = s.repo.CreateComment(ctx, tx, postID, profileID, content)
		if err != nil {
			return err
		}
	} else {
		parentCommentID, err := uuid.Parse(parentCommentIDStr)
		if err != nil {
			return err
		}
		parentPostID, err := s.repo.GetPostIDFromComment(ctx, tx, parentCommentID)
		if err != nil {
			return err
		}
		if parentPostID != postID {
			return errors.New("comment references to another post")
		}
		err = s.repo.CreateCommentWithParent(ctx, tx, postID, profileID, parentCommentID, content)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (s *PostService) CreateCommentLike(commentID, profileID uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = s.repo.CreateCommentLike(ctx, tx, commentID, profileID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PostService) DeleteCommentLike(commentID, profileID uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = s.repo.DeleteCommentLike(ctx, tx, commentID, profileID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PostService) GetComment(commentID, profileID uuid.UUID) (*models.Comment, *models.Profile, bool, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, nil, false, err
	}
	defer tx.Rollback(ctx)
	comment, err := s.repo.GetComment(ctx, tx, commentID)
	if err != nil {
		return nil, nil, false, err
	}
	authorID, err := s.repo.GetAuthorId(ctx, tx, comment.PostID)
	if err != nil {
		return nil, nil, false, err
	}
	author, err := s.accountClient.ProfileMinimum(authorID)
	if err != nil {
		return nil, nil, false, err
	}
	isLiked, err := s.repo.CheckCommentLike(ctx, tx, commentID, profileID)
	if err != nil {
		return nil, nil, false, err
	}
	return comment, author, isLiked, nil
}

func (s *PostService) DeleteComment(commentID, profileID uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	authorID, err := s.repo.GetProfileIDFromComment(ctx, tx, commentID)
	if err != nil {
		return err
	}
	if authorID != profileID {
		return errors.New("not allowed to delete comment")
	}
	err = s.repo.DeleteComment(ctx, tx, commentID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PostService) PublishPost(postID, profileID uuid.UUID) error {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	authorID, err := s.repo.GetAuthorId(ctx, tx, postID)
	if err != nil {
		return err
	}
	if authorID != profileID {
		return errors.New("not allowed to publish post")
	}
	status, err := s.repo.GetPostStatus(ctx, tx, postID)
	if err != nil {
		return err
	}
	if status == "published" {
		return errors.New("already published")
	}
	err = s.repo.PublishPost(ctx, tx, postID)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	log.Println("started to form kafka event")
	event := kafkaEvents.PostEvent{
		EventID:   uuid.New(),
		ProfileID: profileID,
		PostID:    postID,
	}
	err = s.producer.SendPostCreated(ctx, event)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (s *PostService) GetPosts(authorID, profileID uuid.UUID, status string, offset, limit int) ([]models.PostResponse, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if status != "published" && authorID != profileID {
		return nil, errors.New("not allowed to get posts")
	}
	author, err := s.accountClient.ProfileMinimum(authorID)
	if err != nil {
		return nil, err
	}
	posts, err := s.repo.GetPosts(ctx, tx, authorID, status, offset, limit)
	if err != nil {
		return nil, err
	}
	var postsResp []models.PostResponse
	for _, post := range posts {
		images, err := s.repo.GetImages(ctx, tx, post.PostID)
		if err != nil {
			return nil, err
		}
		hashtags, err := s.repo.GetPostHashtags(ctx, tx, post.PostID)
		if err != nil {
			return nil, err
		}
		var postResp models.PostResponse
		if status != "published" {
			postResp = models.PostResponse{
				Author:         *author,
				PostID:         post.PostID,
				Content:        post.Content,
				Images:         images,
				Hashtags:       hashtags,
				LikesAmount:    post.LikesAmount,
				CommentsAmount: post.CommentsAmount,
				RepostsAmount:  post.RepostsAmount,
				LastEditedAt:   post.LastEditedAt,
			}
		} else {
			isLiked, err := s.repo.CheckLike(ctx, tx, post.PostID, profileID)
			if err != nil {
				return nil, err
			}
			isReposted, err := s.repo.CheckRepost(ctx, tx, post.PostID, profileID)
			if err != nil {
				return nil, err
			}
			postResp = models.PostResponse{
				Author:         *author,
				PostID:         post.PostID,
				Content:        post.Content,
				Images:         images,
				Hashtags:       hashtags,
				LikesAmount:    post.LikesAmount,
				CommentsAmount: post.CommentsAmount,
				RepostsAmount:  post.RepostsAmount,
				PublishedAt:    *post.PublishedAt,
				IsLiked:        isLiked,
				IsReposted:     isReposted,
			}
		}
		postsResp = append(postsResp, postResp)
	}
	return postsResp, tx.Commit(ctx)
}

func (s *PostService) GetComments(profileID, postID uuid.UUID, limit, offset int) ([]models.CommentResponse, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	comments, err := s.repo.GetComments(ctx, tx, postID, limit, offset)
	log.Print(comments)
	if err != nil {
		return nil, err
	}
	var commentsResp []models.CommentResponse
	for _, comment := range comments {
		author, err := s.accountClient.ProfileMinimum(comment.ProfileID)
		if err != nil {
			return nil, err
		}
		isLiked, err := s.repo.CheckCommentLike(ctx, tx, comment.CommentID, comment.ProfileID)
		if err != nil {
			return nil, err
		}
		commentsResp = append(commentsResp, models.CommentResponse{
			Author:      *author,
			CommentID:   comment.CommentID,
			Content:     comment.Content,
			LikesAmount: comment.LikesAmount,
			CreatedAt:   comment.CreatedAt,
			IsLiked:     isLiked,
		})
	}
	return commentsResp, tx.Commit(ctx)
}

func (s *PostService) GetCommentAnswers(profileID, commentID uuid.UUID, limit, offset int) ([]models.CommentResponse, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	comments, err := s.repo.GetCommentAnswers(ctx, tx, commentID, limit, offset)
	if err != nil {
		return nil, err
	}
	var commentsResp []models.CommentResponse
	for _, comment := range comments {
		author, err := s.accountClient.ProfileMinimum(comment.ProfileID)
		if err != nil {
			return nil, err
		}
		isLiked, err := s.repo.CheckCommentLike(ctx, tx, comment.CommentID, comment.ProfileID)
		if err != nil {
			return nil, err
		}
		commentsResp = append(commentsResp, models.CommentResponse{
			Author:      *author,
			CommentID:   comment.CommentID,
			Content:     comment.Content,
			LikesAmount: comment.LikesAmount,
			CreatedAt:   comment.CreatedAt,
			IsLiked:     isLiked,
		})
	}
	return commentsResp, tx.Commit(ctx)
}

func (s *PostService) GetFeedPosts(profileID uuid.UUID) ([]uuid.UUID, error) {
	posts, err := s.feedClient.GetFeedPosts(profileID)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) GetRelationship(profileID, authorID uuid.UUID) error {
	relationship, err := s.accountClient.GetRelationship(profileID, authorID)
	if err != nil {
		return err
	}
	if relationship == "blocked" {
		return errors.New("blocked")
	}
	return nil
}

func (s *PostService) GetPostPreviews(posts models.PostIDs) ([]models.PostPreview, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var postPreviews []models.PostPreview
	for _, post := range posts.Posts {
		content, err := s.repo.GetPostsContent(ctx, tx, post)
		if err != nil {
			content = ""
		}
		preview := models.PostPreview{
			PostID:  post,
			Content: content,
		}
		postPreviews = append(postPreviews, preview)
	}
	return postPreviews, tx.Commit(ctx)
}
