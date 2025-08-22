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
	"zavrsni/yo-yo-car/email"
	"zavrsni/yo-yo-car/firebase"
	"zavrsni/yo-yo-car/middleware"
	"zavrsni/yo-yo-car/pusher"
	"zavrsni/yo-yo-car/repositories"
	"zavrsni/yo-yo-car/runtimebag"
	"zavrsni/yo-yo-car/shared/constants"
)

var (
	userController        *controllers.User
	authController        *controllers.Auth
	rideController        *controllers.Ride
	ratingController      *controllers.Rating
	reservationController *controllers.Reservation
	chatController        *controllers.Chat
	messageController     *controllers.Message
	conn                  *database.Connection
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file...")
	}

	// database connection setup
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

	// application layer init
	userApplication := applications.NewUserApplication(
		repositories.NewUserRepository(conn),
		repositories.NewReservationRepository(conn),
	)

	rideApplication := applications.NewRideApplication(
		repositories.NewRideRepository(conn),
		repositories.NewUserRepository(conn),
		repositories.NewReservationRepository(conn),
	)

	ratingApplication := applications.NewRatingApplication(
		repositories.NewRatingRepository(conn),
		repositories.NewRideRepository(conn),
		repositories.NewUserRepository(conn),
	)

	reservationApplication := applications.NewReservationApplication(
		repositories.NewReservationRepository(conn),
		repositories.NewUserRepository(conn),
		repositories.NewRideRepository(conn),
	)

	chatApplication := applications.NewChatApplication(
		repositories.NewChatRepository(conn),
		repositories.NewUserRepository(conn),
		repositories.NewMessageRepository(conn),
		pusher.NewPusherService(runtimebag.GetEnvString("PUSHER_APP_ID", ""), runtimebag.GetEnvString("PUSHER_KEY", ""), runtimebag.GetEnvString("PUSHER_SECRET", ""), runtimebag.GetEnvString("PUSHER_CLUSTER", "eu"), true),
		repositories.NewRideRepository(conn),
	)

	messageApplication := applications.NewMessageApplication(
		repositories.NewMessageRepository(conn),
	)

	// controller init
	userController = controllers.NewUserController(
		userApplication,
	)

	authController = controllers.NewAuthController(
		userApplication,
	)

	rideController = controllers.NewRideController(
		rideApplication,
	)

	ratingController = controllers.NewRatingController(
		ratingApplication,
	)

	reservationController = controllers.NewReservationController(
		reservationApplication,
	)

	chatController = controllers.NewChatController(
		chatApplication,
	)

	messageController = controllers.NewMessageController(
		messageApplication,
	)

	// route setup
	r := gin.Default()

	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}

	r.Use(cors.New(config))

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
		authRoutes.GET("/google/login", authController.OauthGoogleLogin)
		authRoutes.GET("/google/callback", authController.OauthGoogleCallback)
	}

	apiRoutes := r.Group("/api")
	apiRoutes.GET("/reservations/:id/confirm", reservationController.ConfirmReservation)
	apiRoutes.GET("/reservations/:id/decline", reservationController.DeclineReservation)

	apiRoutes.Use(middleware.AuthenticationMiddleware())
	{
		/*user APIs*/
		apiRoutes.GET("/users/:id", userController.GetUser)
		apiRoutes.POST("/users", userController.CreateUser)
		apiRoutes.PUT("/users/:id", userController.UpdateUser)
		apiRoutes.POST("/users/change-password", userController.ChangePassword)
		apiRoutes.DELETE("/users/:id", userController.DeleteUser)
		apiRoutes.POST("/users/:id/profile-picture", userController.UploadProfilePicture)
		apiRoutes.GET("/users/:id/reservations", userController.GetUserReservations)

		/*ride APIs*/
		apiRoutes.GET("/rides", rideController.GetRides)
		apiRoutes.GET("/user/:userId/rides", rideController.GetUserRides)
		apiRoutes.GET("/rides/:id", rideController.GetRideById)
		apiRoutes.POST("/rides", rideController.CreateRide)
		apiRoutes.PUT("/rides/:id", rideController.UpdateRide)
		apiRoutes.DELETE("/rides/:id", rideController.DeleteRide)
		apiRoutes.GET("/rides/search", rideController.SearchRides)
		apiRoutes.GET("/rides/:id/reservations", rideController.GetRideReservations)
		apiRoutes.PATCH("/rides/:id/finish", rideController.FinishRide)
		/*apiRoutes.GET("/rides/:id/cancel", rideController.CancelRide)*/

		/*rating APIs*/
		apiRoutes.GET("/ratings", ratingController.GetRatings)
		apiRoutes.GET("/ratings/:id", ratingController.GetRatingById)
		apiRoutes.GET("/ratings/users/:userId", ratingController.GetRatingsByUserId)
		apiRoutes.POST("/ratings", ratingController.CreateRating)
		apiRoutes.PUT("/ratings/:id", ratingController.UpdateRating)
		apiRoutes.DELETE("/ratings/:id", ratingController.DeleteRating)

		/*reservation APIs*/
		apiRoutes.GET("/reservations", reservationController.GetAllReservations)
		apiRoutes.GET("/reservations/:id", reservationController.GetReservationById)
		apiRoutes.POST("/reservations", reservationController.CreateReservation)
		apiRoutes.PUT("/reservations/:id", reservationController.UpdateReservation)
		apiRoutes.DELETE("/reservations/:id", reservationController.DeleteReservation)

		/*chat APIs*/
		apiRoutes.GET("/chats", chatController.GetUserChats)
		apiRoutes.GET("/chats/:id", chatController.GetChat)
		apiRoutes.POST("/chats", chatController.CreateChat)
		/*apiRoutes.PUT("/chats/:id", chatController.UpdateChat)*/
		apiRoutes.DELETE("/chats/:id", chatController.DeleteChat)
		apiRoutes.GET("/chats/:id/messages", chatController.GetChatMessages)
		apiRoutes.POST("/chats/:id/messages", chatController.SendMessage)

		/*message APIs*/
		apiRoutes.PUT("/messages/:id", messageController.UpdateMessage)

	}

	ctx := context.Background()

	// firebase init
	credentialsFile := runtimebag.GetEnvString("GOOGLE_APPLICATION_CREDENTIALS", "")
	projectID := runtimebag.GetEnvString("FIREBASE_PROJECT_ID", "")
	storageBucket := runtimebag.GetEnvString("FIREBASE_STORAGE_BUCKET", "")

	firebase.InitFirebase(ctx, credentialsFile, projectID, storageBucket)
	defer firebase.Client.Close()

	// mailtrap init
	if err = email.InitEmail(); err != nil {
		log.Println("Error in email: ", err)
	}

	if err = r.Run("localhost:8080"); err != nil {
		log.Fatalf("error while trying to run server: %v\n", err)
	}
}

//func buildDependencies() {}

/*func httpRouter() {}*/
