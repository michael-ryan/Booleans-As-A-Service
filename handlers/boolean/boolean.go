package boolean

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mike-pr.com/booleans_as_a_service/models"
	"mike-pr.com/booleans_as_a_service/security"
)

// @Summary Get a Boolean
// @Produce json
// @Param username path string true "your username"
// @Param boolean path string true "boolean name"
// @Success 200 {object} models.BooleanResponse "Your Booleans"
// @Failure 404 {object} models.MessageResponse "Boolean Not Found"
// @Security ApiKeyAuth
// @Router /users/{username}/booleans/{boolean} [get]
// @Tags boolean
func Get(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		user := security.Authenticate(c, db, username)
		if user == nil {
			return
		}

		boolname := c.Param("boolean")
		var b models.Boolean
		if err := db.Where("user_id = ? AND name = ?", user.ID, boolname).First(&b).Error; err == gorm.ErrRecordNotFound {
			c.JSON(404, models.MessageResponse{Message: "Boolean does not exist"})
			return
		} else if err != nil {
			c.JSON(500, models.MessageResponse{Message: "Database error"})
			return
		}

		c.JSON(200, models.BooleanResponse{
			Name:  b.Name,
			Value: b.Value,
		})
	}
}

// @Summary Delete a Boolean
// @Produce json
// @Param username path string true "your username"
// @Param boolean path string true "boolean name"
// @Success 200 {object} models.MessageResponse "Boolean Deleted"
// @Failure 404 {object} models.MessageResponse "User Not Found"
// @Security ApiKeyAuth
// @Router /users/{username}/booleans/{boolean} [delete]
// @Tags boolean
func Delete(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		user := security.Authenticate(c, db, username)
		if user == nil {
			return
		}

		boolname := c.Param("boolean")

		b := models.Boolean{
			UserID: user.ID,
			Name:   boolname,
		}

		if err := db.Where("user_id = ? AND name = ?", user.ID, boolname).Delete(&b).Error; err != nil {
			fmt.Printf("%s\n", fmt.Errorf("could not delete boolean: %w", err).Error())
			c.JSON(500, models.MessageResponse{Message: "Failed to delete boolean"})
			return
		}

		c.JSON(200, models.MessageResponse{Message: "Boolean deleted"})
	}
}
