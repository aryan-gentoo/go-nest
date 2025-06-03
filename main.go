package main

import (
	"go-nest.com/m/core"
	"go-nest.com/m/modules/user"
)

func main() {
	container := core.Container

	// Create modules
	userModule := user.NewModule(container)

	// Initialize app
	app := core.NewApp(userModule)
	app.Run(":3000")
}
