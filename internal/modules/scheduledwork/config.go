package scheduledwork

import (
	"github.com/pocketbase/pocketbase"
	"github.com/robfig/cron"
)

func SetupCronJobs(app *pocketbase.PocketBase) {
	scheduler := cron.New()

	// scheduler.AddFunc("0 */10 * * * *", func() {
	scheduler.AddFunc("*/30 * * * *", func() {
		// ExecuteScheduledIGReels(app)
	})

	scheduler.Start()
}
