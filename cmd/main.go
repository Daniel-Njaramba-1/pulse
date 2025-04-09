package main

import (
	"github.com/Daniel-Njaramba-1/pulse/internal/app"
	"github.com/Daniel-Njaramba-1/pulse/internal/config"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/logging"
)

func main() {
	config.LoadEnv()

	logging.InitLogging()
	defer logging.CloseLogging()

	app, err := app.NewApp()
	if err != nil {
		logging.LogError("Failed to serve app: %v", err)
		return
	}
	defer app.Close()

	if err := app.Start(); err != nil {
		logging.LogError("Failed to start server: %v", err)
	}

}
