package config

import (
	"os"
)

// MongoDB configuration
var (
	// MongoURI is the connection string for MongoDB
	// Default to a standard local MongoDB URI if environment variable not set
	MongoURI = getEnv("MONGO_URI", "mongodb://localhost:27017")

	// DatabaseName is the name of the MongoDB database
	DatabaseName = getEnv("DB_NAME", "binbag_db")

	// UsersCollection is the name of the users collection in MongoDB
	UsersCollection = getEnv("USERS_COLLECTION", "users")
)

// JWT configuration
var (
	// JWTSecretKey is the secret key used to sign JWT tokens
	JWTSecretKey = getEnv("JWT_SECRET_KEY", "your_default_jwt_secret_key")

	// JWTExpirationHours determines how long a JWT token is valid
	JWTExpirationHours = 24
)

// Helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
