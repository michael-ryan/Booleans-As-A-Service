package users

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mike-pr.com/booleans_as_a_service/models"
)

// @Summary Create a User
// @Produce json
// @Param username query string true "username of user to create"
// @Success 201 {object} models.KeyResponse "User Created"
// @Failure 409 {object} models.MessageResponse "User Already Exists"
// @Router /users [post]
// @Tags users
func Create(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Query("username")
		if username == "" {
			c.JSON(400, models.MessageResponse{Message: "Username is required"})
			return
		}

		apiKey := generateAPIKey()

		user := models.User{
			Username: username,
			APIKey:   apiKey,
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(409, models.MessageResponse{Message: "User already exists"})
			return
		}

		c.JSON(201, models.KeyResponse{Key: apiKey})
	}
}

// generates a 32 character random string
func generateAPIKey() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		panic("failed to generate API key")
	}
	return hex.EncodeToString(bytes)
}
