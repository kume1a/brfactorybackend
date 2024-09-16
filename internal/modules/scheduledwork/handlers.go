package scheduledwork

import (
	"brfactorybackend/internal/modules/igaccount"
	"brfactorybackend/internal/modules/igservice"
	"brfactorybackend/internal/modules/scheduledigreel"
	"brfactorybackend/internal/modules/scheduledigreelupload"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase"
)

func ExecuteScheduledIGReels(app *pocketbase.PocketBase) error {
	scheduledIGReels, err := scheduledigreel.GetAllScheduledIGReels(app)
	if err != nil {
		log.Println("Couldn't find scheduled IG reels, returning", err)
		return err
	}

	for _, scheduledIGReel := range scheduledIGReels {
		igAccount, err := igaccount.GetIGAccountByID(app, scheduledIGReel.IGAccount)
		if err != nil {
			log.Println("Couldn't find IG account, skipping", err)
			return err
		}

		igSessionID, err := igaccount.EnsureIGAccountIGSessionID(app, igAccount.ID)
		if err != nil {
			return err
		}

		latestIGReelUpload, err := scheduledigreelupload.GetLatestSuccessScheduledIGReelUpload(
			app,
			scheduledIGReel.ID,
		)
		if err != nil {
			return err
		}

		if latestIGReelUpload != nil {
			log.Println("latestIGReelUpload is not nil, checking if it's time to upload")

			now := time.Now()
			diffSinceLastUpload := now.Sub(latestIGReelUpload.Created.Time())

			log.Println("diffSinceLastUpload.Seconds() = ", diffSinceLastUpload.Seconds())

			if diffSinceLastUpload.Seconds() < float64(scheduledIGReel.IntervalInSecs) {
				log.Println("It's not time to upload yet, skipping")
				continue
			}
		}

		var nextIndex int
		if latestIGReelUpload == nil {
			nextIndex = 0
		} else {
			nextIndex = latestIGReelUpload.Index + 1
		}

		log.Println("Preparing to upload IGTV video, nextIndex = ", nextIndex)

		videoFileURL, err := scheduledIGReel.VideoFileURL()
		if err != nil {
			log.Println("Couldn't get video file URL, skipping", err)
			return err
		}

		thumbnailFileURL, err := scheduledIGReel.ThumbnailFileURL()
		if err != nil {
			log.Println("Couldn't get thumbnail file URL, skipping", err)
			return err
		}

		title := strings.ReplaceAll(scheduledIGReel.Title, "{indexPlusOne}", strconv.Itoa(nextIndex+1))
		title = strings.ReplaceAll(title, "{index}", strconv.Itoa(nextIndex))

		log.Println("Uploading IGTV video, Title = " + title + ", Caption = " + scheduledIGReel.Caption + ", SessionID = " + igSessionID + ", VideoURL = " + videoFileURL + ", ThumbnailURL = " + thumbnailFileURL)

		igMediaID, uploadIGTVVideoErr := igservice.UploadIGTVVideo(igservice.UploadIGTVVideoArgs{
			Title:        title,
			Caption:      scheduledIGReel.Caption,
			SessionID:    igSessionID,
			VideoURL:     videoFileURL,
			ThumbnailURL: thumbnailFileURL,
		})
		if uploadIGTVVideoErr != nil {
			log.Println("Couldn't upload IG reel, ", uploadIGTVVideoErr)
		}

		if _, err := scheduledigreelupload.CreateScheduledIGReelUpload(
			app,
			scheduledigreelupload.ScheduledIGReelUpload{
				Success:         uploadIGTVVideoErr == nil,
				Index:           nextIndex,
				Title:           title,
				Caption:         scheduledIGReel.Caption,
				IGMediaID:       igMediaID,
				IGAccount:       scheduledIGReel.IGAccount,
				ScheduledIGReel: scheduledIGReel.ID,
			},
		); err != nil {
			log.Println("Couldn't create scheduled IG reel upload, skipping", err)
			return err
		}
	}

	return nil
}
