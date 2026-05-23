package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"postService/client"
	"postService/db"
	"postService/handler"
	"postService/kafka"
	"postService/repository"
	"postService/service"
)

func main() {
	godotenv.Load()
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	postRepository := repository.NewPostRepository(database)
	accountClient := client.NewAccountClient()
	feedClent := client.NewFeedClient()
	postProducer := kafka.NewProducer()
	defer postProducer.Close()
	postService := service.NewPostService(postRepository, accountClient, feedClent, postProducer)
	postHandler := handler.NewPostHandler(postService)

	router := gin.Default()

	router.POST("/posts", postHandler.CreatePost)
	router.PATCH("/posts/:id", postHandler.UpdatePost)
	router.DELETE("/posts/:id", postHandler.DeletePost)
	router.GET("/posts/:id", postHandler.GetPost)
	router.POST("/posts/:id/like", postHandler.CreateLike)
	router.DELETE("/posts/:id/like", postHandler.DeleteLike)
	router.POST("/posts/:id/repost", postHandler.CreateRepost)
	router.DELETE("/posts/:id/repost", postHandler.DeleteRepost)
	router.POST("/posts/:id/report", postHandler.CreateReport)
	router.POST("/posts/:id/comments", postHandler.CreateComment)
	router.POST("/comments/:id/like", postHandler.CreateCommentLike)
	router.DELETE("/comments/:id/like", postHandler.DeleteCommentLike)
	router.GET("/comments/:id", postHandler.GetComment)
	router.DELETE("/comments/:id", postHandler.DeleteComment)
	router.POST("/posts/:id/publish", postHandler.PublishPost)
	router.GET("/users/me/posts/:status", postHandler.GetPostsMe)
	router.GET("/users/:id/posts", postHandler.GetPosts)
	router.GET("/posts/:id/comments", postHandler.GetComments)
	router.GET("/comments/:id/answers", postHandler.GetCommentAnswers)
	router.GET("/feed", postHandler.GetFeed)

	router.Run(":8181")
}
