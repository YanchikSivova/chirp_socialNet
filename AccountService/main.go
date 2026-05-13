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
	repoAuth := repository.NewAuthRepository(database)
	serviceAuth := service.NewAuthService(repoAuth)
	handlerAuth := handler.NewAuthHandler(serviceAuth)

	router := gin.Default()

	router.POST("/auth/register", handlerAuth.Register)
	router.POST("/auth/verify-email", handlerAuth.VerifyEmail)
	router.POST("/auth/resend-verification", handlerAuth.ResendVerificationWithToken) //получает email отправляет temporary token
	router.POST("/auth/login", handlerAuth.Login)
	router.POST("/auth/refresh", handlerAuth.Refresh)
	router.POST("/auth/refresh/logout", handlerAuth.Logout)
	router.POST("/auth/logout-all", handlerAuth.LogoutAll)
	router.POST("/auth/change-email/request", handlerAuth.ChangeEmail)
	router.POST("/auth/change-email/confirm", handlerAuth.ChangeEmailConfirm)
	router.POST("/auth/change-email/resend-verification", handlerAuth.ResendVerificationForChangeEmail)
	router.POST("/auth/change-password/request", handlerAuth.ChangePassword)
	router.POST("/auth/change-password/confirm", handlerAuth.ChangePasswordConfirm)
	router.POST("/auth/change-password/resend-verification", handlerAuth.ResendVerificationForChangePassword)
	router.POST("/auth/reset-password/request", handlerAuth.ResetPasswordRequest)
	router.POST("/auth/reset-password/confirm", handlerAuth.ResetPasswordConfirm)

	repoUsers := repository.NewUsersRepository(database)
	serviceUsers := service.NewUsersService(repoUsers)
	handlerUsers := handler.NewUsersHandler(serviceUsers)

	router.POST("/users/check-username", handlerUsers.CheckUsername)
	router.POST("/users/me/profile", handlerUsers.FillProfile)
	router.PATCH("/users/me/profile", handlerUsers.UpdateProfile)
	router.GET("/users/me", handlerUsers.GetProfileMe)

	router.Run(":8080")

}
