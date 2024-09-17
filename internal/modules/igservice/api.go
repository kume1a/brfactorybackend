package igservice

import (
	"brfactorybackend/internal/config"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/pocketbase/pocketbase"
)

type UploadIGTVVideoArgs struct {
	SessionID    string `json:"sessionId"`
	Title        string `json:"title"`
	Caption      string `json:"caption"`
	ThumbnailURL string `json:"thumbnailURL"`
	VideoURL     string `json:"videoURL"`
}

type GetIGSessionTokenArgs struct {
	IGUsername string `json:"igUsername"`
	IOPassword string `json:"igPassword"`
}

type IGSessionTokenDTO struct {
	SessionID string `json:"sessionId"`
}

type UploadIGTVVideoDTO struct {
	MediaID string `json:"mediaId"`
}

func UploadIGTVVideo(
	app *pocketbase.PocketBase,
	args UploadIGTVVideoArgs,
) (string, error) {
	logger := app.Logger().WithGroup("UploadIGTVVideo")

	envVars, err := config.ParseEnv()
	if err != nil {
		logger.Error("Couldn't parse env vars, returning")
		return "", err
	}

	bodyBytes, err := json.Marshal(args)
	if err != nil {
		logger.Error("Couldn't marshal JSON body, returning")
		return "", err
	}

	client := resty.New()

	resp, err := client.R().
		SetBody(bodyBytes).
		SetResult(&UploadIGTVVideoDTO{}).
		SetHeader("X-Secret", envVars.IGServiceSecret).
		SetHeader("Content-Type", "application/json").
		Post(envVars.IGServiceURL + "/uploadIGTVVideo")

	if err != nil {
		logger.Error("Error sending request", "err", err)
		return "", err
	}

	statusCode := resp.StatusCode()
	if statusCode > 399 {
		logger.Error("invalid status code " + strconv.Itoa(statusCode) + ", response: " + resp.String())
		return "", errors.New("invalid status code " + strconv.Itoa(statusCode))
	}

	res := resp.Result().(*UploadIGTVVideoDTO)

	return res.MediaID, nil
}

func GetIGSessionID(
	app *pocketbase.PocketBase,
	args GetIGSessionTokenArgs,
) (string, error) {
	logger := app.Logger().WithGroup("GetIGSessionID")

	envVars, err := config.ParseEnv()
	if err != nil {
		logger.Error("Couldn't parse env vars, returning")
		return "", err
	}

	bodyBytes, err := json.Marshal(args)
	if err != nil {
		logger.Error("Couldn't marshal JSON body, returning")
		return "", err
	}

	client := resty.New()

	resp, err := client.R().
		SetBody(bodyBytes).
		SetResult(&IGSessionTokenDTO{}).
		SetHeader("X-Secret", envVars.IGServiceSecret).
		SetHeader("Content-Type", "application/json").
		Post(envVars.IGServiceURL + "/getSessionId")

	if err != nil {
		logger.Error("Error sending request", "err", err)
		return "", err
	}

	statusCode := resp.StatusCode()
	if statusCode > 399 {
		logger.Error("invalid status code " + strconv.Itoa(statusCode) + ", response: " + resp.String())
		return "", errors.New("invalid status code " + strconv.Itoa(statusCode))
	}

	res := resp.Result().(*IGSessionTokenDTO)

	return res.SessionID, nil
}
