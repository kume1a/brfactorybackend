package scheduledigreel

import (
	"brfactorybackend/internal/config"
	"brfactorybackend/internal/shared"

	"github.com/pocketbase/pocketbase/tools/types"
)

type ScheduledIGReel struct {
	ID              string         `json:"id"`
	Created         types.DateTime `json:"created"`
	Updated         types.DateTime `json:"updated"`
	StartAt         types.DateTime `json:"startAt"`
	IntervalInSecs  int            `json:"intervalInSeconds"`
	Title           string         `json:"title"`
	Caption         string         `json:"caption"`
	ThumbnailFileID string         `json:"thumbnailFileId"`
	VideoFileID     string         `json:"videoFileId"`
	IGAccount       string         `json:"igAccount"`
}

func (r *ScheduledIGReel) VideoFileURL() (string, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return "", err
	}

	return env.FileURLPrefix + shared.ConstructPBFilePath(shared.CollectionScheduledIGReels, r.ID, r.VideoFileID), nil
}

func (r *ScheduledIGReel) ThumbnailFileURL() (string, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return "", err
	}

	return env.FileURLPrefix + shared.ConstructPBFilePath(shared.CollectionScheduledIGReels, r.ID, r.ThumbnailFileID), nil
}
