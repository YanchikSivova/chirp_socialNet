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
	protected := r.Group("/users")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.Any("/*path", proxy.ProxyHandler(accountProxy))

	}
	r.Run(":8000")
}
