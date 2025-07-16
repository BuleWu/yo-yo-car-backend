package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"zavrsni/yo-yo-car/applications"
	"zavrsni/yo-yo-car/controllers"
	"zavrsni/yo-yo-car/database"
	"zavrsni/yo-yo-car/repositories"
	"zavrsni/yo-yo-car/runtimebag"
	"zavrsni/yo-yo-car/shared/constants"
)

var (
	userController *controllers.User
	authController *controllers.Auth
	conn           *database.Connection
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file...")
	}
	conn, err = database.NewConnection(
		runtimebag.GetEnvString(constants.DatabaseHost, ""),
		runtimebag.GetEnvString(constants.DatabasePort, ""),
		runtimebag.GetEnvString(constants.DatabaseUser, ""),
		runtimebag.GetEnvString(constants.DatabasePassword, ""),
		runtimebag.GetEnvString(constants.DatabaseName, ""),
	)

	if err != nil {
		panic(err)
	}

	userApplication := applications.NewUserApplication(
		repositories.NewUserRepository(conn),
	)

	userController = controllers.NewUserController(
		userApplication,
	)

	authController := controllers.NewAuthController(
		userApplication,
	)

	r := gin.Default()

	/*publicRoutes := r.Group("/public")
	{
		publicRoutes.POST("/login", handlers.Login)
		publicRoutes.POST("/register", handlers.Register)
		//publicRoutes.POST("/user", userController.CreateUser)
	}*/

	/*protectedRoutes := r.Group("/protected")
	protectedRoutes.Use(middleware.AuthenticationMiddleware())
	{
		// Protected routes here
	}*/
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
		authRoutes.GET("/google/login", authController.OauthGoogleLogin)
		authRoutes.GET("/google/callback", authController.OauthGoogleCallback)
	}

	/*TODO: maybe add a profile endpoint*/

	apiRoutes := r.Group("/api")
	{
		apiRoutes.GET("/user/:id", userController.GetUser)
		apiRoutes.POST("/user", userController.CreateUser)
		apiRoutes.DELETE("/user/:id", userController.DeleteUser)
	}

	if err := r.Run("localhost:8080"); err != nil {
		log.Fatalf("error while trying to run server: %v\n", err)
	}
}

//func buildDependencies() {}

/*func httpRouter() {}*/
