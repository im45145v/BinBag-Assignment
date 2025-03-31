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
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatal("Failed to set trusted proxies:", err)
	}

	log.Printf("Connecting to MongoDB at %s", config.MongoURI)

	clientOptions := options.Client().
		ApplyURI(config.MongoURI).
		SetConnectTimeout(15 * time.Second).
		SetServerSelectionTimeout(15 * time.Second)

	var client *mongo.Client
	var err error

	for attempts := 1; attempts <= 3; attempts++ {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		log.Printf("Connection attempt %d...", attempts)
		client, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			break
		}
		log.Printf("Connection attempt %d failed: %v", attempts, err)

		if attempts < 3 {
			time.Sleep(2 * time.Second)
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	log.Println("Successfully connected to MongoDB Atlas")

	routes.SetupRoutes(router, client)

	log.Println("Server running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
