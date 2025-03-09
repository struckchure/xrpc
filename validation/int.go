package validation

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/samber/lo"
)

type IntValidator struct {
	rules []Rule
}

func (i *IntValidator) isEmpty(val reflect.Value) bool {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return true // Nil pointer is considered empty
		}
	}

	return false
}

func (i *IntValidator) isInt(val reflect.Value) bool {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false // Not an integer type
	}
}

func (i *IntValidator) Required(message ...string) *IntValidator {
	if len(message) == 0 {
		message = append(message, "field is required")
	}

	i.TypeCheck()

	i.rules = addRule(i.rules, Rule{
		name: "Required",
		callback: func(v ...any) error {
			val := v[0].(reflect.Value)
			if i.isEmpty(val) {
				return errors.New(message[0])
			}

			return nil
		},
	})

	return i
}

func (i *IntValidator) TypeCheck(message ...string) *IntValidator {
	if len(message) == 0 {
		message = append(message, "input must be a number")
	}

	i.rules = addRule(i.rules, Rule{
		name: "TypeCheck",
		callback: func(v ...any) error {
			val := v[0].(reflect.Value)

			if !i.isInt(val) {
				return errors.New(message[0])
			}

			return nil
		},
	})

	return i
}

func (i *IntValidator) Min(value int, message ...string) *IntValidator {
	if len(message) == 0 {
		message = append(message, fmt.Sprintf("min value required is %d", value))
	}

	i.rules = addRule(i.rules, Rule{
		name: "Min",
		callback: func(v ...any) error {
			if v[0].(int) < value {
				return errors.New(message[0])
			}
			return nil
		},
	})

	return i
}

func (i *IntValidator) Max(value int, message ...string) *IntValidator {
	if len(message) == 0 {
		message = append(message, fmt.Sprintf("max value required is %d", value))
	}

	i.rules = addRule(i.rules, Rule{
		name: "Max",
		callback: func(v ...any) error {
			if v[0].(int) > value {
				return errors.New(message[0])
			}
			return nil
		},
	})

	return i
}

func (i *IntValidator) Validate(value any) error {
	val := reflect.ValueOf(value)

	for _, rule := range i.rules {
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

func Int() *IntValidator {
	return &IntValidator{}
}
