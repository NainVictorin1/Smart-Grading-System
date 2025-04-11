package main

import (
	"github.com/NainVictorin1/smart-grade-system/internal/data"
)

type TemplateData struct {
	Title          string
	HeaderText     string
	FormErrors     map[string]string
	FormData       map[string]string
	Grades         []data.Grade // A slice to hold multiple grades
	IsSubmitted    bool         // Indicates if the form has been submitted
	ID             int
	SuccessMessage string
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:      "Default Title",
		HeaderText: "Default HeaderText",
		FormErrors: make(map[string]string),
		FormData:   make(map[string]string),
		Grades:     []data.Grade{}, // Initialize as an empty slice
	}
}
