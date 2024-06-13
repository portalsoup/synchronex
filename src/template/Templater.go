package template

import (
	"reflect"
	"strings"
)

type Template struct {
	User string `token:"{{USER}}"`
}

func (t Template) Replace(statement string) string {
	val := reflect.ValueOf(t)
	typ := reflect.TypeOf(t)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := typ.Field(i).Tag.Get("token")
		if tag != "" {
			statement = strings.Replace(statement, tag, field.String(), -1)
		}
	}

	return statement
}
