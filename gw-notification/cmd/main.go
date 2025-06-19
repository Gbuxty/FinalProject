package main

import (
	"gw-notification/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("app error:%v",err)
	}
}
