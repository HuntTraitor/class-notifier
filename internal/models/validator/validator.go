package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

// regex for emails recommended by W3C and Web Hypertext Application Technology Working Group
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	FieldErrors map[string]string
}

// checks if form is valid
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// adds a new field error if it does not already exist
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// checks for a specific field error
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// checks if field is not blank
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// checks to see if two fields are equal
func Equal(value1 string, value2 string) bool {
	return value1 == value2
}

// checks for min number of characters
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// checks for max number of characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// checks if field matches regular expression
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// checks if form has a specific value
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
