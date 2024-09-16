package scheduledigreel

import "github.com/pocketbase/pocketbase/models"

func mapRecordToStruct(record *models.Record) ScheduledIGReel {
	id := record.Id
	startAt := record.GetDateTime("startAt")
	intervalInSecs := record.GetInt("intervalInSeconds")
	title := record.GetString("title")
	caption := record.GetString("caption")
	thumbnail := record.GetString("thumbnail")
	video := record.GetString("video")
	igAccount := record.GetString("igAccount")

	return ScheduledIGReel{
		ID:             id,
		StartAt:        startAt,
		IntervalInSecs: intervalInSecs,
		Title:          title,
		Caption:        caption,
		Thumbnail:      thumbnail,
		Video:          video,
		IGAccount:      igAccount,
	}
}
