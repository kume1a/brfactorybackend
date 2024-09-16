package scheduledigreelupload

import "github.com/pocketbase/pocketbase/tools/types"

type ScheduledIGReelUpload struct {
	ID              string         `json:"id"`
	Created         types.DateTime `json:"created"`
	Updated         types.DateTime `json:"updated"`
	Success         bool           `json:"success"`
	Title           string         `json:"title"`
	Caption         string         `json:"caption"`
	IGAccount       string         `json:"igAccount"`
	ScheduledIGReel string         `json:"scheduledIGReel"`
}
