package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"postService/client"
	"postService/db"
	"postService/handler"
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
	postService := service.NewPostService(postRepository, accountClient)
	postHandler := handler.NewPostHandler(postService)

	router := gin.Default()

	router.POST("/posts", postHandler.CreatePost)
	
	router.Run(":8181")
}
