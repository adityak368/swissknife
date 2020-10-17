package playground

import (
	"net/http"

	"github.com/adityak368/swissknife/response"
	"github.com/adityak368/swissknife/validation"
	"github.com/go-playground/validator/v10"
)

// GoPlaygroundValidator Request Validation. Currently Empty But can add any Validator
type GoPlaygroundValidator struct {
	validator *validator.Validate
}

// Validate Implements the validator interface
func (v *GoPlaygroundValidator) Validate(i interface{}) []error {
	err := v.validator.Struct(i)
	if err != nil {
		return toInternalError(err)
	}
	return nil
}

func toInternalError(err error) []error {

	var errors []error

	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			errors = append(errors, response.NewError(http.StatusBadRequest, "FieldRequired", err.Field()))
		case "min":
			errors = append(errors, response.NewError(http.StatusBadRequest, "MustHaveMinCharacters", err.Field(), err.Param()))
		case "max":
			errors = append(errors, response.NewError(http.StatusBadRequest, "MustHaveMaxCharacters", err.Field(), err.Param()))
		case "oneof":
			errors = append(errors, response.NewError(http.StatusBadRequest, "NotAValidValue", err.Field(), err.Param()))
		case "email":
			errors = append(errors, response.NewError(http.StatusBadRequest, "InvalidEmail", err.Field()))
		default:
			errors = append(errors, response.NewError(http.StatusBadRequest, "InvalidField", err.Field()))
		}
	}

	return errors
}

// New Creates a new go-playground validator
func New() validation.Validator {
	return &GoPlaygroundValidator{
		validator: validator.New(),
	}
}

// ValidateStruct is a helper for easy validation
func ValidateStruct(i interface{}) []error {
	validator := New()
	return validator.Validate(i)
}
