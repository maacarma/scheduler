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

// Contains checks if a value is in the slice
func Contains[T comparable](s []T, v T) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}
	return false
}

// ConvertToCronInterval converts a string interval to a cron interval
// Ex: 1h -> @every 1g
func ConvertToCronInterval(interval string) string {
	return fmt.Sprintf("@every %s", interval)
}
