package igaccount

import "github.com/pocketbase/pocketbase/tools/types"

type IGAccount struct {
	ID          string         `json:"id"`
	Created     types.DateTime `json:"created"`
	Updated     types.DateTime `json:"updated"`
	Username    string         `json:"username"`
	Email       string         `json:"email"`
	Password    string         `json:"password"`
	IGSessionID string         `json:"igSessionId"`
	User        string         `json:"user"`
}
