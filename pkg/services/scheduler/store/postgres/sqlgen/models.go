// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sqlgen

type Task struct {
	ID        int64  `json:"id"`
	Url       string `json:"url"`
	Method    string `json:"method"`
	Namespace string `json:"namespace"`
	Headers   []byte `json:"headers"`
	Body      []byte `json:"body"`
}
