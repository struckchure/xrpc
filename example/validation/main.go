package main

import (
	"encoding/json"
	"fmt"

	"github.com/struckchure/xrpc/validation"
)

type CreatePostInput struct {
	AuthorEmail string          `json:"author_email"`
	AuthorId    string          `json:"author_id"`
	Title       *string         `json:"title"`
	Content     string          `json:"content"`
	Likes       *int            `json:"likes"`
	Views       any             `json:"views"`
	Ratings     float32         `json:"ratings"`
	Ratings2    any             `json:"ratings_2"`
	Meta        map[string]any  `json:"meta"`
	Meta2       *map[string]any `json:"meta_2"`
}

func formatJson(data any) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func main() {
	v := validation.NewValidator().
		Field("AuthorEmail", validation.
			String().
			Email().
			Required(),
		).
		Field("AuthorId", validation.String().Required()).
		Field("Title", validation.String().Required()).
		Field("Content", validation.String().Required()).
		Field("Likes", validation.Int().Required()).
		Field("Views", validation.Int().Required()).
		Field("Ratings", validation.Float().Required()).
		Field("Ratings2", validation.Float().Required()).
		Field("Meta", validation.Json().Required()).
		Field("Meta2", validation.Json().TypeCheck())

	err := v.Validate(CreatePostInput{
		Views:    "123",
		Ratings2: 90,
	})
	if err != nil {
		s, _ := formatJson(err)
		fmt.Println(s)
	}
}
