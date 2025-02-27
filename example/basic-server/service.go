package main

import (
	"fmt"

	"github.com/samber/do"
)

type EngineService interface{}

func NewEngineService(i *do.Injector) (EngineService, error) {
	return &engineServiceImplem{}, nil
}

type engineServiceImplem struct{}

// [Optional] Implements do.Healthcheckable.
func (c *engineServiceImplem) HealthCheck() error {
	return fmt.Errorf("engine broken")
}

func NewCarService(i *do.Injector) (*CarService, error) {
	engine := do.MustInvoke[EngineService](i)
	car := CarService{Engine: engine}
	return &car, nil
}

type CarService struct {
	Engine EngineService
}

func (c *CarService) Start() {
	println("car starting")
}

// [Optional] Implements do.Shutdownable.
func (c *CarService) Shutdown() error {
	println("car stopped")
	return nil
}
