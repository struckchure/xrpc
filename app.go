package trpc

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type IApp interface {
	Server() *echo.Echo
	Ctx() Context[any, any]
	Use(...func(Context[any, any]) error) IApp
	Router(string, ...func(IApp, ...*echo.Group)) IApp
	Get(string, func(c echo.Context) error, ...*echo.Group)
	Post(string, func(c echo.Context) error, ...*echo.Group)
	Start(port int) error
}

type App struct {
	srv *echo.Echo

	ctx Context[any, any]
}

func (a *App) Server() *echo.Echo {
	return a.srv
}

func (a *App) Ctx() Context[any, any] {
	return a.ctx
}

func (a *App) Use(middlewares ...func(Context[any, any]) error) IApp {
	for _, m := range middlewares {
		a.srv.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				a.ctx.ec = c
				a.ctx.next = func() error {
					return next(c)
				}

				err := m(a.ctx)
				if err != nil {
					return err
				}

				return nil
			}
		})
	}

	return a
}

func (a *App) Router(path string, procedures ...func(IApp, ...*echo.Group)) IApp {
	group := a.srv.Group(path)

	for _, procedure := range procedures {
		procedure(a, group)
	}

	return a
}

func (a *App) Get(path string, handler func(c echo.Context) error, group ...*echo.Group) {
	if len(group) > 0 {
		group[0].GET(path, handler)
		return
	}

	a.srv.GET(path, handler)
}

func (a *App) Post(path string, handler func(c echo.Context) error, group ...*echo.Group) {
	if len(group) > 0 {
		group[0].POST(path, handler)
		return
	}

	a.srv.POST(path, handler)
}

func (a *App) Start(port int) error {
	return a.srv.Start(fmt.Sprintf(":%d", port))
}

func InitTRPC() IApp {
	srv := echo.New()

	srv.Pre(middleware.AddTrailingSlash())
	srv.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} - ${status} ${latency_human}\n",
	}))

	return &App{
		srv: srv,
		ctx: Context[any, any]{sharedValue: map[string]interface{}{}},
	}
}
