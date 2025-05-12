package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/NainVictorin1/smart-grade-system/internal/data"

	"github.com/justinas/nosurf"

	"github.com/NainVictorin1/smart-grade-system/internal/validator"
)

// Handler for the home page
func (app *application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch all grades from the database
	grades, err := app.grades.GetAllGrades()
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Prepare the template data
	td := NewTemplateData()
	td.Grades = grades
	td.CSRFToken = nosurf.Token(r)
	td.IsAuthenticated = app.isAuthenticated(r) // Check if the user is logged in

	// Render the home page using the "home.tmpl" template
	app.Render(w, http.StatusOK, "home.tmpl", td)
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
	data.CSRFToken = nosurf.Token(r)

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
		// Parse form values
		fullname := r.PostFormValue("fullname")
		email := r.PostFormValue("email")
		subject := r.PostFormValue("subject")
		gradeStr := r.PostFormValue("grade")

		// Log the submitted data
		app.logger.Info("Submitted data", "fullname", fullname, "email", email, "subject", subject, "grade", gradeStr)

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
		data.ValidateGrade(v, grade)

		// If validation fails, re-render the form with errors
		if !v.Valid() {
			td := NewTemplateData()
			td.FormData = url.Values{
				"fullname": []string{fullname},
				"email":    []string{email},
				"subject":  []string{subject},
				"grade":    []string{gradeStr},
			}
			td.FormErrors = v.Errors

			err := app.Render(w, http.StatusOK, "create_grade.tmpl", td)
			if err != nil {
				app.logger.Error("Failed to render create_grade.tmpl", "error", err)
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
		// app.sessions.Put(r, "flash", "Grade successfully add.")
		// Redirect to the grades list
		http.Redirect(w, r, "/grades", http.StatusSeeOther)
		return
	}

	// Render the create grade form
	td := NewTemplateData()
	err := app.Render(w, http.StatusOK, "create_grade.tmpl", td)
	if err != nil {
		app.logger.Error("Failed to render create_grade.tmpl", "error", err)
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
			td.FormData = url.Values{
				"fullname": []string{fullname},
				"email":    []string{email},
				"subject":  []string{subject},
				"grade":    []string{gradeStr},
			}
			td.FormErrors = map[string]string{
				"grade": "Grade must be a number between 0 and 100",
			}
			td.ID = id

			err = app.Render(w, http.StatusOK, "edit_grade.tmpl", td)
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
	td.FormData = url.Values{
		"fullname": []string{grade.Fullname},
		"email":    []string{grade.Email},
		"subject":  []string{grade.Subject},
		"grade":    []string{fmt.Sprintf("%.2f", grade.Grade)},
	}
	td.ID = int(grade.ID)

	err = app.Render(w, http.StatusOK, "edit_grade.tmpl", td)
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

	err = app.Render(w, http.StatusOK, "grade.tmpl", td)
	if err != nil {
		app.logger.Error("Failed to render grade.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
func (app *application) LoginForm(w http.ResponseWriter, r *http.Request) {
	td := NewTemplateData()
	td.CSRFToken = nosurf.Token(r)
	td.IsAuthenticated = app.isAuthenticated(r) // Set the IsAuthenticated field

	err := app.Render(w, http.StatusOK, "login.tmpl", td)
	if err != nil {
		app.logger.Error("failed to render template", "template", "login.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {
	session, err := app.SessionStore.Get(r, SessionName)
	if err != nil {
		return false // If there's an error retrieving the session, assume the user is not authenticated
	}

	// Check if the session contains the user ID
	_, ok := session.Values[SessionUserKey]
	return ok
}

// LoginHandler authenticates the user and starts a session.
func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Extract email and password from the form
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	// Validate email and password
	v := validator.New()
	v.Check(validator.NotBlank(email), "email", "Email is required")
	v.Check(validator.IsValidEmail(email), "email", "Invalid email format")
	v.Check(validator.NotBlank(password), "password", "Password is required")

	if !v.Valid() {
		// Render the form again with validation errors
		td := NewTemplateData()
		td.FormErrors = v.Errors
		td.FormData = url.Values{
			"email": []string{email}, // Preserve the entered email
		}
		app.Render(w, http.StatusOK, "login.tmpl", td)
		return
	}

	// Lookup user by email
	user, err := app.users.GetByEmail(email)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			// Show invalid credentials if user not found
			td := NewTemplateData()
			td.FormErrors = map[string]string{
				"Error": "Invalid email or password",
			}
			td.FormData = url.Values{
				"email": []string{email}, // Preserve the entered email
			}
			app.Render(w, http.StatusOK, "login.tmpl", td)
			return
		}
		app.ServerError(w, err)
		return
	}

	// Compare passwords
	err = user.MatchesPassword(password)
	if err != nil {
		td := NewTemplateData()
		td.FormErrors = map[string]string{
			"Error": "Invalid email or password",
		}
		td.FormData = url.Values{
			"email": []string{email}, // Preserve the entered email
		}
		app.Render(w, http.StatusOK, "login.tmpl", td)
		return
	}

	// Create a session and store user ID
	session, err := app.SessionStore.Get(r, SessionName)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	session.Values[SessionUserKey] = user.ID
	if err := session.Save(r, w); err != nil {
		app.ServerError(w, err)
		return
	}

	// Redirect to the home page after successful login
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// SignupForm displays the signup form.
func (app *application) SignupForm(w http.ResponseWriter, r *http.Request) {
	td := NewTemplateData()
	td.CSRFToken = nosurf.Token(r)              // Add CSRF token for protection
	td.IsAuthenticated = app.isAuthenticated(r) // Set the IsAuthenticated field
	td.ErrorsFromForm = make(map[string]string) // Initialize ErrorsFromForm

	app.Render(w, http.StatusOK, "signup.tmpl", td)
}

// SignupHandler processes user registration.
func (app *application) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Create a new user instance
	user := &data.User{
		Name:      r.PostForm.Get("name"),
		Email:     r.PostForm.Get("email"),
		Activated: true,
	}

	// Validate the user using the centralized validation function
	v := validator.New()
	data.ValidateUser(v, user, r)

	if !v.Valid() {
		// Render the form again with validation errors
		td := NewTemplateData()
		td.FormErrors = v.Errors
		td.FormData = url.Values{
			"name":  []string{user.Name},
			"email": []string{user.Email},
		}
		app.Render(w, http.StatusOK, "signup.tmpl", td)
		return
	}

	// Hash the password
	password := r.PostForm.Get("password")
	err = user.SetPassword(password)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Attempt to insert the user into the database
	err = app.users.Insert(user)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateEmail) {
			// Handle duplicate email error
			v.Errors["email"] = "Email already in use"
			td := NewTemplateData()
			td.FormErrors = v.Errors
			td.FormData = url.Values{
				"name":  []string{user.Name},
				"email": []string{user.Email},
			}
			app.Render(w, http.StatusOK, "signup.tmpl", td)
			return
		}
		app.ServerError(w, err)
		return
	}

	// Redirect to the login page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := app.SessionStore.Get(r, SessionName)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	delete(session.Values, SessionUserKey)
	if err := session.Save(r, w); err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
