package validator

import "regexp"

// This var stores the standard regex pattern for email ID according to https://html.spec.whatwg.org/#valid-e-mail-address
var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Validator is a struct used to hold a map of validation errors
type Validator struct {
	Errors map[string]string
}

// Helper method to create a new Validator instance with an empty errors map
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if the errors map doesn't contain any entries
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the map, as long as no error exists for the existing key
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message to the map if the validation check returns False
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// In returns true if value is within a list of strings, else false
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}

	return false
}

// Matches returns true if value matches a specific regex pattern, else false
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique is a generic function for all comparable types which returns true if all values in the slice are unique
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, values := range values {
		uniqueValues[values] = true
	}

	return len(values) == len(uniqueValues)
}
