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

	ScheduledIGReel_StartAt           = "startAt"
	ScheduledIGReel_IntervalInSeconds = "intervalInSeconds"
	ScheduledIGReel_Title             = "title"
	ScheduledIGReel_Caption           = "caption"
	ScheduledIGReel_ThumbnailFileID   = "thumbnailFileId"
	ScheduledIGReel_VideoFileID       = "videoFileId"
	ScheduledIGReel_IGAccount         = "igAccount"

	CollectionScheduledIGReelUploads = "scheduledIGReelUploads"

	ScheduledIGReelUpload_Success         = "success"
	ScheduledIGReelUpload_Index           = "index"
	ScheduledIGReelUpload_Title           = "title"
	ScheduledIGReelUpload_Caption         = "caption"
	ScheduledIGReelUpload_IGMediaID       = "igMediaId"
	ScheduledIGReelUpload_IGAccount       = "igAccount"
	ScheduledIGReelUpload_ScheduledIGReel = "scheduledIGReel"
)
