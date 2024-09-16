package igaccount

type IGAccount struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IGSessionID string `json:"igSessionId"`
	User        string `json:"user"`
}
