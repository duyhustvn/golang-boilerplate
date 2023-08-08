package main

import (
	"boilerplate/internal/server"
)

func main() {
	app := server.GetApp()
	app.Run()
}
