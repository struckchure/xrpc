package xrpc

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"github.com/samber/lo"
)

type Context[T, R any] struct {
	ec          echo.Context
	next        func() error
	sharedValue map[string]any

	Injector *do.Injector
	Input    T
}

func (c *Context[T, R]) Header(key string) string {
	return c.ec.Request().Header.Get(key)
}

func (c *Context[T, R]) Json(status int, body R) error {
	return c.ec.JSON(status, body)
}

func (c *Context[T, R]) String(status int, body string) error {
	return c.ec.String(status, body)
}

func (c *Context[T, R]) Next() error {
	if c.next != nil {
		return c.next()
	}

	return nil
}

func (c *Context[T, R]) Locals(key string, value ...interface{}) interface{} {
	if len(value) > 0 {
		c.sharedValue[key] = value[0]
		return value[0]
	}

	if lo.HasKey(c.sharedValue, key) {
		return c.sharedValue[key]
	}

	return nil
}
