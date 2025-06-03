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
	mainCode := `package main

import "fmt"

func main() {
	fmt.Println("🚀 Welcome to ` + name + ` - powered by Go Nest")
}
`
	writeFile(name+"/main.go", mainCode)

	// go.mod
	cmd := exec.Command("go", "mod", "init", name)
	cmd.Dir = name
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	fmt.Println("✅ Project created at ./" + name)
}

func generateModule(name string) {
	fmt.Println("📦 Generating module:", name)
	os.MkdirAll("modules/"+name, 0755)

	controller := `package ` + name + `

import "fmt"

func Get() {
	fmt.Println("Hello from ` + name + ` controller")
}
`
	writeFile("modules/"+name+"/controller.go", controller)

	service := `package ` + name + `

func FindAll() []string {
	return []string{"item1", "item2"}
}
`
	writeFile("modules/"+name+"/service.go", service)

	module := `package ` + name + `

// Register your module here
`
	writeFile("modules/"+name+"/module.go", module)

	fmt.Println("✅ Module", name, "generated.")
}

func writeFile(path, content string) {
	os.WriteFile(path, []byte(content), 0644)
}
