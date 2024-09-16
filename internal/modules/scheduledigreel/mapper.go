package scheduledigreel

import (
	"brfactorybackend/internal/shared"

	"github.com/pocketbase/pocketbase/models"
)

func ScheduledIGReelRecordToModel(record *models.Record) ScheduledIGReel {
	return ScheduledIGReel{
		ID:             record.Id,
		Created:        record.Created,
		Updated:        record.Updated,
		StartAt:        record.GetDateTime(shared.ScheduledIGReel_StartAt),
		IntervalInSecs: record.GetInt(shared.ScheduledIGReel_IntervalInSeconds),
		Title:          record.GetString(shared.ScheduledIGReel_Title),
		Caption:        record.GetString(shared.ScheduledIGReel_Caption),
		Thumbnail:      record.GetString(shared.ScheduledIGReel_Thumbnail),
		Video:          record.GetString(shared.ScheduledIGReel_Video),
		IGAccount:      record.GetString(shared.ScheduledIGReel_IGAccount),
	}
}
