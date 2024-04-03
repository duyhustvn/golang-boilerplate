package main

import (
	"boilerplate/internal/server"
	"log"
)

func main() {
	app := server.GetApp()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
