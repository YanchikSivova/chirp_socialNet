package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"messageService/client"
	"messageService/db"
	"messageService/handler"
	"messageService/hub"
	"messageService/repository"
	"messageService/service"
)

func main() {
	godotenv.Load()
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewMessageRepository(database)
	ac := client.NewAccountClient()
	pc := client.NewPostClient()
	wsHub := hub.NewHub()
	messService := service.NewMessageService(repo, ac, pc)
	messHandler := handler.NewMessageHandler(messService, wsHub)
	router := gin.Default()

	router.POST("/conversations", messHandler.GetConversationID)
	router.GET("/conversations/:id", messHandler.GetConversation)
	router.GET("/conversations", messHandler.GetConversationList)
	router.GET("/ws", messHandler.WS)

	router.Run(":8383")
}
