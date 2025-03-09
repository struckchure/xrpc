package xrpc

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/struckchure/xrpc/validation"
)

type ProcedureCallback[T, R any] func(Context[T, R]) error

type IProcedure[T, R any] interface {
	Input(*validation.Validator) IProcedure[T, R]
	Use(...ProcedureCallback[T, R]) IProcedure[T, R]
	Query(ProcedureCallback[T, R]) func(string, IApp)
	Mutation(ProcedureCallback[T, R]) func(string, IApp)
}

type Procedure[T, R any] struct {
	name        string
	validator   *validation.Validator
	ctx         Context[T, R]
	middlewares []ProcedureCallback[T, R]
}

func (p *Procedure[T, R]) Input(v *validation.Validator) IProcedure[T, R] {
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

	p.ctx.ec = c
	p.ctx.Input = input

	err := callback(p.ctx)
	if err != nil {
		switch err := err.(type) {
		case *XRPCError:
			return c.JSON(err.Code, echo.Map{"detail": err.Detail})
		}
	}
	return err
}

func (p *Procedure[T, R]) loadMiddlewares(app IApp) []echo.MiddlewareFunc {
	var middlewareFuncs []echo.MiddlewareFunc = []echo.MiddlewareFunc{}

	for _, middleware := range app.Ctx().rootMiddlewares {
		middlewareFunc := func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				app.Ctx(func(innerCtx Context[any, any]) Context[any, any] {
					innerCtx.ec = c
					return innerCtx
				})

				err := middleware(app.Ctx())
				if err != nil {
					switch err := err.(type) {
					case *XRPCError:
						return c.JSON(err.Code, echo.Map{"detail": err.Detail})
					default:
						return c.JSON(http.StatusInternalServerError, echo.Map{"detail": err})
					}
				}

				return next(c)
			}
		}
		middlewareFuncs = append(middlewareFuncs, middlewareFunc)
	}

	for _, middleware := range p.middlewares {
		middlewareFunc := func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				p.ctx.ec = c

				err := middleware(p.ctx)
				if err != nil {
					switch err := err.(type) {
					case *XRPCError:
						return c.JSON(err.Code, echo.Map{"detail": err.Detail})
					default:
						return c.JSON(http.StatusInternalServerError, echo.Map{"detail": err})
					}
				}
				return next(c)
			}
		}
		middlewareFuncs = append(middlewareFuncs, middlewareFunc)
	}

	return middlewareFuncs
}

func (p *Procedure[T, R]) Query(callback ProcedureCallback[T, R]) func(string, IApp) {
	return func(path string, app IApp) {
		p.ctx.Injector = app.Ctx().Injector
		p.ctx.sharedValue = app.Ctx().sharedValue
		p.ctx.rootMiddlewares = app.Ctx().rootMiddlewares

		path = JoinPath(path, p.name)
		path = app.Get(Route{
			path:        path,
			handler:     func(c echo.Context) error { return p.handler(c, callback) },
			middlewares: p.loadMiddlewares(app),
		})

		app.Spec(func(spec TRPCSpec) TRPCSpec {
			spec.Procedures = append(spec.Procedures, XRPCSpecProcedure{
				Path:   path,
				Type:   XRPCSpecProcedureTypeQuery,
				Input:  createTypeDescriptor[T](),
				Output: createTypeDescriptor[R](),
			})

			return spec
		})

		fmt.Printf("[xRPC] [%s] %s\n", XRPCSpecProcedureTypeQuery, path)
	}
}

func (p *Procedure[T, R]) Mutation(callback ProcedureCallback[T, R]) func(string, IApp) {
	return func(path string, app IApp) {
		p.ctx.Injector = app.Ctx().Injector
		p.ctx.sharedValue = app.Ctx().sharedValue
		p.ctx.rootMiddlewares = app.Ctx().rootMiddlewares

		path = JoinPath(path, p.name)
		path = app.Post(Route{
			path:        path,
			handler:     func(c echo.Context) error { return p.handler(c, callback) },
			middlewares: p.loadMiddlewares(app),
		})

		app.Spec(func(spec TRPCSpec) TRPCSpec {
			spec.Procedures = append(spec.Procedures, XRPCSpecProcedure{
				Path:   path,
				Type:   XRPCSpecProcedureTypeMutation,
				Input:  createTypeDescriptor[T](),
				Output: createTypeDescriptor[R](),
			})

			return spec
		})

		fmt.Printf("[xRPC] [%s] %s\n", XRPCSpecProcedureTypeMutation, path)
	}
}

func NewProcedure[T, R any](name string) IProcedure[T, R] {
	return &Procedure[T, R]{name: StripSlash(name)}
}
