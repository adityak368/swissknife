package validation

// Validator defines a validator's interface
type Validator interface {
	Validate(i interface{}) []error
}
