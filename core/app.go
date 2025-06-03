package core

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Module interface {
	RegisterRoutes(r chi.Router)
}

type App struct {
	router  chi.Router
	modules []Module
}

func NewApp(modules ...Module) *App {
	return &App{
		router:  chi.NewRouter(),
		modules: modules,
	}
}

func (a *App) Run(addr string) {
	for _, m := range a.modules {
		m.RegisterRoutes(a.router)
	}
	log.Println("Server started at", addr)
	http.ListenAndServe(addr, a.router)
}
