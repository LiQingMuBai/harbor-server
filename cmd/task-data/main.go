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
	utils.ServiceStartupBanner("harbor-server task-data", "service", "task", "mode", "data")
	if err := bootstraptask.RunData(); err != nil {
		log.Fatal(err)
	}
}
