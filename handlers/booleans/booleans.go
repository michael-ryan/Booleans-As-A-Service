package booleans

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mike-pr.com/booleans_as_a_service/models"
	"mike-pr.com/booleans_as_a_service/security"
)

// @Summary List your Booleans
// @Produce json
// @Param username path string true "your username"
// @Success 200 {object} models.BooleansResponse "Your Booleans"
// @Security ApiKeyAuth
// @Router /users/{username}/booleans [get]
// @Tags booleans
func Get(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		ok, user := security.Authenticate(c, db, username)
		if !ok {
			return
		}

		// Get all booleans for this user
		var booleans []models.Boolean
		if err := db.Where("user_id = ?", user.ID).Order("name").Find(&booleans).Error; err != nil {
			c.JSON(500, models.MessageResponse{Message: "Failed to fetch booleans"})
			return
		}

		booleanResponses := make([]models.BooleanResponse, 0)
		for _, b := range booleans {
			booleanResponses = append(booleanResponses, models.BooleanResponse{Name: b.Name, Value: b.Value})
		}

		c.JSON(200, models.BooleansResponse{Booleans: booleanResponses})
	}
}

// @Summary Create a Boolean
// @Produce json
// @Param username path string true "your username"
// @Param name query string true "name of boolean to create"
// @Param initialValue query boolean false "initial value of the boolean (defaults to false)"
// @Success 201 {object} models.KeyResponse "Boolean Created"
// @Failure 409 {object} models.MessageResponse "Boolean Already Exists"
// @Security ApiKeyAuth
// @Router /users/{username}/booleans [post]
// @Tags booleans
func Create(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		ok, user := security.Authenticate(c, db, username)

		if !ok {
			return
		}

		name := c.Query("name")
		if name == "" {
			c.JSON(400, models.MessageResponse{Message: "Missing name for boolean"})
			return
		}

		// Check if boolean already exists for this user
		var existing models.Boolean
		if err := db.Where("user_id = ? AND name = ?", user.ID, name).First(&existing).Error; err == nil {
			c.JSON(409, models.MessageResponse{Message: "Boolean already exists for this user"})
			return
		} else if err != gorm.ErrRecordNotFound {
			c.JSON(500, models.MessageResponse{Message: "Database error"})
			return
		}

		value := strings.ToLower(c.Query("initialValue"))
		var valueAsBool bool
		if value == "true" {
			valueAsBool = true
		} else {
			valueAsBool = false
		}

		boolean := models.Boolean{
			UserID: user.ID,
			Name:   name,
			Value:  valueAsBool,
		}

		if err := db.Create(&boolean).Error; err != nil {
			c.JSON(500, models.MessageResponse{Message: "Failed to create boolean"})
			return
		}

		c.JSON(201, models.MessageResponse{Message: fmt.Sprintf("Boolean with name %v created with initial state %v", name, valueAsBool)})
	}
}
