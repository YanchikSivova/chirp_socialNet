package main

import (
	"context"
	"feedService/client"
	"feedService/handler"
	"feedService/kafka"
	"feedService/redisDB"
	"feedService/repository"
	"feedService/service"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisClient := redisDB.NewRedisClient()
	err := redisClient.Ping(ctx).Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to Redis")
	repo := repository.NewFeedRepository(redisClient)
	ac := client.NewAccountClient()
	feedService := service.NewFeedService(repo, ac)
	consumer := kafka.NewConsumer(feedService)
	router := gin.Default()
	feedHandler := handler.NewFeedHandler(feedService)

	router.GET("/internal/feed/:id", feedHandler.GetFeed)

	go consumer.ConsumePostCreated()
	router.Run(":8282")
}
