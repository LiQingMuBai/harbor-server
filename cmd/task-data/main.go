package main

import (
	bootstraptask "cointrade/internal/bootstrap/task"
	"cointrade/utils"
	"log"
)

func main() {
	if err := utils.SetupServiceLogger("task/data"); err != nil {
		log.Fatal(err)
	}
	log.Printf("==================================================")
	log.Printf("START harbor-server task-data")
	log.Printf("service=task mode=data")
	log.Printf("==================================================")
	if err := bootstraptask.RunData(); err != nil {
		log.Fatal(err)
	}
}

