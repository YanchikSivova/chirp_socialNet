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
	authGroup := r.Group("/auth")
	{
		authGroup.Any("/*path", proxy.ProxyHandler(accountProxy))
	}

	//PROTECTED routes
	protected := r.Group("/users")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.Any("/*path", proxy.ProxyHandler(accountProxy))

	}
	r.Run(":8000")
}
