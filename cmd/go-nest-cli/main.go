package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	case "refresh":
		refreshMainFile(".")
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Go Nest CLI")
	fmt.Println("Usage:")
	fmt.Println("  go-nest-cli new <project-name>       Create a new project")
	fmt.Println("  go-nest-cli g module <module-name>   Generate a module")
	fmt.Println("  go-nest-cli refresh                  Regenerate main.go with all modules")
}

func createNewProject(name string) {
	fmt.Println("🚀 Creating new Go Nest project:", name)

	os.Mkdir(name, 0755)
	os.MkdirAll(filepath.Join(name, "modules", "app"), 0755)
	os.MkdirAll(filepath.Join(name, "core"), 0755)

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
}`
	writeFile(filepath.Join(name, "core", "app.go"), appGo)

	// container.go
	containerGo := `package core

import "go.uber.org/dig"

var Container = dig.New()`
	writeFile(filepath.Join(name, "core", "container.go"), containerGo)

	generateModuleInPath("app", filepath.Join(name, "modules"))
	refreshMainFile(name)

	fmt.Println("✅ Project created at ./" + name)
}

func generateModule(name string) {
	fmt.Println("📦 Generating module:", name)
	generateModuleInPath(name, "modules")
	refreshMainFile(".")
	fmt.Println("✅ Module", name, "generated.")
}

func generateModuleInPath(name, basePath string) {
	modulePath := filepath.Join(basePath, name)
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
}`
	writeFile(filepath.Join(modulePath, "controller.go"), controller)

	service := `package ` + name + `

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) FindAll() []string {
	return []string{"item1", "item2"}
}`
	writeFile(filepath.Join(modulePath, "service.go"), service)

	module := `package ` + name + `

import "PROJECT_NAME/core"

func RegisterDependencies() {
	core.Container.Provide(NewService)
	core.Container.Provide(NewController)
}`
	module = strings.ReplaceAll(module, "PROJECT_NAME", getProjectNameFromPath(basePath))
	writeFile(filepath.Join(modulePath, "module.go"), module)
}

func refreshMainFile(projectDir string) {
	modulesDir := filepath.Join(projectDir, "modules")
	entries, err := ioutil.ReadDir(modulesDir)
	if err != nil {
		fmt.Println("❌ Failed to read modules directory:", err)
		return
	}

	imports := []string{}
	invokes := []string{}

	for _, entry := range entries {
		if entry.IsDir() {
			mod := entry.Name()
			imports = append(imports, fmt.Sprintf("\t\"%s/modules/%s\"", projectDir, mod))
			invokes = append(invokes, fmt.Sprintf("\t%s.RegisterDependencies()", mod))
			invokes = append(invokes, fmt.Sprintf("\tcore.Container.Invoke(func(c *%s.Controller) { modules = append(modules, c) })", mod))
		}
	}

	contents := fmt.Sprintf(`package main

import (
	"%s/core"
%s
)

func main() {
%s

	var modules []core.Module
%s
	app := core.NewApp(modules...)
	app.Run(":8080")
}`, projectDir, strings.Join(imports, "\n"), strings.Join(invokes[:len(invokes)/2], "\n"), strings.Join(invokes[len(invokes)/2:], "\n"))

	writeFile(filepath.Join(projectDir, "main.go"), contents)
}

func writeFile(path, content string) {
	os.WriteFile(path, []byte(content), 0644)
}

func getProjectNameFromPath(basePath string) string {
	clean := filepath.Clean(basePath)
	return filepath.Base(clean)
}
