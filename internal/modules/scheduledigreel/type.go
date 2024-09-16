package scheduledigreel

import (
	"brfactorybackend/internal/config"
	"brfactorybackend/internal/shared"

	"github.com/pocketbase/pocketbase/tools/types"
)

type ScheduledIGReel struct {
	ID             string         `json:"id"`
	Created        types.DateTime `json:"created"`
	Updated        types.DateTime `json:"updated"`
	StartAt        types.DateTime `json:"startAt"`
	IntervalInSecs int            `json:"intervalInSeconds"`
	Title          string         `json:"title"`
	Caption        string         `json:"caption"`
	Thumbnail      string         `json:"thumbnail"`
	Video          string         `json:"video"`
	IGAccount      string         `json:"igAccount"`
}

func (r *ScheduledIGReel) VideoFileURL() (string, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return "", err
	}

	return env.FileURLPrefix + shared.ConstructPBFilePath(shared.CollectionScheduledIGReels, r.ID, r.Video), nil
}

func (r *ScheduledIGReel) ThumbnailFileURL() (string, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return "", err
	}

	return env.FileURLPrefix + shared.ConstructPBFilePath(shared.CollectionScheduledIGReels, r.ID, r.Thumbnail), nil
}
