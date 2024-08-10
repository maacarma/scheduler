package task

type MapAny map[string]any

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
