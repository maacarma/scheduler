package errors

type Validation struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}

// defined validation messages
const (
	RequiredFieldMsg = "requires the field"
)

// ErrInvalidTaskPayload is the error returned when a task payload is invalid
func InvalidPayload(key string, desc string) *Validation {
	return &Validation{
		Key:         key,
		Description: desc,
	}
}
