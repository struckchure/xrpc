package trpc

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ProcedureHandler func(IApp, ...*echo.Group)

type IProcedure[T, R any] interface {
	Input(*Validator) IProcedure[T, R]
	Use(...func(Context[T, R]) error) IProcedure[T, R]
	Query(func(Context[T, R]) error) ProcedureHandler
	Mutation(func(Context[T, R]) error) ProcedureHandler
}

type Procedure[T, R any] struct {
	name      string
	validator *Validator
	ctx       Context[T, R]
}

func (p *Procedure[T, R]) Input(v *Validator) IProcedure[T, R] {
	if v != nil {
		p.validator = v
	}

	return p
}

func (p *Procedure[T, R]) Use(middlewares ...func(Context[T, R]) error) IProcedure[T, R] {
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

	p.ctx.ec = c
	p.ctx.Input = input

	return callback(p.ctx)
}

func (p *Procedure[T, R]) Query(callback func(Context[T, R]) error) ProcedureHandler {
	return func(t IApp, groups ...*echo.Group) {
		p.ctx.sharedValue = t.Ctx().sharedValue

		path := t.Get(p.name, func(c echo.Context) error { return p.handler(c, callback) }, groups...)

		t.Spec(func(spec TRPCSpec) TRPCSpec {
			spec.Procedures = append(spec.Procedures, TRPCSpecProcedure{
				Path:   path,
				Type:   TRPCSpecProcedureTypeQuery,
				Input:  createTypeDescriptor[T](),
				Output: createTypeDescriptor[R](),
			})

			return spec
		})
	}
}

func (p *Procedure[T, R]) Mutation(callback func(Context[T, R]) error) ProcedureHandler {
	return func(t IApp, groups ...*echo.Group) {
		p.ctx.sharedValue = t.Ctx().sharedValue

		path := t.Post(p.name, func(c echo.Context) error { return p.handler(c, callback) }, groups...)

		t.Spec(func(spec TRPCSpec) TRPCSpec {
			spec.Procedures = append(spec.Procedures, TRPCSpecProcedure{
				Path:   path,
				Type:   TRPCSpecProcedureTypeMutation,
				Input:  createTypeDescriptor[T](),
				Output: createTypeDescriptor[R](),
			})

			return spec
		})
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
