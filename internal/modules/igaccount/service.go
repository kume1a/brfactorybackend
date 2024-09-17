package igaccount

import (
	"brfactorybackend/internal/modules/igservice"
	"brfactorybackend/internal/shared"

	"github.com/pocketbase/pocketbase"
)

func GetAllIGAccounts(app *pocketbase.PocketBase) ([]IGAccount, error) {
	dao := app.Dao()

	records, err := dao.FindRecordsByExpr(shared.CollectionIGAccounts, nil)
	if err != nil {
		return nil, err
	}

	var mapped []IGAccount
	for _, record := range records {
		mapped = append(mapped, IGAccountRecordToModel(record))
	}

	return mapped, nil
}

func GetIGAccountByID(
	app *pocketbase.PocketBase,
	id string,
) (IGAccount, error) {
	dao := app.Dao()

	record, err := dao.FindRecordById(shared.CollectionIGAccounts, id)
	if err != nil {
		return IGAccount{}, err
	}

	return IGAccountRecordToModel(record), nil
}

func EnsureIGAccountIGSessionID(
	app *pocketbase.PocketBase,
	id string,
) (string, error) {
	dao := app.Dao()

	igAccountRecord, err := dao.FindRecordById(shared.CollectionIGAccounts, id)
	if err != nil {
		return "", err
	}

	igAccount := IGAccountRecordToModel(igAccountRecord)
	if igAccount.IGSessionID != "" {
		return igAccount.IGSessionID, nil
	}

	igSessionID, err := igservice.GetIGSessionID(
		app,
		igservice.GetIGSessionTokenArgs{
			IGUsername: igAccount.Email,
			IOPassword: igAccount.Password,
		},
	)
	if err != nil {
		return "", err
	}

	igAccountRecord.Set(shared.IGAccount_IGSessionID, igSessionID)

	if err := dao.SaveRecord(igAccountRecord); err != nil {
		return "", err
	}

	return igSessionID, nil
}
