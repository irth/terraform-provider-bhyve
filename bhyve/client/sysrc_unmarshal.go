package client

import (
	"fmt"
	"reflect"
	"strings"
)

func has(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

func (s Sysrc) Unmarshal(target interface{}) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("UnmarshalSysrc: target is not a pointer to a struct")
	}

	// Get the type of the struct
	typ := val.Elem().Type()

	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("sysrc")

		tagSplit := strings.Split(tag, ",")
		if len(tagSplit) == 0 {
			continue
		}
		name := tagSplit[0]
		params := tagSplit[1:]

		field := val.Elem().Field(i)
		if !field.CanSet() || !field.CanInterface() {
			continue
		}

		val, err := s.Get(name)
		if err == ErrUnknownVariable && !has(params, "required") {
			// TODO(irth): add variable name to the error
			continue
		}
		if err != nil {
			return err
		}

		switch field.Kind() {
		case reflect.Bool:
			field.SetBool(strings.ToLower(val) == "yes")
		case reflect.String:
			field.SetString(val)
		default:
			return fmt.Errorf("UnmarshalSysrc: unsupported type %s", field.Kind())
		}
	}

	return nil
}
