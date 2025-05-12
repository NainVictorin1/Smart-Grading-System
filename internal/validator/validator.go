package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// emailRegex is a compiled regular expression for validating email format.
var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type Validator struct {
	Errors map[string]string
}

func Matches(value string, rx *regexp.Regexp) bool {

	return rx.MatchString(value)

}

func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) ValidData() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(field string, message string) {
	_, exists := v.Errors[field]
	if !exists {
		v.Errors[field] = message
	}
}

// Check adds an error message if the condition is false.
func (v *Validator) Check(ok bool, field, message string) {
	if !ok {
		v.Errors[field] = message
	}
}

// NotBlank checks if a string is not empty or contains only whitespace.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinLength(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func MaxLength(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// ValidateEmail checks if a string is a valid email address.
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// New creates and returns a new Validator instance.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if there are no errors.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}
