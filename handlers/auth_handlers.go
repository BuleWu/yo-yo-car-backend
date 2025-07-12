package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"zavrsni/yo-yo-car/models"
	"zavrsni/yo-yo-car/utils"
)

func Login(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// TODO: replace this logic with real authentication
	if user.Email == "email" && user.Password == "password" {
		token, err := utils.GenerateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
	}

	hashedPassword, _ := utils.HashPassword(user.Password)

	user.ID = uuid.New()
	user.Password = hashedPassword

	// TODO: Save user to database and check if already exists
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

/*func UserExists(u *models.User) {

	var user models.User
	result := database.fir
}*/
