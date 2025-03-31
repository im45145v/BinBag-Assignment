package controllers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/im45145v/BinBag-Assignment/config"
	"github.com/im45145v/BinBag-Assignment/models"
	"github.com/im45145v/BinBag-Assignment/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Register endpoint hit")

		bodyBytes, err := c.GetRawData()
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}

		log.Printf("Raw request body: %s", string(bodyBytes))

		c.Request.Body = NewBodyReader(bodyBytes)

		var user models.User
		var requestData map[string]interface{}

		if err := c.ShouldBindJSON(&requestData); err != nil {
			log.Printf("Error binding JSON to map: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
			return
		}

		log.Printf("Parsed request data: %+v", requestData)

		passwordVal, exists := requestData["password"]
		if !exists {
			log.Printf("Password field missing from request")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password field is required"})
			return
		}

		passwordStr, ok := passwordVal.(string)
		if !ok {
			log.Printf("Password is not a string type: %T", passwordVal)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be a string"})
			return
		}

		passwordStr = strings.TrimSpace(passwordStr)
		if passwordStr == "" {
			log.Printf("Password is empty or contains only whitespace")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password cannot be empty"})
			return
		}

		c.Request.Body = NewBodyReader(bodyBytes)
		if err := c.ShouldBindJSON(&user); err != nil {
			log.Printf("Error binding JSON to user struct: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data format", "details": err.Error()})
			return
		}

		user.Password = passwordStr
		log.Printf("Password successfully parsed, length: %d", len(user.Password))

		user.ID = primitive.NewObjectID()

		if err := user.HashPassword(user.Password); err != nil {
			log.Printf("Error hashing password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		if user.Password == "" {
			log.Printf("Error: Password is empty after hashing")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Password storage failed"})
			return
		}

		 // Handle bio and profile_picture fields
		if bio, ok := requestData["bio"].(string); ok {
			user.Bio = bio
		}
		if profilePicture, ok := requestData["profile_picture"].(string); ok {
			user.ProfilePicture = profilePicture
		}

		collection := client.Database(config.DatabaseName).Collection(config.UsersCollection)

		var existingUser models.User
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
			return
		} else if err != mongo.ErrNoDocuments {
			log.Printf("Database error during email check: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = collection.InsertOne(ctx, user)
		if err != nil {
			log.Printf("Error inserting user into database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		token, err := utils.GenerateToken(user.ID.Hex(), user.Email)
		if err != nil {
			log.Printf("Error generating token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User registered successfully",
			"token":   token,
		})
	}
}

type bodyReader struct {
	*strings.Reader
}

func (b bodyReader) Close() error {
	return nil
}

func NewBodyReader(body []byte) bodyReader {
	return bodyReader{strings.NewReader(string(body))}
}

func Login(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Login endpoint hit")

		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&credentials); err != nil {
			log.Printf("Error binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
			return
		}

		if credentials.Email == "" || credentials.Password == "" {
			log.Printf("Missing required fields")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
			return
		}

		log.Printf("Credentials parsed successfully - Email: %s, Password length: %d",
			credentials.Email, len(credentials.Password))

		collection := client.Database(config.DatabaseName).Collection(config.UsersCollection)
		var user models.User

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := collection.FindOne(ctx, bson.M{"email": credentials.Email}).Decode(&user)
		if err != nil {
			log.Printf("User not found: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		log.Printf("User found, validating password")

		log.Printf("DEBUG - Stored hash in DB: %s", user.Password)

		if !user.CheckPassword(credentials.Password) {
			log.Printf("Password check failed for user: %s", credentials.Email)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		log.Printf("Password verified successfully for user: %s", credentials.Email)

		token, err := utils.GenerateToken(user.ID.Hex(), user.Email)
		if err != nil {
			log.Printf("Error generating token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func GetProfile(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("GetProfile endpoint hit")

		userID := c.MustGet("userID").(string)
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		collection := client.Database(config.DatabaseName).Collection(config.UsersCollection)
		var user models.User

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
		if err != nil {
			log.Printf("User not found: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":              user.ID.Hex(),
			"name":            user.Name,
			"email":           user.Email,
			"address":         user.Address,
			"bio":             user.Bio,
			"profile_picture": user.ProfilePicture,
		})
	}
}
