package config

import (
	"os"
)

var (
	MongoURI        = getEnv("MONGO_URI", "mongodb://localhost:27017")
	DatabaseName    = getEnv("DB_NAME", "binbag_db")
	UsersCollection = getEnv("USERS_COLLECTION", "users")
	Bio             = getEnv("BIO", "default_bio")
	ProfilePicture  = getEnv("PROFILE_PICTURE", "default_profile_picture_url")
)

var (
	JWTSecretKey       = getEnv("JWT_SECRET_KEY", "your_default_jwt_secret_key")
	JWTExpirationHours = 24
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
