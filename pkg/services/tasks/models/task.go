package task

import (
	"time"

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
	StartUnix int64  `json:"start_unix" bson:"start_unix"`
	EndUnix   int64  `json:"end_unix" bson:"end_unix"`
	Interval  string `json:"interval" bson:"interval"`
	Paused    bool   `json:"paused" bson:"paused"`
}

// TaskPayload is the api payload schema for creating a task.
// Methods accepted are those defined in methods variable array.
// Interval is a string accepted by time.ParseDuration (http://golang.org/pkg/time/#ParseDuration).
type TaskPayload struct {
	Url       string `json:"url" bson:"url"`
	Method    string `json:"method" bson:"method"`
	Namespace string `json:"namespace" bson:"namespace"`
	Params    MapAny `json:"params" bson:"params"`
	Headers   MapAny `json:"headers" bson:"headers"`
	Body      MapAny `json:"body" bson:"body"`
	StartUnix int64  `json:"start_unix" bson:"start_unix"`
	EndUnix   int64  `json:"end_unix" bson:"end_unix"`
	Interval  string `json:"interval" bson:"interval"`
	Paused    bool   `json:"paused" bson:"paused"`
}

// Validate validates the task payload.
// checks if the task payload has all the required fields.
// checks if the task payload has any invalid fields. Ex: http method, interval.
func (t *TaskPayload) Validate() *errors.Validation {
	if t.Url == "" {
		return errors.InvalidPayload("url", errors.RequiredFieldMsg)
	}
	if !utils.Contains(methods, t.Method) {
		return errors.InvalidPayload("method", errors.InvalidFieldMsg)
	}
	_, err := time.ParseDuration(t.Interval)
	if err != nil {
		return errors.InvalidPayload("interval", errors.InvalidFieldMsg, err.Error())
	}

	return nil
}

func (t *TaskPayload) ConvertToTask(id string) Task {
	return Task{
		ID:        id,
		Url:       t.Url,
		Method:    t.Method,
		Namespace: t.Namespace,
		Params:    t.Params,
		Headers:   t.Headers,
		Body:      t.Body,
		StartUnix: t.StartUnix,
		EndUnix:   t.EndUnix,
		Interval:  t.Interval,
		Paused:    t.Paused,
	}
}

// IsActive checks if the task is active.
// A task is active if the current time is between the start and end time.
func (t *Task) IsActive() bool {
	if utils.CurrentUTCUnix() < t.StartUnix || utils.CurrentUTCUnix() > t.EndUnix {
		return false
	}

	return true
}
