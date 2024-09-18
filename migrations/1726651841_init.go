package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		usersCollection, err := dao.FindCollectionByNameOrId("users")
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
					Name:     "igSessionId",
					Type:     schema.FieldTypeText,
					Required: false,
					Options: &schema.TextOptions{
						Max: types.Pointer(4095),
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

		if err := dao.SaveCollection(igAccountsCollection); err != nil {
			return err
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
					Name:     "thumbnailFileId",
					Type:     schema.FieldTypeFile,
					Required: true,
					Options: &schema.FileOptions{
						MaxSelect: 1,
						MaxSize:   10 * 1024 * 1024, // 10 MB
						MimeTypes: []string{"image/jpeg", "image/png"},
					},
				},
				&schema.SchemaField{
					Name:     "videoFileId",
					Type:     schema.FieldTypeFile,
					Required: true,
					Options: &schema.FileOptions{
						MaxSelect: 1,
						MaxSize:   100 * 1024 * 1024, // 100 MB
						MimeTypes: []string{"video/mp4"},
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

		if err := dao.SaveCollection(scheduledIGReelsCollection); err != nil {
			return err
		}

		scheduledIGReelUploadsCollection := &models.Collection{
			Name:       "scheduledIGReelUploads",
			Type:       models.CollectionTypeBase,
			ListRule:   types.Pointer(""),
			ViewRule:   types.Pointer(""),
			CreateRule: types.Pointer(""),
			UpdateRule: types.Pointer(""),
			DeleteRule: types.Pointer(""),
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     "success",
					Type:     schema.FieldTypeBool,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "index",
					Type:     schema.FieldTypeNumber,
					Required: true,
					Options: &schema.NumberOptions{
						Min:       types.Pointer(0.0),
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
					Name:     "igMediaId",
					Type:     schema.FieldTypeText,
					Required: false,
					Options: &schema.TextOptions{
						Max: types.Pointer(512),
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
				&schema.SchemaField{
					Name:     "scheduledIGReel",
					Type:     schema.FieldTypeRelation,
					Required: true,
					Options: &schema.RelationOptions{
						MaxSelect:     types.Pointer(1),
						CollectionId:  scheduledIGReelsCollection.Id,
						CascadeDelete: true,
					},
				},
			),
		}

		if err := dao.SaveCollection(scheduledIGReelUploadsCollection); err != nil {
			return err
		}

		return nil
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		igAccountsCollection, err := dao.FindCollectionByNameOrId("igAccounts")
		if err != nil {
			return err
		}

		scheduledIGReelsCollection, err := dao.FindCollectionByNameOrId("scheduledIGReels")
		if err != nil {
			return err
		}

		scheduledIGReelUploadsCollection, err := dao.FindCollectionByNameOrId("scheduledIGReelUploads")
		if err != nil {
			return err
		}

		if err := dao.DeleteCollection(igAccountsCollection); err != nil {
			return err
		}
		if err := dao.DeleteCollection(scheduledIGReelsCollection); err != nil {
			return err
		}
		if err := dao.DeleteCollection(scheduledIGReelUploadsCollection); err != nil {
			return err
		}

		return nil
	})
}
