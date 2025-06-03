package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "new":
		if len(os.Args) < 3 {
			fmt.Println("❌ Project name required")
			return
		}
		createNewProject(os.Args[2])
	case "g":
		if len(os.Args) < 4 || os.Args[2] != "module" {
			fmt.Println("❌ Usage: go-nest-cli g module <name>")
			return
		}
		generateModule(os.Args[3])
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Go Nest CLI")
	fmt.Println("Usage:")
	fmt.Println("  go-nest-cli new <project-name>       Create a new project")
	fmt.Println("  go-nest-cli g module <module-name>   Generate a module")
}

func createNewProject(name string) {
	fmt.Println("🚀 Creating new Go Nest project:", name)

	// Folder structure
	os.Mkdir(name, 0755)
	os.MkdirAll(name+"/modules/app", 0755)
	os.MkdirAll(name+"/core", 0755)

	// main.go
	mainGo := `package main

import (
	"` + name + `/core"
	"` + name + `/modules/app"
)

func main() {
	app.RegisterDependencies()
	core.Container.Invoke(func(c *app.Controller) {
		a := core.NewApp(c)
		a.Run(":8080")
	})
}
`
	writeFile(name+"/main.go", mainGo)

	// go.mod
	cmd := exec.Command("go", "mod", "init", name)
	cmd.Dir = name
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	// app.go
	appGo := `package core

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
`
	writeFile(name+"/core/app.go", appGo)

	// container.go
	containerGo := `package core

import "go.uber.org/dig"

var Container = dig.New()
`
	writeFile(name+"/core/container.go", containerGo)

	// Default module: app
	generateModuleInPath("app", name+"/modules")

	fmt.Println("✅ Project created at ./" + name)
}

func generateModule(name string) {
	fmt.Println("📦 Generating module:", name)
	generateModuleInPath(name, "modules")
	fmt.Println("✅ Module", name, "generated.")
}

func generateModuleInPath(name, basePath string) {
	modulePath := basePath + "/" + name
	os.MkdirAll(modulePath, 0755)

	controller := `package ` + name + `

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

func (c *Controller) RegisterRoutes(r chi.Router) {
	r.Get("/` + name + `", c.handleGet)
}

func (c *Controller) handleGet(w http.ResponseWriter, r *http.Request) {
	items := c.service.FindAll()
	for _, item := range items {
		fmt.Fprintln(w, item)
	}
}
`
	writeFile(modulePath+"/controller.go", controller)

	service := `package ` + name + `

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) FindAll() []string {
	return []string{"item1", "item2"}
}
`
	writeFile(modulePath+"/service.go", service)

	module := `package ` + name + `

import "PROJECT_NAME/core"

func RegisterDependencies() {
	core.Container.Provide(NewService)
	core.Container.Provide(NewController)
}
`
	module = replace(module, "PROJECT_NAME", getProjectNameFromPath(basePath))
	writeFile(modulePath+"/module.go", module)
}

func writeFile(path, content string) {
	os.WriteFile(path, []byte(content), 0644)
}

func replace(s, old, new string) string {
	return string([]byte(s))
	// You can enhance this to safely replace identifiers
}

func getProjectNameFromPath(basePath string) string {
	parts := []byte(basePath)
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == '/' || parts[i] == '\\' {
			return string(parts[i+1:])
		}
	}
	return string(parts)
}
