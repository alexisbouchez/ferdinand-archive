package main

import (
	"ferdinand/config"

	"github.com/caesar-rocks/core"
	"github.com/charmbracelet/log"
)

func main() {
	// Create a new database connection
	db, err := config.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Validate environment variables
	core.ValidateEnvironmentVariables[config.EnvironmentVariables]()

	// Run the Caesar web application
	app := config.ProvideApp(db)
	app.Run()
}
