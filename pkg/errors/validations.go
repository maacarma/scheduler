package errors

type Validation struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Error       string `json:"error,omitempty"`
}

const (
	RequiredFieldMsg = "requires field"
	InvalidFieldMsg  = "invalid field"
)

// ErrInvalidTaskPayload is the error returned when the payload is invalid
// Expects args array of strings in following format
// args[0]: key of the field
// args[1]: description specified whats wrong with the field
// args[2]: any additional error message
func InvalidPayload(args ...string) *Validation {
	switch len(args) {
	case 2:
		return &Validation{Key: args[0], Description: args[1]}
	case 3:
		return &Validation{Key: args[0], Description: args[1], Error: args[2]}
	default:
		return &Validation{}
	}
}
