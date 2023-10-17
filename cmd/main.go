package main

import (
	"test/config"
	"test/internal/app"
)

func main() {
	// load configs
	cfg := config.MustLoad()

	app.Run(cfg.HTTPServer.Address)
}
