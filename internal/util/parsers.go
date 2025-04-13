package util

import (
	"reflect"
	"strings"
)

func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)

	// If obj is a pointer, get the value it points to
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ensure it's a struct
	if val.Kind() != reflect.Struct {
		panic("Expected struct, got " + val.Kind().String())
	}

	// Iterate over all the struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Name

		// Handle if the field is exported (public)
		if field.CanInterface() {
			result[fieldName] = field.Interface()
		}
	}

	return result
}

func ParseHyperJumpPath(p string) []string {
	path := strings.Split(p, "->")

	return path
}
