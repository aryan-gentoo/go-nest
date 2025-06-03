package profile

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/dig"
)

type Module struct {
	Controller *ProfileController
}

func NewModule(container *dig.Container) *Module {
	container.Provide(NewProfileService)
	container.Provide(NewProfileController)

	var controller *ProfileController
	_ = container.Invoke(func(c *ProfileController) {
		controller = c
	})

	return &Module{Controller: controller}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	m.Controller.RegisterRoutes(r)
}
