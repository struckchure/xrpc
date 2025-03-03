package xrpc

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/do"
	"gopkg.in/yaml.v3"
)

type IApp interface {
	Injector() *do.Injector
	Spec(modifier func(TRPCSpec) TRPCSpec)
	GenerateSpec()
	Server() *echo.Echo
	Ctx(...func(Context[any, any]) Context[any, any]) Context[any, any]
	Middlewares() []ProcedureCallback[any, any]
	Use(...ProcedureCallback[any, any]) IApp
	Router(string, ...func(IApp, ...*echo.Group)) IApp
	Get(string, func(c echo.Context) error, ...*echo.Group) string
	Post(string, func(c echo.Context) error, ...*echo.Group) string
	Start(port int) error
}

type App struct {
	spec        TRPCSpec
	autoGenSpec bool
	specPath    string
	injector    *do.Injector
	srv         *echo.Echo
	ctx         Context[any, any]
	middlewares []ProcedureCallback[any, any]
}

func (a *App) Injector() *do.Injector {
	return a.injector
}

func (a *App) Server() *echo.Echo {
	return a.srv
}

func (a *App) Ctx(modifiers ...func(Context[any, any]) Context[any, any]) Context[any, any] {
	for _, modifier := range modifiers {
		a.ctx = modifier(a.ctx)
	}

	return a.ctx
}

func (a *App) Spec(modifier func(TRPCSpec) TRPCSpec) {
	a.spec = modifier(a.spec)
}

func (a *App) GenerateSpec() {
	a.Spec(func(t TRPCSpec) TRPCSpec {
		yamlData, err := yaml.Marshal(&t)
		if err != nil {
			log.Fatalf("Error marshaling YAML: %v", err)
		}
		err = WriteFile(a.specPath, string(yamlData))
		if err != nil {
			log.Fatalln(err)
		}

		return t
	})
}

func (a *App) Router(path string, procedures ...func(IApp, ...*echo.Group)) IApp {
	group := a.srv.Group(path)

	for _, procedure := range procedures {
		procedure(a, group)
	}

	a.middlewares = []ProcedureCallback[any, any]{}

	return a
}

func (a *App) Middlewares() []ProcedureCallback[any, any] {
	return a.middlewares
}

func (a *App) Use(middlewares ...ProcedureCallback[any, any]) IApp {
	a.middlewares = append(a.middlewares, middlewares...)
	return a
}

func (a *App) Get(path string, handler func(c echo.Context) error, group ...*echo.Group) string {
	if len(group) > 0 {
		return group[0].GET(path, handler).Path
	}

	return a.srv.GET(path, handler).Path
}

func (a *App) Post(path string, handler func(c echo.Context) error, group ...*echo.Group) string {
	if len(group) > 0 {
		return group[0].POST(path, handler).Path
	}

	return a.srv.POST(path, handler).Path
}

func (a *App) Start(port int) error {
	if a.autoGenSpec {
		a.GenerateSpec()
	}

	return a.srv.Start(fmt.Sprintf(":%d", port))
}

type XRPCConfig struct {
	Name            string
	ServerUrl       string
	AutoGenTRPCSpec bool
	SpecPath        string
}

func NewXRPC(cfg ...XRPCConfig) IApp {
	_cfg := XRPCConfig{
		Name:            "xRPC Spec",
		ServerUrl:       "http://localhost:9090",
		AutoGenTRPCSpec: true,
		SpecPath:        "xrpc.yaml",
	}

	if len(cfg) > 0 {
		_cfg = cfg[0]
	}

	srv := echo.New()

	srv.Pre(middleware.AddTrailingSlash())
	srv.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} - ${status} ${latency_human}\n",
	}))

	i := do.New()

	return &App{
		spec: TRPCSpec{
			Name:      _cfg.Name,
			ServerUrl: _cfg.ServerUrl,
		},
		autoGenSpec: _cfg.AutoGenTRPCSpec,
		specPath:    _cfg.SpecPath,
		injector:    i,
		srv:         srv,
		ctx: Context[any, any]{
			sharedValue: map[string]any{},
			Injector:    i,
		},
		middlewares: []ProcedureCallback[any, any]{},
	}
}
