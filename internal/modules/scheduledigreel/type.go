package scheduledigreel

import (
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

func (r *ScheduledIGReel) VideoURL() string {
	return shared.ConstructPBFilePath(shared.CollectionScheduledIGReels, r.ID, r.Video)
}

func (r *ScheduledIGReel) ThumbnailURL() string {
	return shared.ConstructPBFilePath(shared.CollectionScheduledIGReels, r.ID, r.Thumbnail)
}
