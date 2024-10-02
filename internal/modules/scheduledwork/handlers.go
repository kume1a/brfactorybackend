package scheduledwork

import (
	"brfactorybackend/internal/modules/igaccount"
	"brfactorybackend/internal/modules/igservice"
	"brfactorybackend/internal/modules/scheduledigreel"
	"brfactorybackend/internal/modules/scheduledigreelupload"
	"log"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/tools/types"
)

func ExecuteScheduledIGReels(app *pocketbase.PocketBase) error {
	scheduledIGReels, err := scheduledigreel.GetAllScheduledIGReels(app)
	if err != nil {
		log.Println("Couldn't find scheduled IG reels, returning", "err", err)
		return err
	}

	for _, scheduledIGReel := range scheduledIGReels {
		if scheduledIGReel.NextUploadTime.IsZero() {
			intervalDuration := time.Duration(scheduledIGReel.IntervalInSecs) * time.Second
			initialUploadTime, err := types.ParseDateTime(scheduledIGReel.StartAt.Time().Add(intervalDuration))
			if err != nil {
				log.Println("Couldn't parse initial upload time, returning", "err", err)
				return err
			}

			scheduledIGReel.NextUploadTime = initialUploadTime
			if err := scheduledigreel.UpdateScheduledIGReel(app, scheduledIGReel.ID, scheduledIGReel); err != nil {
				log.Println("Couldn't update scheduled IG reel, returning", "err", err)
				return err
			}
		}

		if scheduledIGReel.NextUploadTime.Time().After(time.Now()) {
			log.Println("Scheduled IG Reel " + scheduledIGReel.ID + " is not ready to be uploaded yet, skipping")
			continue
		}

		igAccount, err := igaccount.GetIGAccountByID(app, scheduledIGReel.IGAccount)
		if err != nil {
			log.Println("Couldn't find IG account, returning", "err", err)
			return err
		}

		igSessionID, err := igaccount.EnsureIGAccountIGSessionID(app, igAccount.ID)
		if err != nil {
			log.Println("Couldn't get IG session ID, returning", "err", err)
			return err
		}

		latestIGReelUpload, err := scheduledigreelupload.GetLatestSuccessScheduledIGReelUpload(
			app,
			scheduledIGReel.ID,
		)
		if err != nil {
			log.Println("Couldn't get latest success scheduled IG reel upload, returning", "err", err)
			return err
		}

		var nextIndex int
		if latestIGReelUpload == nil {
			nextIndex = 0
		} else {
			nextIndex = latestIGReelUpload.Index + 1
		}

		videoFileURL, err := scheduledIGReel.VideoFileURL()
		if err != nil {
			log.Println("Couldn't get video file URL, returning", "err", err)
			return err
		}

		thumbnailFileURL, err := scheduledIGReel.ThumbnailFileURL()
		if err != nil {
			log.Println("Couldn't get thumbnail file URL, returning", "err", err)
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
			log.Println("Couldn't upload IG reel, ", "err", uploadIGTVVideoErr)
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
			log.Println("Couldn't create scheduled IG reel upload, returning", "err", err)
			return err
		}

		if uploadIGTVVideoErr == nil {
			intervalDuration := time.Duration(scheduledIGReel.IntervalInSecs) * time.Second
			nextUploadTime, err := types.ParseDateTime(scheduledIGReel.NextUploadTime.Time().Add(intervalDuration))
			if err != nil {
				log.Println("Couldn't parse next upload time, returning", "err", err)
				return err
			}

			scheduledIGReel.NextUploadTime = nextUploadTime

			if err := scheduledigreel.UpdateScheduledIGReel(app, scheduledIGReel.ID, scheduledIGReel); err != nil {
				log.Println("Couldn't update scheduled IG reel, returning", "err", err)
				return err
			}
		}
	}

	return nil
}
