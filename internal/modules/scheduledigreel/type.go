package scheduledigreel

import (
	"github.com/pocketbase/pocketbase/tools/types"
)

type ScheduledIGReel struct {
	ID             string         `json:"id"`
	StartAt        types.DateTime `json:"startAt"`
	IntervalInSecs int            `json:"intervalInSeconds"`
	Title          string         `json:"title"`
	Caption        string         `json:"caption"`
	Thumbnail      string         `json:"thumbnail"`
	Video          string         `json:"video"`
	IGAccount      string         `json:"igAccount"`
}
