package trpc

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type TypeDescriptor[T any] struct {
	TypeName string               `yaml:"type_name,omitempty"`
	Fields   []FieldDescriptor    `yaml:"fields,omitempty"`
	Nillable bool                 `yaml:"nillable"`
	Array    *TypeDescriptor[any] `yaml:"array,omitempty"`
}

type FieldDescriptor struct {
	Name     string `yaml:"name"`
	Alias    string `yaml:"alias"`
	Type     string `yaml:"type"`
	Nillable bool   `yaml:"nillable"`
}

func createTypeDescriptor[T any]() TypeDescriptor[T] {
	var t T
	typeOfT := reflect.TypeOf(t)

	// Dereference pointer if it is a pointer
	isNillable := typeOfT.Kind() == reflect.Ptr || typeOfT.Kind() == reflect.Interface

	if typeOfT.Kind() == reflect.Ptr {
		typeOfT = typeOfT.Elem() // Get the type the pointer points to.
	}

	descriptor := TypeDescriptor[T]{
		TypeName: typeOfT.Name(),
		Nillable: isNillable,
	}

	if typeOfT.Kind() == reflect.Struct {
		for i := 0; i < typeOfT.NumField(); i++ {
			field := typeOfT.Field(i)
			fieldType := field.Type

			isFieldNillable := fieldType.Kind() == reflect.Ptr || fieldType.Kind() == reflect.Interface

			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			descriptor.Fields = append(descriptor.Fields, FieldDescriptor{
				Name:     field.Name,
				Type:     fieldType.String(),
				Alias:    getFieldAlias(field),
				Nillable: isFieldNillable,
			})
		}
	} else if typeOfT.Kind() == reflect.Slice || typeOfT.Kind() == reflect.Array {
		elementType := typeOfT.Elem()

		// Recursively call createTypeDescriptor for the element type
		elementDescriptor := createTypeDescriptorHelper(elementType)

		descriptor.Array = &elementDescriptor
	}

	return descriptor
}

// Helper function to make the recursive call work.
func createTypeDescriptorHelper(typeOfT reflect.Type) TypeDescriptor[any] {
	isNillable := typeOfT.Kind() == reflect.Ptr || typeOfT.Kind() == reflect.Interface

	if typeOfT.Kind() == reflect.Ptr {
		typeOfT = typeOfT.Elem() // Get the type the pointer points to.
	}

	descriptor := TypeDescriptor[any]{
		TypeName: typeOfT.Name(),
		Nillable: isNillable,
	}

	if typeOfT.Kind() == reflect.Struct {
		for i := 0; i < typeOfT.NumField(); i++ {
			field := typeOfT.Field(i)
			fieldType := field.Type

			isFieldNillable := fieldType.Kind() == reflect.Ptr || fieldType.Kind() == reflect.Interface

			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			descriptor.Fields = append(descriptor.Fields, FieldDescriptor{
				Name:     field.Name,
				Type:     fieldType.String(),
				Alias:    getFieldAlias(field),
				Nillable: isFieldNillable,
			})
		}
	} else if typeOfT.Kind() == reflect.Slice || typeOfT.Kind() == reflect.Array {
		elementType := typeOfT.Elem()

		// Recursively call createTypeDescriptor for the element type
		elementDescriptor := createTypeDescriptorHelper(elementType)

		descriptor.Array = &elementDescriptor
	}

	return descriptor
}

func getFieldAlias(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return field.Name
	}
	parts := strings.Split(tag, ",")
	return parts[0]
}

// writeFile creates a new file and writes content to it.
// It returns an error if the file cannot be created or written to.
func writeFile(filename string, content string) error {
	// Create a new file, or truncate an existing one.
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close() // Ensure the file is closed when the function returns

	// Write the content to the file.
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil // Success
}
