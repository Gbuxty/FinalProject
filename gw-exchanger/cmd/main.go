package main

import (
	"gw-exchanger/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Application failed to run: %v", err)
	}
}
