package scheduledwork

import (
	"github.com/pocketbase/pocketbase"
	"github.com/robfig/cron"
)

func SetupCronJobs(app *pocketbase.PocketBase) {
	scheduler := cron.New()

	// scheduler.AddFunc("0 */1 * * * *", func() {
	scheduler.AddFunc("*/10 * * * *", func() {
		ExecuteScheduledIGReels(app)
	})

	scheduler.Start()
}
