package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"messageService/models"
	"messageService/util"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *MessageHandler) WS(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	profileID, err := util.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	log.Println("WS connected", profileID)
	h.hub.AddClient(profileID, conn)
	defer func() {
		h.hub.RemoveClient(profileID, conn)
		conn.Close()
	}()
	for {
		var msg models.WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			break
		}
		switch msg.Type {
		case "message":
			err = h.handleMessage(profileID, msg)
			if err != nil {
				log.Println(err)
			}
		case "read":
			err = h.handleRead(profileID, msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (h *MessageHandler) handleMessage(senderID uuid.UUID, msg models.WSMessage) error {
	response, targetProfileID, err := h.service.CreateMessage(
		senderID,
		msg.ConversationID,
		msg.Content,
		msg.ContentType,
	)
	if err != nil {
		return err
	}
	err = h.hub.Send(*targetProfileID, response)
	if err != nil {
		return err
	}
	return nil
}

func (h *MessageHandler) handleRead(profileID uuid.UUID, msg models.WSMessage) error {
	targetProfileID, err := h.service.MarkConversationRead(profileID, msg.ConversationID)
	if err != nil {
		return err
	}
	resp := models.WSMessageResponse{
		Type:           "read",
		ConversationID: msg.ConversationID,
		ProfileID:      profileID,
	}
	err = h.hub.Send(*targetProfileID, resp)
	if err != nil {
		return err
	}
	err = h.hub.Send(profileID, resp)
	if err != nil {
		return err
	}
	return nil
}
