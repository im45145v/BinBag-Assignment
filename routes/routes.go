package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/im45145v/BinBag-Assignment/config"
	"github.com/im45145v/BinBag-Assignment/controllers"
	"github.com/im45145v/BinBag-Assignment/middlewares"
	"github.com/im45145v/BinBag-Assignment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(router *gin.Engine, client *mongo.Client) {
	// Public routes
	router.POST("/register", controllers.Register(client))
	log.Println("Registered route: POST /register")

	router.POST("/login", controllers.Login(client))
	log.Println("Registered route: POST /login")

	// Protected routes
	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/profile", controllers.GetProfile(client))
		log.Println("Registered route: GET /profile (protected)")
	}

	// Add a test route to verify the server is running
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	log.Println("Registered route: GET /ping")

	// Add a temporary route to reset a password for troubleshooting
	router.POST("/reset-password", func(c *gin.Context) {
		var data struct {
			Email       string `json:"email" binding:"required"`
			NewPassword string `json:"new_password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		collection := client.Database(config.DatabaseName).Collection(config.UsersCollection)
		var user models.User

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := collection.FindOne(ctx, bson.M{"email": data.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Hash the new password
		if err := user.HashPassword(data.NewPassword); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Update the password in the database
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = collection.UpdateOne(
			ctx,
			bson.M{"email": data.Email},
			bson.M{"$set": bson.M{"password": user.Password}},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully", "hash": user.Password})
	})
	log.Println("Registered temporary route: POST /reset-password")

	// Add a test endpoint for password verification
	router.POST("/test-password", func(c *gin.Context) {
		var data struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		collection := client.Database(config.DatabaseName).Collection(config.UsersCollection)
		var user models.User

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := collection.FindOne(ctx, bson.M{"email": data.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Try verification and return detailed results
		log.Printf("TEST - Raw password provided: '%s'", data.Password)
		result := user.CheckPassword(data.Password)

		c.JSON(http.StatusOK, gin.H{
			"success":         result,
			"stored_hash":     user.Password,
			"password_length": len(data.Password),
			"hash_length":     len(user.Password),
		})
	})
	log.Println("Registered test route: POST /test-password")
}
