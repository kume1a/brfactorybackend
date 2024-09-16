package scheduledwork

import (
	"brfactorybackend/internal/modules/igaccount"
	"brfactorybackend/internal/modules/igservice"
	"brfactorybackend/internal/modules/scheduledigreel"
	"log"

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
			continue
		}

		igSessionID, err := igaccount.EnsureIGAccountIGSessionID(app, igAccount.ID)
		if err != nil {
			return err
		}

		if err := igservice.UploadIGTVVideo(igservice.UploadIGTVVideoArgs{
			Title:        scheduledIGReel.Title,
			Caption:      scheduledIGReel.Caption,
			SessionID:    igSessionID,
			VideoURL:     "",
			ThumbnailURL: "",
		}); err != nil {
			log.Println("Couldn't upload IG reel, skipping", err)
			continue
		}
	}
	return nil

	// --------------------------------------------------------------------------------------------------------------

	// igSessionID, err := igservice.GetIGSessionID(igservice.GetIGSessionTokenArgs{
	// 	IGUsername: "testbrainrot000@gmail.com",
	// 	IOPassword: "testbrainrot12345",
	// })
	// if err != nil {
	// 	log.Println("Couldn't get IG session ID, returning")
	// 	return err
	// }

	// igSessionID := "67954392189%3AHwXgIx22CTu3OL%3A29%3AAYesRsUA4Hkd9Aayf4ZBqgcWPRShs185SfHv2Dk1Gw"

	// if err := igservice.UploadIGTVVideo(igservice.UploadIGTVVideoArgs{
	// 	SessionID: igSessionID,
	// 	VideoURL:  "https://firebasestorage.googleapis.com/v0/b/literature-xii.appspot.com/o/Tough%20times%20never%20last%2C%20only%20tough%20people%20last%20ahuefieifehifeh.mp4?alt=media&token=33c70791-12d5-40f3-bc39-f0e6fba4ed05",
	// 	Title:     "Test video",
	// 	Caption:   "Test caption",
	// 	// ThumbnailURL: "https://www.youtube.com/watch?v=3v1n6b7j5Zk",
	// }); err != nil {
	// 	log.Println("Couldn't upload IGTV video, returning")
	// 	return err
	// }

	// log.Println("IG session ID: ", igSessionID)

}
