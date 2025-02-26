package xrpc

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/samber/lo"
)

type Rule struct {
	name     string
	callback func(...any) error
}

func addRule(rules []Rule, r Rule) []Rule {
	_, idx, exists := lo.FindIndexOf(rules, func(r1 Rule) bool { return r1.name == r.name })
	if exists {
		rules[idx] = r
	} else {
		rules = append(rules, r)
	}

	return rules
}

func removeRule(rules []Rule, name string) []Rule {
	_, idx, exists := lo.FindIndexOf(rules, func(r1 Rule) bool { return r1.name == name })
	if exists {
		rules = append(rules[:idx], rules[idx+1:]...)
	}

	return rules
}

type StringValidator struct {
	rules []Rule
}

func (s *StringValidator) Required(message ...string) *StringValidator {
	if len(message) == 0 {
		message = append(message, "field is required")
	}

	s.rules = addRule(s.rules, Rule{
		name: "Required",
		callback: func(v ...any) error {
			val := v[0].(reflect.Value)
			if !val.IsValid() {
				return errors.New(message[0])
			}
			return nil
		},
	})

	return s
}

func (s *StringValidator) TypeCheck(message ...string) *StringValidator {
	if len(message) == 0 {
		message = append(message, "input must be a string")
	}

	s.rules = addRule(s.rules, Rule{
		name: "TypeCheck",
		callback: func(v ...any) error {
			val := v[0].(reflect.Value)
			if reflect.String != val.Kind() {
				return errors.New(message[0])
			}

			return nil
		},
	})

	return s
}

func (s *StringValidator) SkipTypeCheck() *StringValidator {
	s.rules = removeRule(s.rules, "TypeCheck")

	return s
}

func (s *StringValidator) Length(l int, message ...string) *StringValidator {
	if len(message) == 0 {
		message = append(message, fmt.Sprintf("length required is %d", l))
	}

	s.rules = addRule(s.rules, Rule{
		name: "Length",
		callback: func(v ...any) error {
			if len(v[0].(string)) == l {
				return errors.New(message[0])
			}
			return nil
		},
	})

	return s
}

func (s *StringValidator) MinLength(l int, message ...string) *StringValidator {
	if len(message) == 0 {
		message = append(message, fmt.Sprintf("minimum length required is %d", l))
	}

	s.rules = addRule(s.rules, Rule{
		name: "MinLength",
		callback: func(v ...any) error {
			if len(v[0].(string)) < l {
				return errors.New(message[0])
			}
			return nil
		},
	})

	return s
}

func (s *StringValidator) MaxLength(l int, message ...string) *StringValidator {
	if len(message) == 0 {
		message = append(message, fmt.Sprintf("maximum length required is %d", l))
	}

	s.rules = addRule(s.rules, Rule{
		name: "MaxLength",
		callback: func(v ...any) error {
			if len(v[0].(string)) > l {
				return errors.New(message[0])
			}
			return nil
		},
	})

	return s
}

func (s *StringValidator) Email(message ...string) *StringValidator {
	return s
}

func (s *StringValidator) Validate(value interface{}) error {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for _, rule := range s.rules {
		if lo.Contains([]string{"Required", "TypeCheck"}, rule.name) {
			err := rule.callback(val)
			if err != nil {
				return err
			}
		} else {
			if val.IsValid() {
				err := rule.callback(val.String())
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func String() *StringValidator {
	return &StringValidator{}
}

type NumberValidator struct {
	rules []Rule
}

func (n *NumberValidator) Required(message ...string) *NumberValidator {
	if len(message) == 0 {
		message = append(message, "field is required")
	}

	n.rules = addRule(n.rules, Rule{
		name: "Required",
		callback: func(v ...any) error {
			val := v[0].(reflect.Value)
			if !val.IsValid() {
				return errors.New(message[0])
			}
			return nil
		},
	})

	return n
}

func (n *NumberValidator) TypeCheck(message ...string) *NumberValidator {
	if len(message) == 0 {
		message = append(message, "input must be a number")
	}

	n.rules = addRule(n.rules, Rule{
		name: "TypeCheck",
		callback: func(v ...any) error {
			val := v[0].(reflect.Value)

			if !lo.Contains([]reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64}, val.Kind()) {
				return errors.New(message[0])
			}

			return nil
		},
	})

	return n
}

func (n *NumberValidator) SkipTypeCheck(message ...string) *NumberValidator {
	n.rules = removeRule(n.rules, "TypeCheck")

	return n
}

func (n *NumberValidator) Min(value int, message ...string) *NumberValidator {
	if len(message) == 0 {
		message = append(message, fmt.Sprintf("min value required is %d", value))
	}

	n.rules = addRule(n.rules, Rule{
		name: "Min",
		callback: func(v ...any) error {
			if v[0].(int) < value {
				return errors.New(message[0])
			}
			return nil
		},
	})

	return n
}

func (n *NumberValidator) Max(value int, message ...string) *NumberValidator {
	if len(message) == 0 {
		message = append(message, fmt.Sprintf("max value required is %d", value))
	}

	n.rules = addRule(n.rules, Rule{
		name: "Max",
		callback: func(v ...any) error {
			if v[0].(int) > value {
				return errors.New(message[0])
			}
			return nil
		},
	})

	return n
}

func (n *NumberValidator) Validate(value interface{}) error {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for _, rule := range n.rules {
		if lo.Contains([]string{"Required", "TypeCheck"}, rule.name) {
			err := rule.callback(val)
			if err != nil {
				return err
			}
		} else {
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

func Number() *NumberValidator {
	return &NumberValidator{}
}

type Validator struct {
	fields map[string]interface{}
}

func (v *Validator) Field(field string, validator interface{}) *Validator {
	if v.fields == nil {
		v.fields = make(map[string]interface{})
	}
	v.fields[field] = validator
	return v
}

func (v *Validator) Validate(input interface{}) map[string]string {
	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return map[string]string{"error": "input must be a struct"}
	}

	fieldErrors := map[string]string{}

	for name, _validator := range v.fields {
		field := val.FieldByName(name)
		if !field.IsValid() {
			fieldErrors[name] = fmt.Sprintf("no such field: %s in input", name)
			continue
		}

		// Get JSON tag name if available
		structField, _ := val.Type().FieldByName(name)
		jsonTag := structField.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = name
		} else {
			jsonTag = strings.Split(jsonTag, ",")[0]
		}

		validator := reflect.ValueOf(_validator)
		validateMethod := validator.MethodByName("Validate")
		if !validateMethod.IsValid() {
			fieldErrors[jsonTag] = fmt.Sprintf("no Validate method for field: %s", jsonTag)
			continue
		}

		results := validateMethod.Call([]reflect.Value{field})
		if len(results) > 0 && !results[0].IsNil() {
			fieldErrors[jsonTag] = results[0].Interface().(error).Error()
		}
	}

	if len(lo.Keys(fieldErrors)) > 0 {
		return fieldErrors
	}

	return nil
}

func NewValidator() *Validator {
	return &Validator{}
}
