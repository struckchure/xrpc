package trpc

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type IProcedure[T, R any] interface {
	Input(*Validator) IProcedure[T, R]
	Query(func(Context[T, R]) error) func(IApp)
	Mutation(func(Context[T, R]) error) func(IApp)
}

type Procedure[T, R any] struct {
	name      string
	validator *Validator
}

func (p *Procedure[T, R]) Input(v *Validator) IProcedure[T, R] {
	if v != nil {
		p.validator = v
	}

	return p
}

func (p *Procedure[T, R]) handler(c echo.Context, callback func(Context[T, R]) error) error {
	var input T

	if p.validator != nil {
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
		}

		if err := p.validator.Validate(input); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
	}

	return callback(Context[T, R]{ec: c, Input: input})
}

func (p *Procedure[T, R]) Query(callback func(Context[T, R]) error) func(IApp) {
	return func(t IApp) {
		t.Get(p.name, func(c echo.Context) error { return p.handler(c, callback) })

		fmt.Printf("Mapped Query: %s\n", p.name)
	}
}

func (p *Procedure[T, R]) Mutation(callback func(Context[T, R]) error) func(IApp) {
	return func(t IApp) {
		t.Post(p.name, func(c echo.Context) error { return p.handler(c, callback) })

		fmt.Printf("Mapped Mutation: %s\n", p.name)
	}
}

func NewProcedure[T, R any](name string) IProcedure[T, R] {
	if !strings.HasPrefix(name, "/") {
		name = "/" + name
	}

	if !strings.HasSuffix(name, "/") {
		name = name + "/"
	}

	return &Procedure[T, R]{name: name}
}
