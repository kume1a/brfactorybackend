package config

import (
	"brfactorybackend/internal/modules/scheduledigreel"

	"github.com/pocketbase/pocketbase"
	"github.com/robfig/cron"
)

func SetupCronJobs(app *pocketbase.PocketBase) {
	scheduler := cron.New()

	// scheduler.AddFunc("0 */10 * * * *", func() {
	scheduler.AddFunc("*/5 * * * *", func() {
		scheduledigreel.ExecuteScheduledIGReels(app)
	})

	scheduler.Start()
}
