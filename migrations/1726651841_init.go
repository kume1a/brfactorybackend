package migrations

import (
	"brfactorybackend/internal/shared"

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

		usersCollection, err := dao.FindCollectionByNameOrId(shared.CollectionUsers)
		if err != nil {
			return err
		}

		igAccountsCollection := &models.Collection{
			Name:       shared.CollectionIGAccounts,
			Type:       models.CollectionTypeBase,
			ListRule:   types.Pointer(""),
			ViewRule:   types.Pointer(""),
			CreateRule: types.Pointer(""),
			UpdateRule: types.Pointer(""),
			DeleteRule: types.Pointer(""),
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     shared.IGAccount_Username,
					Type:     schema.FieldTypeText,
					Required: true,
					Options: &schema.TextOptions{
						Max: types.Pointer(255),
					},
				},
				&schema.SchemaField{
					Name:     shared.IGAccount_Email,
					Type:     schema.FieldTypeText,
					Required: true,
					Options: &schema.TextOptions{
						Max: types.Pointer(255),
					},
				},
				&schema.SchemaField{
					Name:     shared.IGAccount_Password,
					Type:     schema.FieldTypeText,
					Required: true,
					Options: &schema.TextOptions{
						Max: types.Pointer(255),
					},
				},
				&schema.SchemaField{
					Name:     shared.IGAccount_IGSessionID,
					Type:     schema.FieldTypeText,
					Required: false,
					Options: &schema.TextOptions{
						Max: types.Pointer(4095),
					},
				},
				&schema.SchemaField{
					Name:     shared.IGAccount_User,
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
			Name:       shared.CollectionScheduledIGReels,
			Type:       models.CollectionTypeBase,
			ListRule:   types.Pointer(""),
			ViewRule:   types.Pointer(""),
			CreateRule: types.Pointer(""),
			UpdateRule: types.Pointer(""),
			DeleteRule: types.Pointer(""),
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     shared.ScheduledIGReel_StartAt,
					Type:     schema.FieldTypeDate,
					Required: true,
				},
				&schema.SchemaField{
					Name:     shared.ScheduledIGReel_IntervalInSeconds,
					Type:     schema.FieldTypeNumber,
					Required: true,
					Options: &schema.NumberOptions{
						Min:       types.Pointer(1.0),
						NoDecimal: true,
					},
				},
				&schema.SchemaField{
					Name:     shared.ScheduledIGReel_Title,
					Type:     schema.FieldTypeText,
					Required: true,
					Options: &schema.TextOptions{
						Max: types.Pointer(255),
					},
				},
				&schema.SchemaField{
					Name:     shared.ScheduledIGReel_Caption,
					Type:     schema.FieldTypeText,
					Required: true,
					Options: &schema.TextOptions{
						Max: types.Pointer(65535),
					},
				},
				&schema.SchemaField{
					Name:     shared.ScheduledIGReel_ThumbnailFileID,
					Type:     schema.FieldTypeFile,
					Required: true,
					Options: &schema.FileOptions{
						MaxSelect: 1,
						MaxSize:   10 * 1024 * 1024, // 10 MB
						MimeTypes: []string{"image/jpeg", "image/png"},
					},
				},
				&schema.SchemaField{
					Name:     shared.ScheduledIGReel_VideoFileID,
					Type:     schema.FieldTypeFile,
					Required: true,
					Options: &schema.FileOptions{
						MaxSelect: 1,
						MaxSize:   100 * 1024 * 1024, // 100 MB
						MimeTypes: []string{"video/mp4"},
					},
				},
				&schema.SchemaField{
					Name:     shared.ScheduledIGReel_IGAccount,
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
			Name:       shared.CollectionScheduledIGReelUploads,
			Type:       models.CollectionTypeBase,
			ListRule:   types.Pointer(""),
			ViewRule:   types.Pointer(""),
			CreateRule: types.Pointer(""),
			UpdateRule: types.Pointer(""),
			DeleteRule: types.Pointer(""),
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     shared.ScheduledIGReelUpload_Success,
					Type:     schema.FieldTypeBool,
					Required: true,
				},
				&schema.SchemaField{
					Name:     shared.ScheduledIGReelUpload_Index,
					Type:     schema.FieldTypeNumber,
					Required: true,
					Options: &schema.NumberOptions{
						Min:       types.Pointer(0.0),
						NoDecimal: true,
					},
				},
				&schema.SchemaField{
					Name:     shared.ScheduledIGReelUpload_Title,
					Type:     schema.FieldTypeText,
					Required: true,
					Options: &schema.TextOptions{
						Max: types.Pointer(255),
					},
				},
				&schema.SchemaField{
					Name:     shared.ScheduledIGReelUpload_Caption,
					Type:     schema.FieldTypeText,
					Required: true,
					Options: &schema.TextOptions{
						Max: types.Pointer(65535),
					},
				},
				&schema.SchemaField{
					Name:     shared.ScheduledIGReelUpload_IGMediaID,
					Type:     schema.FieldTypeText,
					Required: false,
					Options: &schema.TextOptions{
						Max: types.Pointer(512),
					},
				},
				&schema.SchemaField{
					Name:     shared.ScheduledIGReelUpload_IGAccount,
					Type:     schema.FieldTypeRelation,
					Required: true,
					Options: &schema.RelationOptions{
						MaxSelect:     types.Pointer(1),
						CollectionId:  igAccountsCollection.Id,
						CascadeDelete: true,
					},
				},
				&schema.SchemaField{
					Name:     shared.ScheduledIGReelUpload_ScheduledIGReel,
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

		igAccountsCollection, err := dao.FindCollectionByNameOrId(shared.CollectionIGAccounts)
		if err != nil {
			return err
		}

		scheduledIGReelsCollection, err := dao.FindCollectionByNameOrId(shared.CollectionScheduledIGReels)
		if err != nil {
			return err
		}

		scheduledIGReelUploadsCollection, err := dao.FindCollectionByNameOrId(shared.CollectionScheduledIGReelUploads)
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
