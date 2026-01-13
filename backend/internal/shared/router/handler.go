package router

import (
	"alloy/internal/shared/config"
	"alloy/internal/shared/constants"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	
	
)

// ModuleServices holds all module services that can be used by handlers
type ModuleServices struct {
	services map[string]interface{}
}

func NewModuleServices() *ModuleServices {
	return &ModuleServices{
		services: make(map[string]interface{}),
	}
}

func (ms *ModuleServices) Register(name string, service interface{}) {
	ms.services[name] = service
}

func (ms *ModuleServices) Get(name string) interface{} {
	return ms.services[name]
}

// Environment provides shared context for all handlers
type Environment struct {
	Config   *config.Config
	Fiber    *fiber.App
	Logger   *zap.Logger
	Stores *constants.DataStores
	Services *ModuleServices
}

// IHandler defines the interface that all module handlers must implement
type IHandler interface {
	Init(basePath string, env *Environment) error
}

func NewEnvironment(cfg *config.Config, fiberApp *fiber.App, 
	logger *zap.Logger, stores *constants.DataStores, services *ModuleServices) *Environment {
	return &Environment{
		Config:   cfg,
		Fiber:    fiberApp,
		Logger:   logger,
		Stores: stores,
		Services: services,
	}
}

func InitHandlers(env *Environment, handlers []IHandler) error {
	for _, handler := range handlers {
		if err := handler.Init("/api", env); err != nil {
			return fmt.Errorf("failed to initialize handler: %v", err)
		}
	}

	return nil
}
