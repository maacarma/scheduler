package errors

type Validation struct {
	key         string
	description string
}

// ErrInvalidTaskPayload is the error returned when a task payload is invalid
func InvalidPayload(key, desc string) *Validation {
	return &Validation{
		key:         key,
		description: desc,
	}
}
