package scheduledigreel

import (
	"log"

	"github.com/pocketbase/pocketbase"
)

func ExecuteScheduledIGReels(app *pocketbase.PocketBase) error {
	scheduledIGReels, err := GetAllScheduledIGReels(app)
	if err != nil {
		log.Println("Couldn't find scheduled IG reels, returning", err)
		return err
	}

	return nil

	// for _, scheduledIGReel := range scheduledIGReels {
	// 	igAccountID, ok := scheduledIGReel["igAccountId"].(string)
	// 	if !ok {
	// 		log.Println("Couldn't get ig account ID from scheduled IG reel, skipping")
	// 		continue
	// 	}

	// 	igAccount, err := dao.FindRecordByID(shared.CollectionIGAccounts, igAccountID)
	// 	if err != nil {
	// 		log.Println("Couldn't find IG account, skipping", err)
	// 		continue
	// 	}

	// 	igSessionID, ok := igAccount["igSessionId"].(string)
	// 	if !ok {
	// 		log.Println("Couldn't get IG session ID from IG account, skipping")
	// 		continue
	// 	}

	// 	igReel, ok := scheduledIGReel["igReel"].(map[string]interface{})
	// 	if !ok {
	// 		log.Println("Couldn't get IG reel from scheduled IG reel, skipping")
	// 		continue
	// 	}

	// 	if err := igservice.UploadIGReel(igservice.UploadIGReelArgs{
	// 		SessionID: igSessionID,
	// 		Reel:      igReel,
	// 	}); err != nil {
	// 		log.Println("Couldn't upload IG reel, skipping", err)
	// 		continue
	// 	}

	// 	log.Println("Successfully uploaded IG reel")
	// }

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
