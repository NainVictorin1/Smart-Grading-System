package main

import (
	"net/url"

	"github.com/NainVictorin1/smart-grade-system/internal/data"
)

type TemplateData struct {
	Title           string
	HeaderText      string
	FormErrors      map[string]string
	ErrorsFromForm  map[string]string
	FormData        url.Values   // A map to hold form data
	Grades          []data.Grade // A slice to hold multiple grades
	IsSubmitted     bool         // Indicates if the form has been submitted
	ID              int
	SuccessMessage  string
	Flash           string
	CSRFToken       string // Added CSRFToken field
	IsAuthenticated bool
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:      "Default Title",
		HeaderText: "Default HeaderText",
		FormErrors: make(map[string]string),
		FormData:   url.Values{},   // Initialize as an empty url.Values
		Grades:     []data.Grade{}, // Initialize as an empty slice
	}
}
