package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "mike-pr.com/booleans_as_a_service/docs"
	"mike-pr.com/booleans_as_a_service/handlers/boolean"
	"mike-pr.com/booleans_as_a_service/handlers/booleans"
	"mike-pr.com/booleans_as_a_service/handlers/user"
	"mike-pr.com/booleans_as_a_service/handlers/users"
	"mike-pr.com/booleans_as_a_service/models"
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

	// boolean
	r.GET("/users/:username/booleans/:boolean", boolean.Get(db))
	r.DELETE("/users/:username/booleans/:boolean", boolean.Delete(db))

	if err := r.Run(fmt.Sprintf(":%v", os.Getenv("PORT"))); err != nil {
		panic(fmt.Errorf("Could not start server: %w", err))
	}
}

func initDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	retries := 5
	dbOk := false
	var db *gorm.DB
	var err error
	for i := 1; i < retries; i++ {
		fmt.Printf("Attempt %v of %v to connect to database...\n", i, retries)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			dbOk = true
			break
		} else {
			fmt.Printf("Database connection failed. Wait 5s and try again.\n")
			time.Sleep(time.Second * 5)
		}
	}
	if !dbOk {
		panic(fmt.Sprintf("failed to connect database after %v attempts: ", retries) + err.Error())
	}
	err = db.AutoMigrate(&models.User{}, &models.Boolean{})
	if err != nil {
		panic(fmt.Errorf("failed to initialise DB: %w", err))
	}
	return db
}
