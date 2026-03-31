package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator represents a wrapper for go-playground validator
type Validator struct {
	validate *validator.Validate
}

// New returns a new instance of Validator with initialized properties
func New() *Validator {
	v := validator.New()

	// Register structural tag function to retrieve the json tag of the field.
	// This will make validator.FieldError.Field() return the JSON tag instead of the Go struct field name.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return fld.Name
		}
		return name
	})

	return &Validator{validate: v}
}

// Validate executes validation on struct input gracefully formatting to Indonesian language custom messages
func (v *Validator) Validate(i interface{}) map[string]string {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			// The key will be format snake_case from JSON tag
			errors[e.Field()] = msgForTag(e)
		}
	} else {
		// Handle the case where the error isn't a validation error
		errors["general"] = "Kesalahan internal pada validasi payload"
	}

	return errors
}

// msgForTag returns Indonesian human-readable string based on the current field error tag
func msgForTag(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "wajib diisi"
	case "min":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("minimal %s karakter", err.Param())
		}
		return fmt.Sprintf("minimal %s", err.Param())
	case "max":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("maksimal %s karakter", err.Param())
		}
		return fmt.Sprintf("maksimal %s", err.Param())
	case "email":
		return "format email tidak valid"
	case "oneof":
		return fmt.Sprintf("harus salah satu dari: %s", err.Param())
	case "uuid4":
		return "format UUID tidak valid"
	case "len":
		return fmt.Sprintf("harus tepat %s karakter", err.Param())
	case "numeric":
		return "harus berupa angka"
	default:
		return "tidak valid"
	}
}
