package data

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/NainVictorin1/smart-grade-system/internal/validator"
)

type UserModel struct {
	DB *sql.DB
}

var (
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrRecordNotFound     = errors.New("record not found")
	ErrInvalidCredentials = errors.New(" invalid credentials")
)

func ValidateUser(v *validator.Validator, user *User, r *http.Request) {
	password := r.PostForm.Get("password")
	v.Check(validator.NotBlank(user.Name), "name", "Name must be provided")
	v.Check(validator.MaxLength(user.Name, 50), "name", "Name must not be more than 50 characters")
	v.Check(validator.NotBlank(user.Email), "email", "Email must be provided")
	v.Check(validator.IsValidEmail(user.Email), "email", "Invalid email address")
	v.Check(validator.NotBlank(password), "password", "Password must be provided")
	v.Check(len(password) >= 8, "password", "Password must be at least 8 characters")
}

func (m *UserModel) Insert(user *User) error {
	query := `
        INSERT INTO users (name, email, password_hash, activated)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at`

	err := m.DB.QueryRow(query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Activated,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_email_key"`) {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

// GetByEmail fetches a user from the database by email address
func (m *UserModel) GetByEmail(email string) (*User, error) {
	// Query the database and scan the row into a User struct
	query := `
		SELECT id, email, password_hash, created_at,
		FROM users
		WHERE email = $1`

	var user User
	err := m.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	// If no rows returned, wrap and return ErrRecordNotFound
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}
func (m *UserModel) GetByID(id int) (*User, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE id = $1`
	// Execute the query and populate the user struct
	var user User
	err := m.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	// Handle case where user ID does not exist in the database
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}
