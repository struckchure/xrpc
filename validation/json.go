package validation

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/samber/lo"
)

type JsonValidator struct {
	rules []Rule
}

func (j *JsonValidator) isJsonEmpty(val reflect.Value) bool {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return true // Nil pointer is considered empty
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Map, reflect.Slice:
		return val.Len() == 0
	case reflect.String:
		var temp any
		err := json.Unmarshal([]byte(val.String()), &temp)
		if err != nil {
			return false // Invalid JSON string
		}

		// Check if the unmarshaled value is an empty map or slice
		if mapVal, ok := temp.(map[string]any); ok && len(mapVal) == 0 {
			return true
		}
		if sliceVal, ok := temp.([]any); ok && len(sliceVal) == 0 {
			return true
		}
		return false
	default:
		return false // Not a JSON-compatible type
	}
}

func (j *JsonValidator) isJson(val reflect.Value) bool {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return false
		}
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Map, reflect.Slice:
		return true
	case reflect.String:
		var temp any
		err := json.Unmarshal([]byte(val.String()), &temp)
		return err == nil
	default:
		return false
	}
}

func (j *JsonValidator) Required(message ...string) *JsonValidator {
	if len(message) == 0 {
		message = append(message, "field is required")
	}

	j.rules = addRule(j.rules, Rule{
		name: "Required",
		callback: func(v ...any) error {
			val := v[0].(reflect.Value)
			if j.isJsonEmpty(val) {
				return errors.New(message[0])
			}

			return nil
		},
	})

	j.TypeCheck()

	return j
}

func (j *JsonValidator) TypeCheck(message ...string) *JsonValidator {
	if len(message) == 0 {
		message = append(message, "input must be json")
	}

	j.rules = addRule(j.rules, Rule{
		name: "TypeCheck",
		callback: func(v ...any) error {
			val := v[0].(reflect.Value)
			if !j.isJson(val) {
				return errors.New(message[0])
			}

			return nil
		},
	})

	return j
}

func (j *JsonValidator) Validate(value any) error {
	val := reflect.ValueOf(value)

	for _, rule := range j.rules {
		if lo.Contains([]string{"Required", "TypeCheck"}, rule.name) {
			err := rule.callback(val)
			if err != nil {
				return err
			}
		} else {
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}

			if val.IsValid() {
				err := rule.callback(int(val.Int()))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func Json() *JsonValidator {
	return &JsonValidator{}
}
