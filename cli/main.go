package main

import (
	"fmt"
	"os"
	"strings"
)

func pascalCase(s string) string {
	return strings.Title(s)
}

func generateModule(name string) {
	moduleName := strings.ToLower(name)
	pascal := pascalCase(name)
	dir := fmt.Sprintf("modules/%s", moduleName)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Println("❌ Failed to create module directory:", err)
		return
	}

	write := func(filename, content string) {
		err := os.WriteFile(fmt.Sprintf("%s/%s.go", dir, filename), []byte(content), 0644)
		if err != nil {
			fmt.Println("❌ Failed to write", filename, ":", err)
		}
	}

	write("service", fmt.Sprintf(`package %s

type %sService struct {}

func New%sService() *%sService {
	return &%sService{}
}
`, moduleName, pascal, pascal, pascal, pascal))

	write("controller", fmt.Sprintf(`package %s

import (
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
)

type %sController struct {
	service *%sService
}

func New%sController(service *%sService) *%sController {
	return &%sController{service}
}

func (c *%sController) RegisterRoutes(r chi.Router) {
	r.Route("/%s", func(r chi.Router) {
		r.Get("/", c.GetAll)
	})
}

func (c *%sController) GetAll(w http.ResponseWriter, r *http.Request) {
	data := []string{"example"}
	json.NewEncoder(w).Encode(data)
}
`, moduleName, pascal, pascal, pascal, pascal, pascal, pascal, pascal, moduleName, pascal))

	write("module", fmt.Sprintf(`package %s

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/dig"
)

type Module struct {
	Controller *%sController
}

func NewModule(container *dig.Container) *Module {
	container.Provide(New%sService)
	container.Provide(New%sController)

	var controller *%sController
	_ = container.Invoke(func(c *%sController) {
		controller = c
	})

	return &Module{Controller: controller}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	m.Controller.RegisterRoutes(r)
}
`, moduleName, pascal, pascal, pascal, pascal, pascal))

	fmt.Println("✅ Module", name, "created at", dir)
}

func main() {
	if len(os.Args) < 4 || os.Args[1] != "g" || os.Args[2] != "module" {
		fmt.Println("Usage: go run cli/main.go g module <name>")
		return
	}
	generateModule(os.Args[3])
}
