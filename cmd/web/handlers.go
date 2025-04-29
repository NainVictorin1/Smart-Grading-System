package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/NainVictorin1/smart-grade-system/internal/data"

	"github.com/NainVictorin1/smart-grade-system/internal/validator"
)

// ValidateGrade checks the grade fields for validity
type Grade struct {
	Fullname string
	Email    string
	Subject  string
	Grade    float64
}

func ValidateGrade(v *validator.Validator, grade *Grade) {
	v.Check(validator.NotBlank(grade.Fullname), "fullname", "must be provided")
	v.Check(validator.MaxLength(grade.Fullname, 50), "fullname", "must not be more than 50 characters")
	v.Check(validator.NotBlank(grade.Subject), "subject", "must be provided")
	v.Check(validator.MaxLength(grade.Subject, 50), "subject", "must not be more than 50 characters")
	v.Check(validator.NotBlank(grade.Email), "email", "must be provided")
	v.Check(validator.IsValidEmail(grade.Email), "email", "invalid email address")
	v.Check(grade.Grade >= 0 && grade.Grade <= 100, "grade", "must be between 0 and 100")
}

// Handler for the home page
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := app.templateCache["home.tmpl"] // Note: no "./ui/html/" here
	if !ok {
		app.logger.Error("Unable to load template", "template", "home.tmpl")
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}
	err := tmpl.Execute(w, nil)
	if err != nil {
		app.logger.Error("Unable to render template", "template", "home.tmpl", "error", err)
		http.Error(w, fmt.Sprintf("Unable to render template: %v", err), http.StatusInternalServerError)
	}
}

// Handler to view grades
func (app *application) viewGrade(w http.ResponseWriter, r *http.Request) {
	grades, err := app.grades.GetAllGrades()
	if err != nil {
		app.logger.Error("Failed to fetch grades from the database", "error", err)
		http.Error(w, "Internal Server Error: Unable to fetch grades", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData()
	data.Title = "View Grades"
	data.HeaderText = "Student Grades"
	data.Grades = grades

	tmpl, ok := app.templateCache["grade.tmpl"]
	if !ok {
		app.logger.Error("Template not found in cache", "template", "grade.tmpl")
		http.Error(w, "Internal Server Error: Template not found", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		app.logger.Error("Failed to render the grade template", "template", "grade.tmpl", "error", err)
		http.Error(w, "Internal Server Error: Unable to render template", http.StatusInternalServerError)
		return
	}
}

// Handler to create a grade
func (app *application) createGrade(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			app.logger.Error("Failed to parse grade form", "error", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Extract form values
		fullname := r.PostForm.Get("fullname")
		email := r.PostForm.Get("email")
		subject := r.PostForm.Get("subject")
		gradeStr := r.PostForm.Get("grade")

		// Log the submitted data
		app.logger.Info("Form values", "fullname", fullname, "email", email, "subject", subject, "grade", gradeStr)

		// Convert grade string to float64
		gradeValue, err := strconv.ParseFloat(gradeStr, 64)
		if err != nil {
			gradeValue = -1 // Invalid grade value
		}

		// Create a Grade object
		grade := &data.Grade{
			Fullname: fullname,
			Email:    email,
			Subject:  subject,
			Grade:    gradeValue,
		}

		// Validate the Grade object
		v := validator.New()
		ValidateGrade(v, &Grade{
			Fullname: grade.Fullname,
			Email:    grade.Email,
			Subject:  grade.Subject,
			Grade:    grade.Grade,
		})

		// If validation fails, re-render the form with errors
		if !v.Valid() {
			data := NewTemplateData()
			data.FormErrors = v.Errors
			data.FormData = map[string]string{
				"fullname": fullname,
				"email":    email,
				"subject":  subject,
				"grade":    gradeStr,
			}
			tmpl, ok := app.templateCache["create_grade.tmpl"]
			if !ok {
				app.logger.Error("Template not found in cache", "template", "create_grade.tmpl")
				http.Error(w, "Internal Server Error: Template not found", http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, data)
			if err != nil {
				app.logger.Error("Failed to render the create_grade template", "template", "create_grade.tmpl", "error", err)
				http.Error(w, "Internal Server Error: Unable to render template", http.StatusInternalServerError)
			}
			return
		}

		// Insert the grade into the database
		err = app.grades.Insert(grade)
		if err != nil {
			app.logger.Error("Failed to insert grade", "error", err)
			http.Error(w, "Internal Server Error: Unable to save grade", http.StatusInternalServerError)
			return
		}

		// Redirect to the grades list
		http.Redirect(w, r, "/grade", http.StatusSeeOther)
		return
	}

	// Render the create grade form
	data := NewTemplateData()
	tmpl, ok := app.templateCache["create_grade.tmpl"]
	if !ok {
		app.logger.Error("Template not found in cache", "template", "create_grade.tmpl")
		http.Error(w, "Internal Server Error: Template not found", http.StatusInternalServerError)
		return
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		app.logger.Error("Failed to render the create_grade template", "template", "create_grade.tmpl", "error", err)
		http.Error(w, "Internal Server Error: Unable to render template", http.StatusInternalServerError)
	}
}
func (app *application) editGrade(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		// Parse form values
		fullname := r.PostFormValue("fullname")
		email := r.PostFormValue("email")
		subject := r.PostFormValue("subject")
		gradeStr := r.PostFormValue("grade")

		// Log the submitted data
		app.logger.Info("Submitted data", "fullname", fullname, "email", email, "subject", subject, "grade", gradeStr)

		// Convert grade
		gradeValue, err := strconv.ParseFloat(gradeStr, 64)
		if err != nil {
			http.Error(w, "Invalid grade value", http.StatusBadRequest)
			return
		}

		// Validate grade range
		if gradeValue < 0 || gradeValue > 100 {
			td := NewTemplateData()
			td.FormData = map[string]string{
				"fullname": fullname,
				"email":    email,
				"subject":  subject,
				"grade":    gradeStr,
			}
			td.FormErrors = map[string]string{
				"grade": "Grade must be a number between 0 and 100",
			}
			td.ID = id

			err = app.render(w, http.StatusOK, "edit_grade.tmpl", td)
			if err != nil {
				app.logger.Error("Failed to render edit_grade.tmpl", "error", err)
				http.Error(w, "Internal Server Error: Unable to render template", http.StatusInternalServerError)
			}
			return
		}

		// Construct grade object
		updatedGrade := &data.Grade{
			ID:       int64(id),
			Fullname: fullname,
			Email:    email,
			Subject:  subject,
			Grade:    gradeValue,
		}

		// Update the grade in the database
		err = app.grades.UpdateGrade(updatedGrade)
		if err != nil {
			app.logger.Error("Failed to update grade", "error", err)
			http.Error(w, "Internal Server Error: Unable to update grade", http.StatusInternalServerError)
			return
		}

		// Redirect to the grades list
		http.Redirect(w, r, "/grade", http.StatusSeeOther)
		return
	}

	// GET: Fetch the grade from the database
	grade, err := app.grades.GetGradeByID(id)
	if err != nil {
		app.logger.Error("Grade not found", "id", id, "error", err)
		http.Error(w, "Grade not found", http.StatusNotFound)
		return
	}

	// Render the edit form
	td := NewTemplateData()
	td.FormData = map[string]string{
		"fullname": grade.Fullname,
		"email":    grade.Email,
		"subject":  grade.Subject,
		"grade":    fmt.Sprintf("%.2f", grade.Grade),
	}
	td.ID = int(grade.ID)

	err = app.render(w, http.StatusOK, "edit_grade.tmpl", td)
	if err != nil {
		app.logger.Error("Failed to render edit_grade.tmpl", "error", err)
		http.Error(w, "Internal Server Error: Unable to render template", http.StatusInternalServerError)
	}
}

// Handler to delete a grade
func (app *application) deleteGrade(w http.ResponseWriter, r *http.Request) {
	// Extract ID from the query string (e.g., /delete?id=123)
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}

	// Convert the ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Delete the grade using the extracted ID
	err = app.grades.DeleteGrade(id)
	if err != nil {
		app.logger.Error("Failed to delete grade", "id", id, "error", err)
		http.Error(w, "Unable to delete grade", http.StatusInternalServerError)
		return
	}

	// Fetch the updated list of grades
	grades, err := app.grades.GetAllGrades()
	if err != nil {
		app.logger.Error("Failed to fetch grades after deletion", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the updated grades list with a success message
	td := NewTemplateData()
	td.Title = "Grades"
	td.HeaderText = "Grades List"
	td.Grades = grades
	td.SuccessMessage = "Grade successfully deleted."

	err = app.render(w, http.StatusOK, "grade.tmpl", td)
	if err != nil {
		app.logger.Error("Failed to render grade.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
