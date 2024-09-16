package scheduledigreel

import (
	"brfactorybackend/internal/shared"

	"github.com/pocketbase/pocketbase/models"
)

func ScheduledIGReelRecordToModel(record *models.Record) ScheduledIGReel {
	id := record.Id
	startAt := record.GetDateTime(shared.ScheduledIGReel_StartAt)
	intervalInSecs := record.GetInt(shared.ScheduledIGReel_IntervalInSeconds)
	title := record.GetString(shared.ScheduledIGReel_Title)
	caption := record.GetString(shared.ScheduledIGReel_Caption)
	thumbnail := record.GetString(shared.ScheduledIGReel_Thumbnail)
	video := record.GetString(shared.ScheduledIGReel_Video)
	igAccount := record.GetString(shared.ScheduledIGReel_IGAccount)

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
