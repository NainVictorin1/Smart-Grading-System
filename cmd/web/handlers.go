package main

import (
	"net/http"

	"github.com/NainVictorin/smart-grade-system/internal/data"
	"github.com/NainVictorin1/smart-grade-system/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Welcome"
	data.HeaderText = "We are here to help"
	err := app.render(w, http.StatusOK, "home.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render home page", "template", "home.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) submitGrade(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	subject := r.PostForm.Get("subject")
	gradeStr := r.PostForm.Get("grade")

	grade := &data.Grade{
		Fullname: name,
		Email:    email,
		Subject:  subject,
		Grade:    gradeStr,
	}

	// validate data
	v := validator.NewValidator()
	data.ValidateGrade(v, grade)

	if !v.ValidData() {
		td := NewTemplateData()
		td.Title = "Submit Grade"
		td.HeaderText = "Enter Student Grade"
		td.FormErrors = v.Errors
		td.FormData = map[string]string{
			"name":    name,
			"email":   email,
			"subject": subject,
			"grade":   gradeStr,
		}

		err := app.render(w, http.StatusUnprocessableEntity, "grade.tmpl", td)
		if err != nil {
			app.logger.Error("failed to render grade page", "template", "grade.tmpl", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.grades.Insert(grade)
	if err != nil {
		app.logger.Error("failed to insert grade", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/grades", http.StatusSeeOther)
}
