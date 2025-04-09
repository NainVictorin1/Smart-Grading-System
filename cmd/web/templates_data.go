package main

type TemplateData struct {
	Title      string
	HeaderText string
	FormErrors map[string]string
	FormData   map[string]string
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:      "Default Title",
		HeaderText: "Default HeaderText",
		FormErrors: map[string]string{},
		FormData:   map[string]string{},
	}
}
