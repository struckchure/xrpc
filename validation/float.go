package validation

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/samber/lo"
)

type FloatValidator struct {
	rules []Rule
}

func (f *FloatValidator) isEmpty(val reflect.Value) bool {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return true // Nil pointer is considered empty
		}
	}

	return false
}

func (f *FloatValidator) isFloat(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false // Not a float type
	}
}

func (f *FloatValidator) Required(message ...string) *FloatValidator {
	if len(message) == 0 {
		message = append(message, "field is required")
	}

	f.TypeCheck()

	f.rules = addRule(f.rules, Rule{
		name: "Required",
		callback: func(v ...any) error {
			val := v[0].(reflect.Value)
			if f.isEmpty(val) {
				return errors.New(message[0])
			}

			return nil
		},
	})

	return f
}

func (f *FloatValidator) TypeCheck(message ...string) *FloatValidator {
	if len(message) == 0 {
		message = append(message, "input must be a float")
	}

	f.rules = addRule(f.rules, Rule{
		name: "TypeCheck",
		callback: func(v ...any) error {
			val := v[0].(reflect.Value)

			if !f.isFloat(val) {
				return errors.New(message[0])
			}

			return nil
		},
	})

	return f
}

func (f *FloatValidator) Min(value float64, message ...string) *FloatValidator {
	if len(message) == 0 {
		message = append(message, fmt.Sprintf("min value required is %f", value))
	}

	f.rules = addRule(f.rules, Rule{
		name: "Min",
		callback: func(v ...any) error {
			if v[0].(float64) < value {
				return errors.New(message[0])
			}
			return nil
		},
	})

	return f
}

func (f *FloatValidator) Max(value float64, message ...string) *FloatValidator {
	if len(message) == 0 {
		message = append(message, fmt.Sprintf("max value required is %f", value))
	}

	f.rules = addRule(f.rules, Rule{
		name: "Max",
		callback: func(v ...any) error {
			if v[0].(float64) > value {
				return errors.New(message[0])
			}
			return nil
		},
	})

	return f
}

func (f *FloatValidator) Validate(value any) error {
	val := reflect.ValueOf(value)

	for _, rule := range f.rules {
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
				err := rule.callback(val.Float())
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func Float() *FloatValidator {
	return &FloatValidator{}
}
