package security

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mike-pr.com/booleans_as_a_service/models"
)

// Returns the database model of the user if the username is paired with a correct API key header, nil otherwise.
// If authentication is unsuccessful, this function will write appropriate error codes and messages to the gin context.
// You MUST NOT continue to process a request if this function returns false.
func Authenticate(c *gin.Context, db *gorm.DB, username string) *models.User {
	if username == "" {
		c.JSON(400, models.MessageResponse{Message: "Username is required"})
		return nil
	}

	apiKey := c.GetHeader("X-API-KEY")
	if apiKey == "" {
		c.JSON(401, models.MessageResponse{Message: "Authentication required"})
		return nil
	}

	var user models.User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, models.MessageResponse{Message: fmt.Sprintf("User not found: %v", username)})
		} else {
			c.JSON(500, models.MessageResponse{Message: "Database error"})
		}
		return nil
	}

	if user.APIKey != apiKey {
		c.JSON(403, models.MessageResponse{Message: "API key does not match"})
		return nil
	}

	return &user
}
