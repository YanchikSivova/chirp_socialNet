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

	router.POST("/auth/register", handler.Register)
	router.POST("/auth/verify-email", handler.VerifyEmail)
	router.POST("/auth/resend-verification", handler.ResendVerificationEmail)
	router.POST("/auth/login", handler.Login)
	router.POST("/auth/refresh", handler.Refresh)
	router.POST("/auth/refresh/logout", handler.Logout)
	router.POST("/auth/logout-all", handler.LogoutAll)

	router.Run(":8080")

}
