package scheduledigreel

import (
	"brfactorybackend/internal/shared"

	"github.com/pocketbase/pocketbase/models"
)

func ScheduledIGReelRecordToModel(record *models.Record) ScheduledIGReel {
	return ScheduledIGReel{
		ID:              record.Id,
		Created:         record.Created,
		Updated:         record.Updated,
		StartAt:         record.GetDateTime(shared.ScheduledIGReel_StartAt),
		NextUploadTime:  record.GetDateTime(shared.ScheduledIGReel_NextUploadTime),
		IntervalInSecs:  record.GetInt(shared.ScheduledIGReel_IntervalInSeconds),
		Title:           record.GetString(shared.ScheduledIGReel_Title),
		Caption:         record.GetString(shared.ScheduledIGReel_Caption),
		ThumbnailFileID: record.GetString(shared.ScheduledIGReel_ThumbnailFileID),
		VideoFileID:     record.GetString(shared.ScheduledIGReel_VideoFileID),
		IGAccount:       record.GetString(shared.ScheduledIGReel_IGAccount),
	}
}

func ScheduledIGReelSetRecordFields(record *models.Record, e ScheduledIGReel) *models.Record {
	record.Set(shared.ScheduledIGReel_StartAt, e.StartAt)
	record.Set(shared.ScheduledIGReel_NextUploadTime, e.NextUploadTime)
	record.Set(shared.ScheduledIGReel_IntervalInSeconds, e.IntervalInSecs)
	record.Set(shared.ScheduledIGReel_Title, e.Title)
	record.Set(shared.ScheduledIGReel_Caption, e.Caption)
	record.Set(shared.ScheduledIGReel_ThumbnailFileID, e.ThumbnailFileID)
	record.Set(shared.ScheduledIGReel_VideoFileID, e.VideoFileID)
	record.Set(shared.ScheduledIGReel_IGAccount, e.IGAccount)

	return record
}
