package igservice

import (
	"brfactorybackend/internal/config"
	"errors"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type UploadIGTVVideoArgs struct {
	SessionID    string
	Title        string
	Caption      string
	ThumbnailURL string
	VideoURL     string
}

type GetIGSessionTokenArgs struct {
	IGUsername string
	IOPassword string
}
type IGSessionTokenDTO struct {
	SessionID string `json:"sessionId"`
}

func UploadIGTVVideo(args UploadIGTVVideoArgs) error {
	envVars, err := config.ParseEnv()
	if err != nil {
		log.Println("Couldn't parse env vars, returning")
		return err
	}

	body := `{"sessionId": "` + args.SessionID +
		`", "videoURL": "` + args.VideoURL +
		`", "title": "` + args.Title +
		`", "caption": "` + args.Caption +
		`", "thumbnailURL": "` + args.ThumbnailURL + `"}`

	client := resty.New()

	resp, err := client.R().
		SetBody(body).
		SetHeader("X-Secret", envVars.IGServiceSecret).
		SetHeader("Content-Type", "application/json").
		Post(envVars.IGServiceURL + "/uploadIGTVVideo")

	if err != nil {
		return err
	}

	statusCode := resp.StatusCode()
	if statusCode > 399 {
		log.Println("invalid status code " + strconv.Itoa(statusCode) + ", response: " + resp.String())
		return errors.New("invalid status code " + strconv.Itoa(statusCode))
	}

	return nil
}

func GetIGSessionID(args GetIGSessionTokenArgs) (string, error) {
	envVars, err := config.ParseEnv()
	if err != nil {
		log.Println("Couldn't parse env vars, returning")
		return "", err
	}

	body := `{"igUsername": "` + args.IGUsername + `", "igPassword": "` + args.IOPassword + `"}`

	client := resty.New()

	resp, err := client.R().
		SetBody(body).
		SetResult(&IGSessionTokenDTO{}).
		SetHeader("X-Secret", envVars.IGServiceSecret).
		SetHeader("Content-Type", "application/json").
		Post(envVars.IGServiceURL + "/getSessionId")

	if err != nil {
		return "", err
	}

	statusCode := resp.StatusCode()
	if statusCode > 399 {
		log.Println("invalid status code " + strconv.Itoa(statusCode) + ", response: " + resp.String())
		return "", errors.New("invalid status code " + strconv.Itoa(statusCode))
	}

	sessionTokenResponse := resp.Result().(*IGSessionTokenDTO)

	return sessionTokenResponse.SessionID, nil
}
