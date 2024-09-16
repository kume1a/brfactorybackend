package scheduledigreel

import (
	"brfactorybackend/internal/shared"

	"github.com/pocketbase/pocketbase/models"
)

func ScheduledIGReelRecordToModel(record *models.Record) ScheduledIGReel {
	id := record.Id
	startAt := record.GetDateTime(shared.ScheduledIGReels_StartAt)
	intervalInSecs := record.GetInt(shared.ScheduledIGReels_IntervalInSeconds)
	title := record.GetString(shared.ScheduledIGReels_Title)
	caption := record.GetString(shared.ScheduledIGReels_Caption)
	thumbnail := record.GetString(shared.ScheduledIGReels_Thumbnail)
	video := record.GetString(shared.ScheduledIGReels_Video)
	igAccount := record.GetString(shared.ScheduledIGReels_IGAccount)

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
