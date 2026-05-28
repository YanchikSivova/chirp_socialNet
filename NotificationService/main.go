package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"notificationService/client"
	"notificationService/db"
	"notificationService/handler"
	"notificationService/kafka"
	"notificationService/repository"
	"notificationService/service"
)

func main() {
	godotenv.Load()
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewNotificationRepository(database)
	consumer := kafka.NewConsumer()
	accountClient := client.NewAccountClient()
	notifService := service.NewNotificationService(repo, accountClient)
	notifhandler := handler.NewNotificationHandler(notifService)
	router := gin.Default()
	router.GET("/notifications", notifhandler.GetNotifications)
	router.DELETE("/notifications", notifhandler.DeleteAllNotifications)
	router.DELETE("/notifications/:id", notifhandler.DeleteNotification)
	go consumer.ConsumeNotification(repo)
	router.Run(":8484")
}
