package main

import (
	"gw-currency-wallet/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("app error:%v", err)
	}
}
