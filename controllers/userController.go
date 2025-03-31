package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/im45145v/BinBag-Assignment/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateProfile(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(string)
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var updateData map[string]interface{}
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
			return
		}

		collection := client.Database(config.DatabaseName).Collection(config.UsersCollection)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		updateFields := bson.M{}
		if name, ok := updateData["name"].(string); ok {
			updateFields["name"] = name
		}
		if address, ok := updateData["address"].(string); ok {
			updateFields["address"] = address
		}
		if bio, ok := updateData["bio"].(string); ok {
			updateFields["bio"] = bio
		}
		if profilePicture, ok := updateData["profile_picture"].(string); ok {
			updateFields["profile_picture"] = profilePicture
		}

		update := bson.M{"$set": updateFields}
		_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
	}
}
