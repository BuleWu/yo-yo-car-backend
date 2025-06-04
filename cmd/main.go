package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"zavrsni/yo-yo-car/db"
	"zavrsni/yo-yo-car/handlers"
	"zavrsni/yo-yo-car/middleware"
	"zavrsni/yo-yo-car/runtimebag"
	"zavrsni/yo-yo-car/shared/constants"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file...")
	}

	var (
		host     = runtimebag.GetEnvString(constants.DatabaseHost, "")
		port     = runtimebag.GetEnvString(constants.DatabasePort, "")
		user     = runtimebag.GetEnvString(constants.DatabaseUser, "")
		password = runtimebag.GetEnvString(constants.DatabasePassword, "")
		database = runtimebag.GetEnvString(constants.DatabaseName, "")
		driver   = runtimebag.GetEnvString(constants.DatabaseDriver, "")
	)

	db.Connect(host, port, user, password, database, driver)

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
