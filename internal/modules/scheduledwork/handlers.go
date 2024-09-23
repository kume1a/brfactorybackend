package scheduledwork

import (
	"brfactorybackend/internal/modules/igaccount"
	"brfactorybackend/internal/modules/igservice"
	"brfactorybackend/internal/modules/scheduledigreel"
	"brfactorybackend/internal/modules/scheduledigreelupload"
	"log"
	"time"

	"github.com/pocketbase/pocketbase"
)

func ExecuteScheduledIGReels(app *pocketbase.PocketBase) error {
	logger := app.Logger().WithGroup("ExecuteScheduledIGReels")

	log.Println("Executing scheduled IG reels")

	scheduledIGReels, err := scheduledigreel.GetAllScheduledIGReels(app)
	if err != nil {
		logger.Error("Couldn't find scheduled IG reels, returning", "err", err)
		return err
	}

	log.Println("Found ", len(scheduledIGReels), " scheduled IG reels")

	for _, scheduledIGReel := range scheduledIGReels {
		log.Println("Processing scheduled IG reel: ", scheduledIGReel.ID)

		if scheduledIGReel.StartAt.Time().Before(time.Now()) {
			logger.Debug("Scheduled IG Reel " + scheduledIGReel.ID + " is not ready to be uploaded yet, skipping")
			continue
		}

		igAccount, err := igaccount.GetIGAccountByID(app, scheduledIGReel.IGAccount)
		if err != nil {
			logger.Error("Couldn't find IG account, skipping", "err", err)
			return err
		}

		igSessionID, err := igaccount.EnsureIGAccountIGSessionID(app, igAccount.ID)
		if err != nil {
			logger.Error("Couldn't get IG session ID, skipping", "err", err)
			return err
		}

		latestIGReelUpload, err := scheduledigreelupload.GetLatestSuccessScheduledIGReelUpload(
			app,
			scheduledIGReel.ID,
		)
		if err != nil {
			logger.Error("Couldn't get latest success scheduled IG reel upload, skipping", "err", err)
			return err
		}

		if latestIGReelUpload != nil {
			logger.Debug("scheduledIGReelID: " + scheduledIGReel.ID + ", latestIGReelUpload is not nil, checking if it's time to upload")

			now := time.Now()
			diffSinceLastUpload := now.Sub(latestIGReelUpload.Created.Time())

			logger.Debug("scheduledIGReelID: " + scheduledIGReel.ID + ", diffSinceLastUpload.Seconds() = " + diffSinceLastUpload.String())

			if diffSinceLastUpload.Seconds() < float64(scheduledIGReel.IntervalInSecs) {
				logger.Debug("scheduledIGReelID: " + scheduledIGReel.ID + ", It's not time to upload yet, skipping")
				continue
			}
		}

		var nextIndex int
		if latestIGReelUpload == nil {
			nextIndex = 0
		} else {
			nextIndex = latestIGReelUpload.Index + 1
		}

		videoFileURL, err := scheduledIGReel.VideoFileURL()
		if err != nil {
			logger.Error("Couldn't get video file URL, skipping", "err", err)
			return err
		}

		thumbnailFileURL, err := scheduledIGReel.ThumbnailFileURL()
		if err != nil {
			logger.Error("Couldn't get thumbnail file URL, skipping", "err", err)
			return err
		}

		title := scheduledIGReel.FormattedTitle(nextIndex)
		caption := scheduledIGReel.FormattedCaption(nextIndex)

		igMediaID, uploadIGTVVideoErr := igservice.UploadIGTVVideo(
			app,
			igservice.UploadIGTVVideoArgs{
				Title:        title,
				Caption:      caption,
				SessionID:    igSessionID,
				VideoURL:     videoFileURL,
				ThumbnailURL: thumbnailFileURL,
			},
		)
		if uploadIGTVVideoErr != nil {
			logger.Error("Couldn't upload IG reel, ", "err", uploadIGTVVideoErr)
		}

		if _, err := scheduledigreelupload.CreateScheduledIGReelUpload(
			app,
			scheduledigreelupload.ScheduledIGReelUpload{
				Success:         uploadIGTVVideoErr == nil,
				Index:           nextIndex,
				Title:           title,
				Caption:         caption,
				IGMediaID:       igMediaID,
				IGAccount:       scheduledIGReel.IGAccount,
				ScheduledIGReel: scheduledIGReel.ID,
			},
		); err != nil {
			logger.Error("Couldn't create scheduled IG reel upload, skipping", "err", err)
			return err
		}
	}

	return nil
}
