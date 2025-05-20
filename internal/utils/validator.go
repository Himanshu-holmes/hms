package util

import (
	"fmt"
	"reflect"
	"strings"
	"time" // Required for custom datetime validation message formatting

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register a custom function to use JSON field names in error messages.
	// This makes error messages more aligned with what the API consumer sends.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// If the JSON tag is "-", it means the field should be ignored by JSON,
		// so we don't want it in validation error messages either.
		if name == "-" {
			return ""
		}
		return name
	})

	// Example: Register a custom validation for a specific format if needed.
	// For instance, if you had a custom 'uuid_custom' tag:
	// validate.RegisterValidation("uuid_custom", func(fl validator.FieldLevel) bool {
	// 	_, err := uuid.Parse(fl.Field().String())
	// 	return err == nil
	// })
}

// ValidateStruct performs validation on a given struct `s`.
// It returns an error if validation fails, otherwise nil.
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// FormatValidationErrors converts validator.ValidationErrors into a map[string]string.
// The keys of the map are the field names (using JSON tags), and values are
// user-friendly error messages.
func FormatValidationErrors(err error) map[string]string {
	formattedErrors := make(map[string]string)

	// Check if the error is of type validator.ValidationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// If it's not validation errors, return a generic message or the error itself
		// For simplicity, we'll return an empty map here, but you might want to handle it.
		// formattedErrors["_error"] = "An unexpected validation error occurred."
		return formattedErrors
	}

	for _, fieldErr := range validationErrors {
		// Use fieldErr.Field() which, due to RegisterTagNameFunc, gives the JSON field name.
		fieldName := fieldErr.Field()
		formattedErrors[fieldName] = formatFieldError(fieldErr)
	}
	return formattedErrors
}

// formatFieldError creates a user-friendly message for a single validation error.
func formatFieldError(err validator.FieldError) string {
	// The 'Tag' is the validation rule that failed (e.g., "required", "min", "email").
	// The 'Param' is the parameter associated with the tag (e.g., "6" for "min=6").
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return "This field is required."
	case "min":
		switch err.Kind() {
		case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
			return fmt.Sprintf("This field must have at least %s items/characters.", param)
		default: // numbers
			return fmt.Sprintf("This field must be at least %s.", param)
		}
	case "max":
		switch err.Kind() {
		case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
			return fmt.Sprintf("This field must have at most %s items/characters.", param)
		default: // numbers
			return fmt.Sprintf("This field must be at most %s.", param)
		}
	case "email":
		return "Invalid email format."
	case "oneof":
		return fmt.Sprintf("This field must be one of: [%s].", strings.ReplaceAll(param, " ", ", "))
	case "datetime":
		// The param for datetime is the layout string (e.g., "2006-01-02")
		// We can try to format this layout into a more human-readable example.
		exampleFormat := param
		if parsedTime, err := time.Parse(param, time.Now().Format(param)); err == nil {
			// If the param itself is a valid layout like "2006-01-02",
			// formatting time.Now() with it will give an example in that format.
			exampleFormat = parsedTime.Format(param)
		}
		return fmt.Sprintf("Invalid date/time format. Expected format like: %s.", exampleFormat)
	case "e164": // Example for a custom tag or one you might expect for phone numbers
		return "Invalid phone number format. Expected E.164 format (e.g., +12125551234)."
	// Add more cases for other common validation tags as needed:
	// case "url":
	// 	return "Invalid URL format."
	// case "uuid": // if using the built-in uuid tag
	//  return "Invalid UUID format."
	// case "eqfield", "nefield":
	// 	return fmt.Sprintf("This field must be equal/not equal to the %s field.", param)
	default:
		return fmt.Sprintf("Invalid value. (Validation rule: %s)", tag)
	}
}

// You could add other utility functions here as your project grows.
// For example:
// - Random string generators
// - Slugify functions
// - Date/Time manipulation helpers (though standard library is often enough)

/*
Example of adding a custom validator (if you used `e164` for phone numbers):

import "github.com/nyaruka/phonenumbers"

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(...) // as before

	validate.RegisterValidation("e164", func(fl validator.FieldLevel) bool {
		// fl.Field() gives access to the field value
		numStr := fl.Field().String()
		if numStr == "" { // Allow empty if not 'required'
			return true
		}
		// Parse the number. The second argument is the default region.
		// For E.164, the region might not be strictly necessary if the number includes the country code.
		// However, providing a default region can help parse numbers that might be missing it.
		// If your numbers are always expected to have a '+', you might not need a default region.
		num, err := phonenumbers.Parse(numStr, "") // Use "" for region if E.164 is strictly enforced with '+'
		if err != nil {
			return false
		}
		// Check if the number is valid and in E.164 format
		return phonenumbers.IsValidNumber(num) &&
		       phonenumbers.Format(num, phonenumbers.E164) == numStr
	})
}

To use this, you would need: go get github.com/nyaruka/phonenumbers
*/