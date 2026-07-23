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
	utils.ServiceStartupBanner("harbor-server task-jobs", "service", "task", "mode", "task")
	if err := bootstraptask.RunJobs(); err != nil {
		log.Fatal(err)
	}
}
