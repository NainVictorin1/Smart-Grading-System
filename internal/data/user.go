package data

// Package data provides the data models and database interactions for the application.
import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// A struct to hold a user
type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Activated    bool
}

func (u *User) SetPassword(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	return nil
}

// MatchesPassword compares the stored hash with a plaintext password to verify authentication
func (u *User) MatchesPassword(plaintext string) error {
	// bcrypt.CompareHashAndPassword returns nil if the password matches the hash
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plaintext))
}
