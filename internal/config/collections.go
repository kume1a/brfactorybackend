package config

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func CreateCollections(app *pocketbase.PocketBase) error {
	usersCollection, err := app.Dao().FindCollectionByNameOrId("users")

	if err != nil {
		return err
	}

	igAccountsCollection := &models.Collection{
		Name:       "igAccounts",
		Type:       models.CollectionTypeBase,
		ListRule:   types.Pointer(""),
		ViewRule:   types.Pointer(""),
		CreateRule: types.Pointer(""),
		UpdateRule: types.Pointer(""),
		DeleteRule: types.Pointer(""),
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
			&schema.SchemaField{
				Name:     "user",
				Type:     schema.FieldTypeRelation,
				Required: true,
				Options: &schema.RelationOptions{
					MaxSelect:     types.Pointer(1),
					CollectionId:  usersCollection.Id,
					CascadeDelete: true,
				},
			},
		),
	}

	scheduledIGReelsCollection := &models.Collection{
		Name:       "scheduledIGReels",
		Type:       models.CollectionTypeBase,
		ListRule:   types.Pointer(""),
		ViewRule:   types.Pointer(""),
		CreateRule: types.Pointer(""),
		UpdateRule: types.Pointer(""),
		DeleteRule: types.Pointer(""),
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name:     "startAt",
				Type:     schema.FieldTypeDate,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "intervalInSeconds",
				Type:     schema.FieldTypeNumber,
				Required: true,
				Options: &schema.NumberOptions{
					Min:       types.Pointer(1.0),
					NoDecimal: true,
				},
			},
			&schema.SchemaField{
				Name:     "title",
				Type:     schema.FieldTypeText,
				Required: true,
				Options: &schema.TextOptions{
					Max: types.Pointer(255),
				},
			},
			&schema.SchemaField{
				Name:     "caption",
				Type:     schema.FieldTypeText,
				Required: true,
				Options: &schema.TextOptions{
					Max: types.Pointer(65535),
				},
			},
			&schema.SchemaField{
				Name:     "igAccount",
				Type:     schema.FieldTypeRelation,
				Required: true,
				Options: &schema.RelationOptions{
					MaxSelect:     types.Pointer(1),
					CollectionId:  igAccountsCollection.Id,
					CascadeDelete: true,
				},
			},
		),
	}

	if err := app.Dao().SaveCollection(igAccountsCollection); err != nil {
		return err
	}

	if err := app.Dao().SaveCollection(scheduledIGReelsCollection); err != nil {
		return err
	}

	return nil
}
