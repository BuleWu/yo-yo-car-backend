package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"zavrsni/yo-yo-car/db"
	"zavrsni/yo-yo-car/handlers"
	"zavrsni/yo-yo-car/middleware"
)

func main() {
	db.Connect()

	r := gin.Default()

	publicRoutes := r.Group("/public")
	{
		publicRoutes.POST("/login", handlers.Login)
		publicRoutes.POST("/register", handlers.Register)
	}

	protectedRoutes := r.Group("/protected")
	protectedRoutes.Use(middleware.AuthenticationMiddleware())
	{
		// Protected routes here
	}

	if err := r.Run(":8080"); err != nil {
		fmt.Errorf("error while trying to run server: %v\n", err)
	}
}
