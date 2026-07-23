package main

import (
	"cointrade/internal/tools/mongoinit"
	"cointrade/utils"
	"log"
)

func main() {
	if err := utils.SetupServiceLogger("tools/mongo-init"); err != nil {
		log.Fatal(err)
	}
	mongoinit.Run()
}
