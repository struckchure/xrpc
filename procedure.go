package xrpc

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ProcedureHandler func(IApp, ...*echo.Group)
type ProcedureCallback[T, R any] func(Context[T, R]) error

type IProcedure[T, R any] interface {
	Input(*Validator) IProcedure[T, R]
	Use(...ProcedureCallback[T, R]) IProcedure[T, R]
	Query(ProcedureCallback[T, R]) ProcedureHandler
	Mutation(ProcedureCallback[T, R]) ProcedureHandler
}

type Procedure[T, R any] struct {
	name        string
	validator   *Validator
	ctx         Context[T, R]
	middlewares []ProcedureCallback[T, R]
}

func (p *Procedure[T, R]) Input(v *Validator) IProcedure[T, R] {
	if v != nil {
		p.validator = v
	}

	return p
}

func (p *Procedure[T, R]) Use(middlewares ...ProcedureCallback[T, R]) IProcedure[T, R] {
	p.middlewares = append(p.middlewares, middlewares...)

	return p
}

func (p *Procedure[T, R]) handler(c echo.Context, callback ProcedureCallback[T, R]) error {
	var input T

	if p.validator != nil {
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"detail": err})
		}

		if err := p.validator.Validate(input); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"detail": err})
		}
	}

	var middlewareErr error = nil
	proceedToNext := false

	p.ctx.ec = c
	p.ctx.Input = input
	p.ctx.next = func() error {
		proceedToNext = true
		return nil
	}

	for _, middleware := range p.middlewares {
		if middlewareErr != nil && !proceedToNext {
			break
		}

		err := middleware(p.ctx)
		if err != nil {
			middlewareErr = err
		}

		proceedToNext = false
	}

	if middlewareErr != nil {
		switch err := middlewareErr.(type) {
		case *XRPCError:
			return c.JSON(err.Code, echo.Map{"detail": err.Detail})
		default:
			return c.JSON(http.StatusInternalServerError, echo.Map{"detail": err})
		}
	}

	return callback(p.ctx)
}

func (p *Procedure[T, R]) Query(callback ProcedureCallback[T, R]) ProcedureHandler {
	return func(t IApp, groups ...*echo.Group) {
		p.ctx.Injector = t.Ctx().Injector
		p.ctx.sharedValue = t.Ctx().sharedValue

		path := t.Get(p.name, func(c echo.Context) error { return p.handler(c, callback) }, groups...)

		t.Spec(func(spec TRPCSpec) TRPCSpec {
			spec.Procedures = append(spec.Procedures, XRPCSpecProcedure{
				Path:   path,
				Type:   XRPCSpecProcedureTypeQuery,
				Input:  createTypeDescriptor[T](),
				Output: createTypeDescriptor[R](),
			})

			return spec
		})
	}
}

func (p *Procedure[T, R]) Mutation(callback ProcedureCallback[T, R]) ProcedureHandler {
	return func(t IApp, groups ...*echo.Group) {
		p.ctx.Injector = t.Ctx().Injector
		p.ctx.sharedValue = t.Ctx().sharedValue

		path := t.Post(p.name, func(c echo.Context) error { return p.handler(c, callback) }, groups...)

		t.Spec(func(spec TRPCSpec) TRPCSpec {
			spec.Procedures = append(spec.Procedures, XRPCSpecProcedure{
				Path:   path,
				Type:   XRPCSpecProcedureTypeMutation,
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
