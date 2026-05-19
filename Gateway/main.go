package main

import (
	"gateway/middleware"
	"gateway/proxy"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()

	accountProxy, err := proxy.NewReverseProxy("http://account-service:8080")
	postProxy, err := proxy.NewReverseProxy("http://post-service:8181")
	if err != nil {
		log.Fatal(err)
	}

	//PUBLIC routes
	publicAuth := r.Group("/auth")
	{
		publicAuth.POST("/register", proxy.ProxyHandler(accountProxy))
		publicAuth.POST("/verify-email", proxy.ProxyHandler(accountProxy))
		publicAuth.POST("/resend-verification", proxy.ProxyHandler(accountProxy))
		publicAuth.POST("/login", proxy.ProxyHandler(accountProxy))
		publicAuth.POST("/refresh", proxy.ProxyHandler(accountProxy))
		publicAuth.POST("/refresh/logout", proxy.ProxyHandler(accountProxy))
		publicAuth.POST("/reset-password", proxy.ProxyHandler(accountProxy))
		publicAuth.POST("/reset-password/*path", proxy.ProxyHandler(accountProxy))
	}

	protectedAuth := r.Group("/auth")
	protectedAuth.Use(middleware.AuthMiddleware())
	{
		protectedAuth.POST("/logout-all", proxy.ProxyHandler(accountProxy))
		protectedAuth.POST("/change-email/*path", proxy.ProxyHandler(accountProxy))
		protectedAuth.POST("/change-password/*path", proxy.ProxyHandler(accountProxy))
		protectedAuth.POST("/me/delete/*path", proxy.ProxyHandler(accountProxy))
	}
	//PROTECTED routes
	protectedUsers := r.Group("/users")
	protectedUsers.Use(middleware.AuthMiddleware())
	{
		protectedUsers.POST("/check-username", proxy.ProxyHandler(accountProxy))
		protectedUsers.Any("/me/profile", proxy.ProxyHandler(accountProxy))
		protectedUsers.GET("/me", proxy.ProxyHandler(accountProxy))
		protectedUsers.GET("/:id", proxy.ProxyHandler(accountProxy))
		protectedUsers.Any("/:id/follow", proxy.ProxyHandler(accountProxy))
		protectedUsers.Any("/:id/block", proxy.ProxyHandler(accountProxy))
		protectedUsers.Any("/:id/followers", proxy.ProxyHandler(accountProxy))
		protectedUsers.Any("/:id/following", proxy.ProxyHandler(accountProxy))
		protectedUsers.GET("/me/followers", proxy.ProxyHandler(accountProxy))
		protectedUsers.GET("/me/following", proxy.ProxyHandler(accountProxy))
		protectedUsers.GET("/me/blacklist", proxy.ProxyHandler(accountProxy))
		protectedUsers.GET("/search/by-name", proxy.ProxyHandler(accountProxy))
		protectedUsers.GET("/search", proxy.ProxyHandler(accountProxy))
		protectedUsers.GET("/me/posts/*path", proxy.ProxyHandler(postProxy))
		protectedUsers.GET("/:id/posts", proxy.ProxyHandler(postProxy))
	}

	protectedPosts := r.Group("/posts")
	protectedPosts.Use(middleware.AuthMiddleware())
	{
		protectedPosts.Any("", proxy.ProxyHandler(postProxy))
		protectedPosts.Any("/*path", proxy.ProxyHandler(postProxy))
	}
	protectedComments := r.Group("/comments")
	protectedComments.Use(middleware.AuthMiddleware())
	{
		protectedComments.Any("/*path", proxy.ProxyHandler(postProxy))
	}
	r.Run(":8000")
}
