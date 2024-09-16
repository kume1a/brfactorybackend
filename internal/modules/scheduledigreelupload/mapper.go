package scheduledigreelupload

import (
	"brfactorybackend/internal/shared"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

func ScheduledIGReelUploadRecordToModel(r *models.Record) ScheduledIGReelUpload {
	return ScheduledIGReelUpload{
		ID:              r.Id,
		Created:         r.Created,
		Updated:         r.Updated,
		Success:         r.GetBool(shared.ScheduledIGReelUpload_Success),
		Index:           r.GetInt(shared.ScheduledIGReelUpload_Index),
		Title:           r.GetString(shared.ScheduledIGReelUpload_Title),
		Caption:         r.GetString(shared.ScheduledIGReelUpload_Caption),
		IGMediaID:       r.GetString(shared.ScheduledIGReelUpload_IGMediaID),
		IGAccount:       r.GetString(shared.ScheduledIGReelUpload_IGAccount),
		ScheduledIGReel: r.GetString(shared.ScheduledIGReelUpload_ScheduledIGReel),
	}
}

func ScheduledIGReelUploadModelToRecord(app *pocketbase.PocketBase, m ScheduledIGReelUpload) (*models.Record, error) {
	collection, err := app.Dao().FindCollectionByNameOrId(shared.CollectionScheduledIGReelUploads)
	if err != nil {
		return nil, err
	}

	record := models.NewRecord(collection)

	record.Set(shared.ScheduledIGReelUpload_Success, m.Success)
	record.Set(shared.ScheduledIGReelUpload_Index, m.Index)
	record.Set(shared.ScheduledIGReelUpload_Title, m.Title)
	record.Set(shared.ScheduledIGReelUpload_Caption, m.Caption)
	record.Set(shared.ScheduledIGReelUpload_IGMediaID, m.IGMediaID)
	record.Set(shared.ScheduledIGReelUpload_IGAccount, m.IGAccount)
	record.Set(shared.ScheduledIGReelUpload_ScheduledIGReel, m.ScheduledIGReel)

	return record, nil
}
