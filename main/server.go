package main

import (
	"github.com/GoPex/unleash"
)

func main() {
	// Create a new Unleash application
	application := unleash.New()

	// Parse the configuration
	config, err := unleash.ParseConfiguration()
	if err != nil {
		panic("Not able to parse the configuration ! Cause: " + err.Error())
	}

	// Initialize the application
	if err = application.Initialize(&config); err != nil {
		panic("Not able to initialize the application ! Cause: " + err.Error())
	}

	// Listen and serve on port defined by environment variable UNLEASH_PORT
	application.Engine.Run(":" + unleash.Config.Port)
}
