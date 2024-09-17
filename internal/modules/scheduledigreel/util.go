package scheduledigreel

import (
	"brfactorybackend/internal/config"
	"brfactorybackend/internal/shared"
	"strconv"
	"strings"
)

func (r *ScheduledIGReel) VideoFileURL() (string, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return "", err
	}

	path := shared.ConstructPBFilePath(
		shared.CollectionScheduledIGReels,
		r.ID,
		r.VideoFileID,
	)

	return env.FileURLPrefix + path, nil
}

func (r *ScheduledIGReel) ThumbnailFileURL() (string, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return "", err
	}

	path := shared.ConstructPBFilePath(
		shared.CollectionScheduledIGReels,
		r.ID,
		r.ThumbnailFileID,
	)

	return env.FileURLPrefix + path, nil
}

func (r *ScheduledIGReel) FormattedCaption(index int) string {
	return formatWithVars(r.Caption, index)
}

func (r *ScheduledIGReel) FormattedTitle(index int) string {
	return formatWithVars(r.Title, index)
}

func formatWithVars(value string, index int) string {
	res := strings.ReplaceAll(value, "{indexPlusOne}", strconv.Itoa(index+1))
	res = strings.ReplaceAll(res, "{index}", strconv.Itoa(index))

	return res
}
