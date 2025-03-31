package models

import (
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	Email          string             `json:"email" bson:"email"`
	Password       string             `json:"-" bson:"password"` // Changed to exclude password from JSON responses
	Address        string             `json:"address,omitempty" bson:"address,omitempty"`
	Bio            string             `json:"bio,omitempty" bson:"bio,omitempty"`
	ProfilePicture string             `json:"profile_picture,omitempty" bson:"profile_picture,omitempty"`
}

func (u *User) HashPassword(password string) error {
	// Prevent empty passwords
	if password == "" {
		return errors.New("password cannot be empty")
	}

	log.Printf("Hashing password with bcrypt (length: %d)", len(password))

	// Generate bcrypt hash from password with cost factor of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error generating bcrypt hash: %v", err)
		return err
	}

	// Store the hashed password
	u.Password = string(hashedPassword)

	log.Printf("Password hashed and stored successfully (hash length: %d)", len(u.Password))
	return nil
}

func (u *User) CheckPassword(password string) bool {
	// Prevent matching with empty stored passwords
	if u.Password == "" {
		log.Printf("Stored password hash is empty, rejecting authentication")
		return false
	}

	log.Printf("Comparing password: input length: %d, stored hash length: %d", len(password), len(u.Password))

	// Compare provided password with stored hash
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		log.Printf("Password comparison failed: %v", err)
		return false
	}

	log.Printf("Password verification successful")
	return true
}
