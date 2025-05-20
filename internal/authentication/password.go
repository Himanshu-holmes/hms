package authentication

import (
	"log" // For logging errors

	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultCost = 12
)

func HashPassword(password string) (string, error) {
	// bcrypt.GenerateFromPassword returns []byte, so we convert to string.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		log.Printf("Error generating password hash: %v", err) 
		return "", err                                       
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	
	if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
		// Log unexpected errors during password comparison, as they might indicate issues.
		log.Printf("Error comparing password hash (potential hash issue, not just mismatch): %v", err)
	}
	return err == nil
}
