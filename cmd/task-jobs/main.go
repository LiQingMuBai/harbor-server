package main

import (
	bootstraptask "cointrade/internal/bootstrap/task"
	"cointrade/utils"
	"log"
)

func main() {
	if err := utils.SetupServiceLogger("task/jobs"); err != nil {
		log.Fatal(err)
	}
	log.Printf("==================================================")
	log.Printf("START harbor-server task")
	log.Printf("service=task mode=task")
	log.Printf("==================================================")
	if err := bootstraptask.RunJobs(); err != nil {
		log.Fatal(err)
	}
}

