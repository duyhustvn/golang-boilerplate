package main

import (
	"boilerplate/internal/handlers"
	"boilerplate/internal/server"
)

func main() {
	app := handlers.GetApp()
	server.Run(app)
}
