package cron

import (
	"log"
	"time"

	"github.com/bionicosmos/submgr/api"
	"github.com/bionicosmos/submgr/services"
	"github.com/go-co-op/gocron"
)

func Init() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(30).Minutes().Do(func() {
		if err := api.ResetAllNodes(); err != nil {
			log.Print(err)
		}
		if err := services.CheckUserBill(); err != nil {
			log.Print(err)
		}
	})
	scheduler.StartAsync()
}
