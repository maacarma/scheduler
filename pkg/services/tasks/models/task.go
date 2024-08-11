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

// Task represents a task entity.
type Task struct {
	ID        string `json:"_id" bson:"_id"`
	Url       string `json:"url" bson:"url"`
	Method    string `json:"method" bson:"method"`
	Namespace string `json:"namespace" bson:"namespace"`
	Params    MapAny `json:"params" bson:"params"`
	Headers   MapAny `json:"headers" bson:"headers"`
	Body      MapAny `json:"body" bson:"body"`
}

// TaskPayload is the api payload schema for creating a task.
type TaskPayload struct {
	Url       string `json:"url" bson:"url"`
	Method    string `json:"method" bson:"method"`
	Namespace string `json:"namespace" bson:"namespace"`
	Params    MapAny `json:"params" bson:"params"`
	Headers   MapAny `json:"headers" bson:"headers"`
	Body      MapAny `json:"body" bson:"body"`
}

func (t *TaskPayload) Validate() *errors.Validation {
	if t.Url == "" {
		return errors.InvalidPayload("url", errors.RequiredFieldMsg)
	}
	if !utils.Contains(methods, t.Method) {
		return errors.InvalidPayload("method", "method is invalid")
	}
	return nil
}
