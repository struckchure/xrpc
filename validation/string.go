package validation

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/samber/lo"
)

type StringValidator struct {
	rules []Rule
}

func (s *StringValidator) isEmpty(val reflect.Value) bool {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return true // Nil pointer is considered empty
		}
		val = val.Elem()
	}

	return len(val.String()) == 0
}

func (s *StringValidator) isString(val reflect.Value) bool {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val.Kind() == reflect.String
}

func (s *StringValidator) Required(message ...string) *StringValidator {
	if len(message) == 0 {
		message = append(message, "field is required")
	}

	s.TypeCheck()

	s.rules = addRule(s.rules, Rule{
		name: "Required",
		callback: func(v ...any) error {
			val := v[0].(reflect.Value)
			if s.isEmpty(val) {
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
			if !s.isString(val) {
				return errors.New(message[0])
			}

			return nil
		},
	})

	return s
}

func (s *StringValidator) Length(l int, message ...string) *StringValidator {
	if len(message) == 0 {
		message = append(message, fmt.Sprintf("length required is %d", l))
	}

	s.rules = addRule(s.rules, Rule{
		name: "Length",
		callback: func(v ...any) error {
			if len(v[0].(string)) != l {
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

func (s *StringValidator) Regex(pattern string, message ...string) *StringValidator {
	if len(message) == 0 {
		message = append(message, fmt.Sprintf("string does not match regex: %s", pattern))
	}

	s.rules = addRule(s.rules, Rule{
		name: "Regex",
		callback: func(v ...any) error {
			matched, err := regexp.MatchString(pattern, v[0].(string))
			if err != nil {
				return fmt.Errorf("invalid regex pattern: %w", err)
			}
			if !matched {
				return errors.New(message[0])
			}
			return nil
		},
	})

	return s
}

func (s *StringValidator) Email(message ...string) *StringValidator {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if len(message) == 0 {
		message = append(message, "invalid email format")
	}

	_, isRequired := lo.Find(s.rules, func(r Rule) bool { return r.name == "Required" })
	if !isRequired {
		return s
	}

	return s.Regex(emailRegex, message...)
}

func (s *StringValidator) Validate(value any) error {
	val := reflect.ValueOf(value)

	for _, rule := range s.rules {
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
