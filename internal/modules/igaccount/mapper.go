package igaccount

import (
	"brfactorybackend/internal/shared"

	"github.com/pocketbase/pocketbase/models"
)

func IGAccountRecordToModel(record *models.Record) IGAccount {
	return IGAccount{
		ID:          record.Id,
		Created:     record.Created,
		Updated:     record.Updated,
		Username:    record.GetString(shared.IGAccount_Username),
		Email:       record.GetString(shared.IGAccount_Email),
		Password:    record.GetString(shared.IGAccount_Password),
		IGSessionID: record.GetString(shared.IGAccount_IGSessionID),
		User:        record.GetString(shared.IGAccount_User),
	}
}
