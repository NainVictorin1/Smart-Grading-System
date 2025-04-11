package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NainVictorin1/smart-grade-system/internal/validator"
)

type Grade struct {
	ID        int64
	CreatedAt time.Time
	Fullname  string
	Subject   string
	Grade     float64
	Email     string
}

type GradeModel struct {
	DB *sql.DB
}

// Insert a new grade into the database
func (m *GradeModel) Insert(grade *Grade) error {
	query := `
		INSERT INTO grade (fullname, subject, grade, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query,
		grade.Fullname,
		grade.Subject,
		grade.Grade,
		grade.Email,
	).Scan(&grade.ID, &grade.CreatedAt)
}

// ValidateGrade checks the grade fields for validity
func ValidateGrade(v *validator.Validator, grade *Grade) {
	v.Check(validator.NotBlank(grade.Fullname), "fullname", "must be provided")
	v.Check(validator.MaxLength(grade.Fullname, 50), "fullname", "must not be more than 50 characters")
	v.Check(validator.NotBlank(grade.Subject), "subject", "must be provided")
	v.Check(validator.MaxLength(grade.Subject, 50), "subject", "must not be more than 50 characters")
	v.Check(validator.NotBlank(grade.Email), "email", "must be provided")
	v.Check(validator.IsValidEmail(grade.Email), "email", "invalid email address")
	v.Check(grade.Grade >= 0 && grade.Grade <= 100, "grade", "must be between 0 and 100")
}

// GetAllGrades retrieves all grades from the database
func (m *GradeModel) GetAllGrades() ([]Grade, error) {
	query := `SELECT id, fullname, subject, grade, email, created_at FROM grade`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grades []Grade
	for rows.Next() {
		var grade Grade
		if err := rows.Scan(&grade.ID, &grade.Fullname, &grade.Subject, &grade.Grade, &grade.Email, &grade.CreatedAt); err != nil {
			return nil, err
		}
		grades = append(grades, grade)
	}
	return grades, nil
}

// CreateGrade adds a new grade to the database
func (m *GradeModel) CreateGrade(student string, grade float64, subject string) (int, error) {
	stmt := `INSERT INTO grade (fullname, subject, grade, email) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := m.DB.QueryRow(stmt, student, grade, subject).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// DeleteGrade deletes a grade by its ID
func (m *GradeModel) DeleteGrade(id int) error {
	stmt := `DELETE FROM grade WHERE id = $1`
	_, err := m.DB.Exec(stmt, id)
	return err
}

// UpdateGrade updates an existing grade in the database

func (m *GradeModel) UpdateGrade(g *Grade) error {
	query := `
        UPDATE grade
        SET fullname = $1, email = $2, subject = $3, grade = $4
        WHERE id = $5
    `
	fmt.Printf("Executing query: %s\n", query)
	fmt.Printf("Parameters: Fullname=%s, Email=%s, Subject=%s, Grade=%.2f, ID=%d\n",
		g.Fullname, g.Email, g.Subject, g.Grade, g.ID)

	_, err := m.DB.Exec(query, g.Fullname, g.Email, g.Subject, g.Grade, g.ID)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
	}
	return err
}

func (m *GradeModel) GetGradeByID(id int) (*Grade, error) {
	query := `SELECT id, fullname, email, subject, grade FROM grade WHERE id = $1`

	row := m.DB.QueryRow(query, id)

	grade := &Grade{}

	err := row.Scan(&grade.ID, &grade.Fullname, &grade.Email, &grade.Subject, &grade.Grade)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("grade not found")
		}
		return nil, err
	}

	return grade, nil
}
