package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/im45145v/BinBag-Assignment/config"
	"github.com/im45145v/BinBag-Assignment/routes"
)

func main() {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Initialize Gin router
	router := gin.Default()

	// Configure trusted proxies
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatal("Failed to set trusted proxies:", err)
	}

	// Initialize MongoDB connection with timeout
	log.Printf("Connecting to MongoDB at %s", config.MongoURI)

	// Set MongoDB client options with longer timeout and direct connection
	clientOptions := options.Client().
		ApplyURI(config.MongoURI).
		SetConnectTimeout(15 * time.Second).
		SetServerSelectionTimeout(15 * time.Second)

	// Try to connect with retries
	var client *mongo.Client
	var err error

	for attempts := 1; attempts <= 3; attempts++ {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		log.Printf("Connection attempt %d...", attempts)
		client, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			// Successfully connected
			break
		}
		log.Printf("Connection attempt %d failed: %v", attempts, err)

		if attempts < 3 {
			time.Sleep(2 * time.Second) // Wait before retrying
		}
	}

	if err != nil {
		log.Fatal("Failed to connect to MongoDB after all attempts:", err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal("Failed to disconnect MongoDB:", err)
		}
	}()

	// Ping the database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	log.Println("Successfully connected to MongoDB Atlas")

	// Set up routes
	routes.SetupRoutes(router, client)

	// Start the server
	log.Println("Server running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
