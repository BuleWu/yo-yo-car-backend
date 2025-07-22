package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"time"
	"zavrsni/yo-yo-car/applications"
	"zavrsni/yo-yo-car/controllers"
	"zavrsni/yo-yo-car/database"
	"zavrsni/yo-yo-car/firebase"
	"zavrsni/yo-yo-car/repositories"
	"zavrsni/yo-yo-car/runtimebag"
	"zavrsni/yo-yo-car/shared/constants"
)

var (
	userController *controllers.User
	authController *controllers.Auth
	rideController *controllers.Ride
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

	rideApplication := applications.NewRideApplication(
		repositories.NewRideRepository,
	)

	userController = controllers.NewUserController(
		userApplication,
	)

	authController = controllers.NewAuthController(
		userApplication,
	)

	rideController = controllers.NewRideController(
		rideApplication,
	)

	r := gin.Default()

	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}

	r.Use(cors.New(config))

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

	apiRoutes := r.Group("/api") /*TODO: add auth middleware*/
	{
		/*user APIs*/
		apiRoutes.GET("/user/:id", userController.GetUser)
		apiRoutes.POST("/user", userController.CreateUser)
		apiRoutes.DELETE("/user/:id", userController.DeleteUser)
		apiRoutes.POST("/user/:id/profile-picture", userController.UploadProfilePicture)

		/*ride APIs*/
		apiRoutes.POST("/ride", rideController.CreateRide)

		/*rating APIs*/

		/*conversation APIs*/

		/*message APIs*/
	}

	ctx := context.Background()

	credentialsFile := runtimebag.GetEnvString("GOOGLE_APPLICATION_CREDENTIALS", "")
	projectID := runtimebag.GetEnvString("FIREBASE_PROJECT_ID", "")
	storageBucket := runtimebag.GetEnvString("FIREBASE_STORAGE_BUCKET", "")

	firebase.InitFirebase(ctx, credentialsFile, projectID, storageBucket)
	defer firebase.Client.Close()

	if err := r.Run("localhost:8080"); err != nil {
		log.Fatalf("error while trying to run server: %v\n", err)
	}
}

//func buildDependencies() {}

/*func httpRouter() {}*/
