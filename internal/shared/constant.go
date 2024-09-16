package shared

const (
	CollectionUsers = "users"

	CollectionIGAccounts = "igAccounts"

	IGAccount_Username    = "username"
	IGAccount_Email       = "email"
	IGAccount_Password    = "password"
	IGAccount_IGSessionID = "igSessionId"
	IGAccount_User        = "user"

	CollectionScheduledIGReels = "scheduledIGReels"

	ScheduledIGReels_StartAt           = "startAt"
	ScheduledIGReels_IntervalInSeconds = "intervalInSeconds"
	ScheduledIGReels_Title             = "title"
	ScheduledIGReels_Caption           = "caption"
	ScheduledIGReels_Thumbnail         = "thumbnail"
	ScheduledIGReels_Video             = "video"
	ScheduledIGReels_IGAccount         = "igAccount"

	CollectionScheduledIGReelUploads = "scheduledIGReelUploads"

	ScheduledIGReelUploads_Success         = "success"
	ScheduledIGReelUploads_IGAccount       = "igAccount"
	ScheduledIGReelUploads_ScheduledIGReel = "scheduledIGReel"
)
