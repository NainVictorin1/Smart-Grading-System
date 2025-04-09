package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/NainVictorin1/smart-grade-system/internal/validator"
)

type Grade struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Fullname  string    `json:"fullname"`
	Subject   string    `json:"subject"`
	Grade     string    `json:"grade"`
	Email     string    `json:"email"`
}

// ValidateGrade validates the input grade form
func ValidateGrade(v *validator.Validator, grade *Grade) {
	v.Check(validator.NotBlank(grade.Fullname), "fullname", "must be provided")
	v.Check(validator.MaxLength(grade.Fullname, 50), "fullname", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(grade.Subject), "subject", "must be provided")
	v.Check(validator.MaxLength(grade.Subject, 50), "subject", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(grade.Email), "email", "must be provided")
	v.Check(validator.IsValidEmail(grade.Email), "email", "invalid email address")
	v.Check(validator.MaxLength(grade.Email, 100), "email", "must not be more than 100 bytes long")
	validGrades := map[string]bool{
		"A": true, "A-": true,
		"B+": true, "B": true, "B-": true,
		"C+": true, "C": true, "C-": true,
		"D+": true, "D": true, "D-": true,
		"F": true,
	}
	v.Check(validGrades[grade.Grade], "grade", "must be a valid letter grade (A to F)")
}

type GradeModel struct {
	DB *sql.DB
}

// Insert adds a new grade record to the database
func (m *GradeModel) Insert(grade *Grade) error {
	query := `
		INSERT INTO garde (fullname, subject, grade, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(
		ctx,
		query,
		grade.Fullname,
		grade.Subject,
		grade.Grade,
		grade.Email,
	).Scan(&grade.ID, &grade.CreatedAt)
}
