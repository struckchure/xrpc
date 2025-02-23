package trpc

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type IApp interface {
	Get(string, func(c echo.Context) error)
	Post(string, func(c echo.Context) error)
	Start(port int) error
}

type App struct {
	srv *echo.Echo
}

func (a *App) Get(path string, handler func(c echo.Context) error) {
	a.srv.GET(path, handler)
}

func (a *App) Post(path string, handler func(c echo.Context) error) {
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

	return &App{srv: srv}
}
