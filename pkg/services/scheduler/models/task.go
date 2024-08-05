package task

import (
	errors "github.com/maacarma/scheduler/pkg/errors"
	utils "github.com/maacarma/scheduler/utils"
)

type MapAny map[string]any

// valid HTTP methods for a task
const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	PATCH  = "PATCH"
)

var methods = []string{GET, POST, PUT, DELETE, PATCH}

const (
	json = "json"
)

type Task struct {
	ID        int64  `json:"id"`
	Url       string `json:"url"`
	Method    string `json:"method"`
	Namespace string `json:"namespace"`
	Headers   MapAny `json:"headers"`
	Body      MapAny `json:"body"`
}

type TaskPayload struct {
	Url       string `json:"url"`
	Method    string `json:"method"`
	Namespace string `json:"namespace"`
	Headers   MapAny `json:"headers"`
	Body      MapAny `json:"body"`
}

// Validate validates the task payload
func (t *TaskPayload) Validate() *errors.Validation {
	if t.Url == "" {
		name := utils.GetStructTag(t, "Url", json)
		return errors.InvalidPayload(name, "url is required")
	}
	if utils.Contains(methods, t.Method) {
		name := utils.GetStructTag(t, "Method", json)
		return errors.InvalidPayload(name, "method is required")
	}
	if t.Namespace == "" {
		name := utils.GetStructTag(t, "Namespace", json)
		return errors.InvalidPayload(name, "namespace is required")
	}
	return nil
}
