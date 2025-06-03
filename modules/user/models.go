package user

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/dig"
)

type Module struct {
	Controller *UserController
}

func NewModule(container *dig.Container) *Module {
	container.Provide(NewUserService)
	container.Provide(NewUserController)

	var controller *UserController
	_ = container.Invoke(func(c *UserController) {
		controller = c
	})

	return &Module{Controller: controller}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	m.Controller.RegisterRoutes(r)
}
