package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Add(a, b int) int { return a + b }

func PrettyPrint(data interface{}) {
	var p []byte
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}

// GetStructTag returns the value of a tag in a struct
// if the s is not a struct or the pointer to a struct, it returns an empty string
// if the field does not exist, it returns an empty string
func GetStructTag(s interface{}, field string, tag string) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ""
	}

	r, found := t.FieldByName(field)
	if !found {
		return ""
	}

	return r.Tag.Get(tag)
}

// Contains checks if a value is in a slice
func Contains[T comparable](s []T, v T) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}
	return false
}
