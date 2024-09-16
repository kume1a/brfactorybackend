package scheduledigreel

import (
	"brfactorybackend/internal/shared"

	"github.com/pocketbase/pocketbase"
)

func GetAllScheduledIGReels(app *pocketbase.PocketBase) ([]ScheduledIGReel, error) {
	dao := app.Dao()

	records, err := dao.FindRecordsByExpr(shared.CollectionScheduledIGReels, nil)
	if err != nil {
		return nil, err
	}

	var mapped []ScheduledIGReel
	for _, record := range records {
		mapped = append(mapped, ScheduledIGReelRecordToModel(record))
	}

	return mapped, nil
}
