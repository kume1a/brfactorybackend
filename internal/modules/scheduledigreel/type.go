package scheduledigreel

import (
	"github.com/pocketbase/pocketbase/tools/types"
)

type ScheduledIGReel struct {
	ID              string         `json:"id"`
	Created         types.DateTime `json:"created"`
	Updated         types.DateTime `json:"updated"`
	StartAt         types.DateTime `json:"startAt"`
	NextUploadTime  types.DateTime `json:"nextUploadTime"`
	IntervalInSecs  int            `json:"intervalInSeconds"`
	Title           string         `json:"title"`
	Caption         string         `json:"caption"`
	ThumbnailFileID string         `json:"thumbnailFileId"`
	VideoFileID     string         `json:"videoFileId"`
	IGAccount       string         `json:"igAccount"`
}
