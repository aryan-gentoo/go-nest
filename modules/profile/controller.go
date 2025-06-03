package profile

import (
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
)

type ProfileController struct {
	service *ProfileService
}

func NewProfileController(service *ProfileService) *ProfileController {
	return &ProfileController{service}
}

func (c *ProfileController) RegisterRoutes(r chi.Router) {
	r.Route("/profile", func(r chi.Router) {
		r.Get("/", c.GetAll)
	})
}

func (c *ProfileController) GetAll(w http.ResponseWriter, r *http.Request) {
	data := []string{"example"}
	json.NewEncoder(w).Encode(data)
}
