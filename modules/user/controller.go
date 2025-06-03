package user

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserController struct {
	service *UserService
}

func NewUserController(service *UserService) *UserController {
	return &UserController{service}
}

func (c *UserController) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", c.GetAll)
	})
}

func (c *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	users := c.service.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}
