package igaccount

import (
	"brfactorybackend/internal/shared"

	"github.com/pocketbase/pocketbase/models"
)

func IGAccountRecordToModel(record *models.Record) IGAccount {
	id := record.Id
	username := record.GetString(shared.IGAccount_Username)
	email := record.GetString(shared.IGAccount_Email)
	password := record.GetString(shared.IGAccount_Password)
	igSessionID := record.GetString(shared.IGAccount_IGSessionID)
	user := record.GetString(shared.IGAccount_User)

	return IGAccount{
		ID:          id,
		Username:    username,
		Email:       email,
		Password:    password,
		IGSessionID: igSessionID,
		User:        user,
	}
}
