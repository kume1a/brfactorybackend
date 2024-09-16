package igservice

import (
	"brfactorybackend/internal/config"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
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

func UploadIGTVVideo(args UploadIGTVVideoArgs) (string, error) {
	envVars, err := config.ParseEnv()
	if err != nil {
		log.Println("Couldn't parse env vars, returning")
		return "", err
	}

	bodyBytes, err := json.Marshal(args)
	if err != nil {
		log.Println("Couldn't marshal JSON body, returning")
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
		return "", err
	}

	statusCode := resp.StatusCode()
	if statusCode > 399 {
		log.Println("invalid status code " + strconv.Itoa(statusCode) + ", response: " + resp.String())
		return "", errors.New("invalid status code " + strconv.Itoa(statusCode))
	}

	res := resp.Result().(*UploadIGTVVideoDTO)

	return res.MediaID, nil
}

func GetIGSessionID(args GetIGSessionTokenArgs) (string, error) {
	envVars, err := config.ParseEnv()
	if err != nil {
		log.Println("Couldn't parse env vars, returning")
		return "", err
	}

	bodyBytes, err := json.Marshal(args)
	if err != nil {
		log.Println("Couldn't marshal JSON body, returning")
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
		return "", err
	}

	statusCode := resp.StatusCode()
	if statusCode > 399 {
		log.Println("invalid status code " + strconv.Itoa(statusCode) + ", response: " + resp.String())
		return "", errors.New("invalid status code " + strconv.Itoa(statusCode))
	}

	res := resp.Result().(*IGSessionTokenDTO)

	return res.SessionID, nil
}
