package xrpc

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type TypeDescriptor struct {
	TypeName string            `yaml:"type_name,omitempty"`
	Fields   []FieldDescriptor `yaml:"fields,omitempty"`
	Nillable bool              `yaml:"nillable"`
	Array    *TypeDescriptor   `yaml:"array,omitempty"`
}

type FieldDescriptor struct {
	Name     string `yaml:"name"`
	Alias    string `yaml:"alias"`
	Type     string `yaml:"type"`
	Nillable bool   `yaml:"nillable"`
}

func createTypeDescriptor[T any]() TypeDescriptor {
	var t T
	typeOfT := reflect.TypeOf(t)

	if typeOfT == nil {
		// Handle nil case (e.g., return a default or error)
		return TypeDescriptor{
			TypeName: "nil",
			Nillable: true,
		}
	}

	// Dereference pointer if it is a pointer
	isNillable := typeOfT.Kind() == reflect.Ptr || typeOfT.Kind() == reflect.Interface

	if typeOfT.Kind() == reflect.Ptr {
		typeOfT = typeOfT.Elem() // Get the type the pointer points to.
	}

	descriptor := TypeDescriptor{
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
func createTypeDescriptorHelper(typeOfT reflect.Type) TypeDescriptor {
	isNillable := typeOfT.Kind() == reflect.Ptr || typeOfT.Kind() == reflect.Interface

	if typeOfT.Kind() == reflect.Ptr {
		typeOfT = typeOfT.Elem() // Get the type the pointer points to.
	}

	descriptor := TypeDescriptor{
		TypeName: typeOfT.Name(),
		Nillable: isNillable,
	}

	if typeOfT.Kind() == reflect.Struct {
		for i := range typeOfT.NumField() {
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

// WriteFile creates a new file and writes content to it.
// It returns an error if the file cannot be created or written to.
func WriteFile(filename string, content string) error {
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

func StripSlash(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	return path
}

func JoinPath(paths ...string) string {
	// Filter out empty strings
	var nonEmptyPaths []string
	for _, path := range paths {
		if path != "" {
			nonEmptyPaths = append(nonEmptyPaths, path)
		}
	}

	// Join paths with "/"
	joined := strings.Join(nonEmptyPaths, "/")

	// Remove duplicate slashes
	for strings.Contains(joined, "//") {
		joined = strings.ReplaceAll(joined, "//", "/")
	}

	// Ensure path starts with "/"
	if !strings.HasPrefix(joined, "/") {
		joined = "/" + joined
	}

	// Ensure path ends with "/"
	if !strings.HasSuffix(joined, "/") {
		joined = joined + "/"
	}

	return joined
}
