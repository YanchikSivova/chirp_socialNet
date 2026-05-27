package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log"
	"messageService/client"
	"messageService/models"
	"messageService/repository"
	"time"
)

type MessageService struct {
	repo          *repository.MessageRepository
	accountClient *client.AccountClient
	postClient    *client.PostClient
}

func NewMessageService(repo *repository.MessageRepository, ac *client.AccountClient, pc *client.PostClient) *MessageService {
	return &MessageService{repo: repo, accountClient: ac, postClient: pc}
}

func (s *MessageService) GetConversationByProfile(meProfileID, targetProfileID uuid.UUID) (*uuid.UUID, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	conversationID, err := s.repo.GetConversationID(ctx, tx, meProfileID, targetProfileID)
	if err != nil {
		if err == pgx.ErrNoRows {
			conversationID = uuid.New()
			err = s.repo.CreateConversation(ctx, tx, meProfileID, targetProfileID, conversationID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &conversationID, tx.Commit(ctx)
}

func (s *MessageService) GetProfileByID(meProfileID, conversationID uuid.UUID) (*models.Profile, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	targetProfileID, err := s.repo.GetTargetProfileID(ctx, tx, meProfileID, conversationID)
	if err != nil {
		return nil, err
	}
	log.Println(targetProfileID)
	profile, err := s.accountClient.ProfileMinimum(*targetProfileID)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *MessageService) GetConversation(meProfileID, conversationID uuid.UUID, cursor *time.Time) ([]models.MessageResponse, *time.Time, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)
	messages, err := s.repo.GetMessages(ctx, tx, conversationID, cursor)
	if err != nil {
		log.Println("Get messages err:", err)
		return nil, nil, err
	}
	length := len(messages)
	if length == 0 {
		return []models.MessageResponse{}, nil, tx.Commit(ctx)
	}
	newestMessageID := messages[0].MessageID
	err = s.repo.MarkConversationRead(ctx, tx, conversationID, meProfileID, newestMessageID)
	if err != nil {
		log.Println("Mark conversation read err", err)
		return nil, nil, err
	}
	var posts models.PostIDs
	for _, message := range messages {
		if message.Type == "post" {
			postID, err := uuid.Parse(message.Content)
			log.Println("If message type post err:", err)
			if err != nil {
				continue
			}
			posts.Posts = append(posts.Posts, postID)
		}
	}
	var postPreviews *models.PostPreviews
	if len(posts.Posts) > 0 {
		postPreviews, err = s.postClient.BatchPosts(posts)
		if err != nil {
			log.Println("Batch posts err", err)
			return nil, nil, err
		}
	} else {
		postPreviews = &models.PostPreviews{}
	}

	previewMap := make(map[uuid.UUID]models.PostPreview)
	for _, preview := range postPreviews.Posts {
		previewMap[preview.PostID] = preview
	}
	var resp []models.MessageResponse
	for _, message := range messages {
		m := models.MessageResponse{
			MessageID: message.MessageID,
			SenderID:  message.SenderID,
			Content:   message.Content,
			Type:      message.Type,
			CreatedAt: message.CreatedAt,
		}
		if message.SenderID == meProfileID {
			m.IsMine = true
		}
		if message.Type == "post" {
			postID, err := uuid.Parse(message.Content)

			if err != nil {
				log.Println("message type post", err)
				resp = append(resp, m)
				continue
			}
			preview, ok := previewMap[postID]
			if !ok {
				resp = append(resp, m)
				continue
			}
			m.Post = &preview
		}
		resp = append(resp, m)
	}

	nextCursor := &messages[len(messages)-1].CreatedAt
	return resp, nextCursor, tx.Commit(ctx)
}

func (s *MessageService) GetConversationList(meProfileID uuid.UUID) ([]models.ConversationPreview, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	conversations, err := s.repo.GetConversationList(ctx, tx, meProfileID)
	if err != nil {
		log.Println("GetConversationList failed:", err)
		return nil, err
	}
	for i := range conversations {
		unread, err := s.repo.GetUnreadCount(ctx, tx, conversations[i].ConversationID, meProfileID, conversations[i].LastReadMessageID)
		log.Println("GetUnreadCount failed:", err)
		if err != nil {
			return nil, err
		}
		conversations[i].UnreadCount = unread
	}
	return conversations, tx.Commit(ctx)
}

func (s *MessageService) CreateMessage(senderID, conversationID uuid.UUID, content string, messageType string) (*models.WSMessageResponse, *uuid.UUID, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	targetProfileID, err := s.repo.GetTargetProfileID(ctx, tx, senderID, conversationID)
	if err != nil {
		return nil, nil, err
	}
	messageID := uuid.New()
	err = s.repo.CreateMessage(ctx, tx, messageID, conversationID, senderID, content, messageType)
	if err != nil {
		return nil, nil, err
	}
	err = s.repo.UpdateLastMessageCreatedAt(ctx, tx, conversationID, time.Now())
	if err != nil {
		return nil, nil, err
	}
	resp := &models.WSMessageResponse{
		MessageID:      messageID,
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		ContentType:    messageType,
		CreatedAt:      time.Now(),
	}
	return resp, targetProfileID, tx.Commit(ctx)
}

func (s *MessageService) MarkConversationRead(profileID, conversationID uuid.UUID) (*uuid.UUID, error) {
	ctx := context.Background()
	tx, err := s.repo.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	targetProfileID, err := s.repo.GetTargetProfileID(ctx, tx, profileID, conversationID)
	if err != nil {
		return nil, err
	}
	lastMessageID, err := s.repo.GetLastMessageID(ctx, tx, conversationID)
	if err != nil {
		return nil, err
	}
	err = s.repo.MarkConversationRead(ctx, tx, conversationID, profileID, lastMessageID)
	if err != nil {
		return nil, err
	}
	return targetProfileID, tx.Commit(ctx)
}
