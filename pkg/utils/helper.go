package utils

import (
	"reflect"
	"time"
)

func IsFieldInitialized(v interface{}, fieldName string) bool {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return false
	}

	structVal := val.Elem()
	structType := structVal.Type()

	for i := 0; i < structType.NumField(); i++ {
		if structType.Field(i).Name == fieldName {
			fieldVal := structVal.Field(i)

			if fieldVal.Kind() == reflect.Bool {
				return fieldVal.Bool()
			}

			if fieldVal.Kind() == reflect.Int {
				return fieldVal.Int() != 0
			}

			if fieldVal.Kind() == reflect.Struct && fieldVal.Type() == reflect.TypeOf(time.Time{}) {
				return !fieldVal.Interface().(time.Time).IsZero()
			}

			return fieldVal.Interface() != reflect.Zero(fieldVal.Type()).Interface()
		}
	}

	return false
}
