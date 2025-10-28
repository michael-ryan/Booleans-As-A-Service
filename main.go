package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"mike-pr.com/booleans_as_a_service/booleans"
	_ "mike-pr.com/booleans_as_a_service/docs"
	"mike-pr.com/booleans_as_a_service/models"
	"mike-pr.com/booleans_as_a_service/user"
	"mike-pr.com/booleans_as_a_service/users"
)

var db *gorm.DB

// @title Booleans as a Service API
// @description This is a simple API to manage booleans. Supports user-level private booleans.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-KEY
func main() {
	db = initDB()
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// users
	r.POST("/users", users.Create(db))

	// user
	r.DELETE("/users/:username", user.Delete(db))

	// booleans
	r.GET("/users/:username/booleans", booleans.Get(db))
	r.POST("/users/:username/booleans", booleans.Create(db))

	// // boolean
	// r.GET("/users/:username/booleans/:boolean", getBoolean)
	// r.POST("/users/:username/booleans/:boolean", createBoolean)
	// r.DELETE("/users/:username/booleans/:boolean", deleteBoolean)

	if err := r.Run(":8080"); err != nil {
		panic(fmt.Errorf("Could not start server: %w", err))
	}
}

func initDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	err = db.AutoMigrate(&models.User{}, &models.Boolean{})
	if err != nil {
		panic(fmt.Errorf("failed to initialise DB: %w", err))
	}
	return db
}
