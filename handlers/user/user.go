package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mike-pr.com/booleans_as_a_service/models"
	"mike-pr.com/booleans_as_a_service/security"
)

// @Summary Delete a User
// @Description Requires you to have the corresponding key of the user you are deleting.
// @Produce json
// @Param username path string true "username of user to delete"
// @Success 200 {object} models.MessageResponse "User Deleted"
// @Failure 404 {object} models.MessageResponse "User Not Found"
// @Security ApiKeyAuth
// @Router /users/{username} [delete]
// @Tags user
func Delete(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		user := security.Authenticate(c, db, username)
		if user == nil {
			return
		}

		if err := db.Delete(&user).Error; err != nil {
			c.JSON(500, models.MessageResponse{Message: "Failed to delete user"})
			return
		}

		c.JSON(200, models.MessageResponse{Message: "User deleted"})
	}
}
