package scheduledigreelupload

import (
	"brfactorybackend/internal/shared"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

func GetLatestSuccessScheduledIGReelUpload(
	app *pocketbase.PocketBase,
	scheduledIGReelID string,
) (*ScheduledIGReelUpload, error) {
	dao := app.Dao()

	query := dao.RecordQuery(shared.CollectionScheduledIGReelUploads).
		Where(dbx.HashExp{shared.ScheduledIGReelUpload_Success: true}).
		AndWhere(dbx.NewExp(
			shared.ScheduledIGReelUpload_ScheduledIGReel+" = {:0}",
			dbx.Params{"0": scheduledIGReelID}),
		).
		OrderBy("created DESC").
		Limit(1)

	records := []*models.Record{}
	if err := query.All(&records); err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, nil
	}

	mapped := ScheduledIGReelUploadRecordToModel(records[0])

	return &mapped, nil
}

func CreateScheduledIGReelUpload(
	app *pocketbase.PocketBase,
	scheduledIGReelUpload ScheduledIGReelUpload,
) (*ScheduledIGReelUpload, error) {
	record, err := ScheduledIGReelUploadModelToRecord(app, scheduledIGReelUpload)
	if err != nil {
		return nil, err
	}

	if err := app.Dao().SaveRecord(record); err != nil {
		return nil, err
	}

	scheduledIGReelUpload.ID = record.Id

	return &scheduledIGReelUpload, nil
}
