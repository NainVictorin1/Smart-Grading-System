package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/NainVictorin1/smart-grade-system/internal/data"

	"github.com/NainVictorin1/smart-grade-system/internal/validator"
)

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

	app.logger.Info("Form values", "fullname", fullname, "email", email, "subject", subject, "grade", gradeStr)

	// Convert grade string to float64
	gradeValue, err := strconv.ParseFloat(gradeStr, 64)
	if err != nil {
		app.logger.Error("Invalid grade format", "grade", gradeStr, "error", err)
	}

	// Create a Grade object
	grade := &data.Grade{
		Fullname: fullname,
		Email:    email,
		Subject:  subject,
		Grade:    gradeValue,
	}

	// Validate the grade object
	v := validator.NewValidator()
	data.ValidateGrade(v, grade)

	// Check for validation errors
	if !v.ValidData() {
		td := NewTemplateData()
		td.Title = "Submit Grade"
		td.HeaderText = "Enter Grade Details"
		td.FormErrors = v.Errors
		td.FormData = map[string]string{
			"fullname": fullname,
			"email":    email,
			"subject":  subject,
			"grade":    gradeStr,
		}
		td.IsSubmitted = true // Indicate that the form has been submitted

		err := app.render(w, http.StatusUnprocessableEntity, "create_grade.tmpl", td)
		if err != nil {
			app.logger.Error("Failed to render create_grade page", "template", "create_grade.tmpl", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Insert the grade into the database
	err = app.grades.Insert(grade)
	if err != nil {
		app.logger.Error("Failed to insert grade", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect to the success page
	http.Redirect(w, r, "/grade/success", http.StatusSeeOther)
}

// Helper function to parse and validate grade value
// func parseGrade(grade string) float64 {
// 	gradeValue := 0.0
// 	_, err := fmt.Sscanf(grade, "%f", &gradeValue)
// 	if err != nil {
// 		return 0.0
// 	}
// 	return gradeValue
// }

// Handler to edit a grade
// Handler to edit a grade
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

		// Construct grade object
		updatedGrade := &data.Grade{
			ID:       int64(id),
			Fullname: fullname,
			Email:    email,
			Subject:  subject,
			Grade:    gradeValue,
		}

		// Log the update operation
		app.logger.Info("Updating grade in database", "id", updatedGrade.ID, "fullname", updatedGrade.Fullname, "email", updatedGrade.Email, "subject", updatedGrade.Subject, "grade", updatedGrade.Grade)

		// Update in DB
		err = app.grades.UpdateGrade(updatedGrade)
		if err != nil {
			app.logger.Error("Failed to update grade", "error", err)
			http.Error(w, "Failed to update grade", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/grades", http.StatusSeeOther)
		return
	}

	// GET: fetch the grade from DB
	grade, err := app.grades.GetGradeByID(id)
	if err != nil {
		app.logger.Error("Grade not found", "id", id, "error", err)
		http.Error(w, "Grade not found", http.StatusNotFound)
		return
	}

	td := NewTemplateData()
	td.Title = "Edit Grade"
	td.HeaderText = "Edit Grade Details"
	td.FormData = map[string]string{
		"fullname": grade.Fullname,
		"email":    grade.Email,
		"subject":  grade.Subject,
		"grade":    fmt.Sprintf("%.2f", grade.Grade),
	}

	err = app.render(w, http.StatusOK, "edit_grade.tmpl", td)
	if err != nil {
		app.logger.Error("Failed to render edit_grade.tmpl", "error", err)
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
