package utils

import (
	"fmt"
	"reflect"
	errors "webdp/internal/api/http"
)

const (
	NON_EMPTY = "non-empty-string"
)

/*
	OBS OBS OBS
	important to use the right tag
	this validation works only for STRING fields

*/

func ValidateNonEmptyString(obj any) error {

	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("dpvalidation")
		if tag != "" && tag != NON_EMPTY {
			panic("WRONG TAG")
		}
		val := v.FieldByName(field.Name)
		if val.Interface() == reflect.ValueOf("").Interface() {
			return fmt.Errorf("%w: field \"%s\" in struct \"%s\" should not be empty", errors.ErrBadFormatting, field.Name, t.Name())
		}

	}
	return nil
}
