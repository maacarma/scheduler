package task

type MapAny map[string]any

type Task struct {
	ID        int64  `json:"id"`
	Url       string `json:"url"`
	Method    string `json:"method"`
	Namespace string `json:"namespace"`
	Params    MapAny `json:"params"`
	Headers   MapAny `json:"headers"`
	Body      MapAny `json:"body"`
}

type TaskPayload struct {
	Url       string `json:"url"`
	Method    string `json:"method"`
	Namespace string `json:"namespace"`
	Params    MapAny `json:"params"`
	Headers   MapAny `json:"headers"`
	Body      MapAny `json:"body"`
}
