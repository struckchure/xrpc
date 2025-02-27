package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/samber/lo"
)

type Validator struct {
	fields map[string]any
}

func (v *Validator) Field(field string, validator any) *Validator {
	if v.fields == nil {
		v.fields = make(map[string]any)
	}
	v.fields[field] = validator
	return v
}

func (v *Validator) Validate(input any) map[string]string {
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
