package main

import (
	"accountService/db"
	"accountService/handler"
	"accountService/repository"
	"accountService/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	godotenv.Load()
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewAuthRepository(database)
	service := service.NewAuthService(repo)
	handler := handler.NewAuthHandler(service)

	router := gin.Default()

	router.POST("/register", handler.Register)
	router.POST("/verify-email", handler.VerifyEmail)
	router.POST("/resend-verification", handler.ResendVerificationEmail)
	router.POST("/login", handler.Login)

	router.Run(":8080")

}
