package utils

import (
	"encoding/json"
	"fmt"
)

func Add(a, b int) int { return a + b }

// PrintStruct prints a givens struct in pretty format with indent
func PrintStruct(v any) {
	res, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(res))
}
