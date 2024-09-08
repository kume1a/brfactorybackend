package config

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func CreateCollections(app *pocketbase.PocketBase) error {

	collection := &models.Collection{
		Name:       "igAccounts",
		Type:       models.CollectionTypeBase,
		ListRule:   nil,
		ViewRule:   nil,
		CreateRule: nil,
		UpdateRule: nil,
		DeleteRule: nil,
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name:     "username",
				Type:     schema.FieldTypeText,
				Required: true,
				Options: &schema.TextOptions{
					Max: types.Pointer(255),
				},
			},
			&schema.SchemaField{
				Name:     "email",
				Type:     schema.FieldTypeText,
				Required: true,
				Options: &schema.TextOptions{
					Max: types.Pointer(255),
				},
			},
			&schema.SchemaField{
				Name:     "password",
				Type:     schema.FieldTypeText,
				Required: true,
				Options: &schema.TextOptions{
					Max: types.Pointer(255),
				},
			},
			// &schema.SchemaField{
			// 	Name:     "user",
			// 	Type:     schema.FieldTypeRelation,
			// 	Required: true,
			// 	Options: &schema.RelationOptions{
			// 		MaxSelect:     types.Pointer(1),
			// 		CollectionId:  "users",
			// 		CascadeDelete: true,
			// 	},
			// },
		),
	}

	if err := app.Dao().SaveCollection(collection); err != nil {
		return err
	}

	return nil
}
