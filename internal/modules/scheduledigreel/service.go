package scheduledigreel

import (
	"brfactorybackend/internal/shared"

	"github.com/pocketbase/pocketbase"
)

func GetAllScheduledIGReels(app *pocketbase.PocketBase) ([]ScheduledIGReel, error) {
	dao := app.Dao()

	scheduledIGReels, err := dao.FindRecordsByExpr(shared.CollectionScheduledIGReels, nil)
	if err != nil {
		return nil, err
	}

	var records []ScheduledIGReel
	for _, record := range scheduledIGReels {
		scheduledIGReel := mapRecordToStruct(record)
		records = append(records, scheduledIGReel)
	}

	return records, nil
}
