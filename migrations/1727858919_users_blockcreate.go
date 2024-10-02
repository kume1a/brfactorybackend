package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(
		func(db dbx.Builder) error {
			dao := daos.New(db)

			usersCollection, err := dao.FindCollectionByNameOrId("users")
			if err != nil {
				return err
			}

			scheduledIGReelsCollection, err := dao.FindCollectionByNameOrId("scheduledIGReels")
			if err != nil {
				return err
			}

			usersCollection.CreateRule = nil

			scheduledIGReelsCollection.Schema.AddField(&schema.SchemaField{
				Name:     "nextUploadTime",
				Type:     schema.FieldTypeDate,
				Required: false,
			})

			if err := dao.SaveCollection(usersCollection); err != nil {
				return err
			}

			if err := dao.SaveCollection(scheduledIGReelsCollection); err != nil {
				return err
			}

			return nil
		},
		func(db dbx.Builder) error {
			dao := daos.New(db)

			usersCollection, err := dao.FindCollectionByNameOrId("users")
			if err != nil {
				return err
			}

			scheduledIGReelsCollection, err := dao.FindCollectionByNameOrId("scheduledIGReels")
			if err != nil {
				return err
			}

			usersCollection.CreateRule = types.Pointer("")

			scheduledIGReelsCollection.Schema.RemoveField("nextUploadTime")

			if err := dao.SaveCollection(usersCollection); err != nil {
				return err
			}

			if err := dao.SaveCollection(scheduledIGReelsCollection); err != nil {
				return err
			}

			return nil
		},
	)
}
