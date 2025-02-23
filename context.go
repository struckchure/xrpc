package trpc

import "github.com/labstack/echo/v4"

type Context[T, R any] struct {
	ec echo.Context

	Input T
}

func (c *Context[T, R]) Json(status int, body R) error {
	return c.ec.JSON(status, body)
}

func (c *Context[T, R]) String(status int, body string) error {
	return c.ec.String(status, body)
}
