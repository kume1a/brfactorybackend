package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"brfactorybackend/internal/config"
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/google/uuid"

	"crypto/hmac"
	"encoding/hex"
)

type MicmonsterVoiceDTO struct {
	ID             string      `json:"id"`
	UserID         string      `json:"user_id"`
	ProjectID      string      `json:"project_id"`
	Text           string      `json:"text"`
	HTML           interface{} `json:"html"`
	Characters     string      `json:"characters"`
	Audio          string      `json:"audio"`
	Language       string      `json:"language"`
	VoiceUsed      string      `json:"voice_used"`
	VoiceGender    string      `json:"voice_gender"`
	UserVoiceName  string      `json:"user_voice_name"`
	IsDeleted      string      `json:"is_deleted"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at"`
	IsUpload       string      `json:"is_upload"`
	ConversionData interface{} `json:"conversion_data"`
	TextMergingID  string      `json:"text_merging_id"`
	MergingStatus  interface{} `json:"merging_status"`
	Status         interface{} `json:"status"`
	AudioQuality   interface{} `json:"audio_quality"`
	IsPending      int         `json:"is_pending"`
}

type MicmonsterListVoicesDTO struct {
	Voices       []MicmonsterVoiceDTO `json:"voices"`
	TotalRecords int                  `json:"totalRecords"`
	TotalPages   int                  `json:"totalPages"`
	CurrentPage  int                  `json:"currentPage"`
	Start        string               `json:"start"`
	Limit        string               `json:"limit"`
	Search       string               `json:"search"`
	SortBy       string               `json:"sortBy"`
	SortOrder    string               `json:"sortOrder"`
	Data         struct {
		Token string `json:"token"`
	} `json:"data"`
}

// Step 1: Split text into sentences and handle long sentences
func splitTextIntoSentences(text string) []string {
	re := regexp.MustCompile(`(?m)([^.!?]+[.!?]*)`)
	sentences := re.FindAllString(text, -1)

	var result []string
	for _, sentence := range sentences {
		if len(sentence) > 200 { // Assuming 200 characters as the threshold for a long sentence
			parts := splitLongSentence(sentence, 200)
			result = append(result, parts...)
		} else {
			result = append(result, sentence)
		}
	}
	return result
}

// Helper function to split long sentences into smaller parts
func splitLongSentence(sentence string, maxLength int) []string {
	words := regexp.MustCompile(`\s+`).Split(sentence, -1)
	var parts []string
	var currentPart string

	for _, word := range words {
		if len(currentPart)+len(word)+1 > maxLength {
			parts = append(parts, currentPart)
			currentPart = word
		} else {
			if currentPart != "" {
				currentPart += " "
			}
			currentPart += word
		}
	}
	if currentPart != "" {
		parts = append(parts, currentPart)
	}
	return parts
}

// Step 3: Get duration of audio file
func getAudioDuration(audioFile string) (float64, error) {
	cmd := exec.Command("ffmpeg", "-i", audioFile, "-f", "null", "-")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	// Use regex to extract the duration from the ffmpeg output
	re := regexp.MustCompile(`Duration: (\d+):(\d+):(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) != 5 {
		return 0, fmt.Errorf("could not parse duration from ffmpeg output")
	}

	hours, _ := strconv.Atoi(matches[1])
	minutes, _ := strconv.Atoi(matches[2])
	seconds, _ := strconv.Atoi(matches[3])
	milliseconds, _ := strconv.Atoi(matches[4])

	duration := float64(hours*3600+minutes*60+seconds) + float64(milliseconds)/1000
	return duration, nil
}

// Step 4: Create SRT file
func createSRTFile(sentences []string, durations []float64, srtFile string) error {
	file, err := os.Create(srtFile)
	if err != nil {
		return err
	}
	defer file.Close()

	currentTime := 0.0
	for i, sentence := range sentences {
		startTime := formatTimeForSRT(currentTime)
		currentTime += durations[i]
		endTime := formatTimeForSRT(currentTime)

		srtEntry := fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, startTime, endTime, sentence)
		file.WriteString(srtEntry)
	}
	return nil
}

// Format time for SRT subtitles
func formatTimeForSRT(seconds float64) string {
	hrs := int(seconds / 3600)
	mins := int(seconds/60) % 60
	secs := int(seconds) % 60
	millis := int((seconds - float64(int(seconds))) * 1000)
	return fmt.Sprintf("%02d:%02d:%02d,%03d", hrs, mins, secs, millis)
}

// Step 5: Stitch video with audio and subtitles
func stitchVideoWithAudioAndSubtitles(videoFile string, audioFiles []string, srtFile string, outputFile string) error {
	audioInputs := []string{}
	for _, audioFile := range audioFiles {
		audioInputs = append(audioInputs, "-i", audioFile)
	}

	cmd := exec.Command("ffmpeg", append(audioInputs, "-i", videoFile, "-vf", fmt.Sprintf("subtitles=%s", srtFile), outputFile)...)
	return cmd.Run()
}

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Fatal("Couldn't load env vars, returning")
	}

	// story := "This is the first sentence. This is the second sentence."
	// sentences := splitTextIntoSentences(story)

	// log.Println("sentences", sentences)
	// log.Println("len", len(sentences))

	lastGeneratedAudio, err := GetLastGeneratedVoice()
	if err != nil {
		log.Fatal(err)
	}

	DownloadAudio(lastGeneratedAudio.Audio, "audio.mp3")

	// var durations []float64
	// var audioFiles []string

	// // Fetch TTS for each sentence
	// for i, sentence := range sentences {
	// 	filename := "audio_" + strconv.Itoa(i+1) + ".mp3"
	// 	audioFiles = append(audioFiles, filename)

	// 	audioID, err := GenerateVoice(sentence)
	// 	if err != nil {
	// 		fmt.Println("Error generating audio:", err)
	// 		return
	// 	}

	// 	log.Println("audioID", audioID)

	// 	duration, err := getAudioDuration(filename)
	// 	if err != nil {
	// 		fmt.Println("Error getting audio duration:", err)
	// 		return
	// 	}
	// 	durations = append(durations, duration)
	// }

	// // Create the SRT file
	// srtFile := "output.srt"
	// err := createSRTFile(sentences, durations, srtFile)
	// if err != nil {
	// 	fmt.Println("Error creating SRT file:", err)
	// 	return
	// }

	// // Stitch video with audio and subtitles
	// videoFile := "background.mp4"
	// outputFile := "final_output.mp4"
	// err = stitchVideoWithAudioAndSubtitles(videoFile, audioFiles, srtFile, outputFile)
	// if err != nil {
	// 	fmt.Println("Error stitching video:", err)
	// } else {
	// 	fmt.Println("Video created successfully:", outputFile)
	// }
}

// func NewMicmonsterClient(environmentVariableProvider EnvironmentVariableProvider) *MicmonsterClient {
// 	jar, _ := cookiejar.New(nil)
// 	client := &http.Client{Jar: jar}
// 	return &MicmonsterClient{
// 		client:                      client,
// 		environmentVariableProvider: environmentVariableProvider,
// 	}
// }

func GenerateVoice(text string) (bool, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return false, err
	}

	data := url.Values{}
	data.Set("text", text)
	data.Set("language", "English (US)")
	data.Set("language_code", "en-US")
	data.Set("voice", "en-US-ChristopherNeural")
	data.Set("voice_gender", "Male")
	data.Set("project_id", env.MicmonsterProjectID)
	data.Set("project_id_plain", env.MicmonsterProjectIDPlain)
	data.Set("humanname", "Christopher")
	data.Set("audio_quality", "medium")
	data.Set("tts_type", "MS")
	data.Set("voice_style", "")
	data.Set("msspeed", "1")
	data.Set("mspitch", "1")
	data.Set("csrf_test_name", "acf58af7136d4c79805aaad5cd0f51e5")
	data.Set("featureType", "generate")
	data.Set("user_voice_name", uuid.New().String())

	req, err := http.NewRequest("POST", env.MicmonsterApiURL+"/generate-voice", bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Println("Error generating voice, creating request", err)
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err := attachCookiesToRequest(req); err != nil {
		log.Println("Error generating voice, attaching cookies", err)
		return false, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error generating voice", err)
		return false, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if status, ok := result["status"].(string); ok && status == "success" {
		return true, nil
	}

	log.Println("Error generating voice, status not success")
	return false, errors.New("error generating voice, status not success")
}

// func (mc *MicmonsterClient) DeleteVoice(voiceId string) error {
// 	data := url.Values{}
// 	data.Set("id", voiceId)

// 	req, err := http.NewRequest("POST", mc.environmentVariableProvider.MicMonsterApiUrl()+"/delete-voice", bytes.NewBufferString(data.Encode()))
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	mc.addCookies(req)

// 	_, err = mc.client.Do(req)
// 	return err
// }

type MicmonsterListVoicesBodyDTO struct {
	Start          string `json:"start"`
	Limit          string `json:"limit"`
	SortBy         string `json:"sortBy"`
	SortOrder      string `json:"sortOrder"`
	ProjectID      string `json:"project_id"`
	ProjectIDPlain string `json:"project_id_plain"`
	CsrfTestName   string `json:"csrf_test_name"`
	Timezone       string `json:"timezone"`
}

func GetLastGeneratedVoice() (MicmonsterVoiceDTO, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return MicmonsterVoiceDTO{}, err
	}

	cookies, err := getCookies()
	if err != nil {
		log.Println("Error getting last generated voice id, getting cookies", err)
		return MicmonsterVoiceDTO{}, err
	}

	data := url.Values{}
	data.Set("start", "0")
	data.Set("limit", "1")
	data.Set("sortBy", "created_at")
	data.Set("sortOrder", "DESC")
	data.Set("project_id", env.MicmonsterProjectID)
	data.Set("project_id_plain", env.MicmonsterProjectIDPlain)
	data.Set("csrf_test_name", "5d0f455a0c885b5d26401731eae2b85a")
	data.Set("timezone", "Asia/Tbilisi")

	req, err := http.NewRequest("POST", env.MicmonsterApiURL+"/list-voices", bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Println("Error getting last generated voice id, creating request", err)
		return MicmonsterVoiceDTO{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error getting last generated voice id, sending request", err)
		return MicmonsterVoiceDTO{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("invalid status code", resp.StatusCode, "response:", resp.Status)
		return MicmonsterVoiceDTO{}, errors.New("invalid status code " + strconv.Itoa(resp.StatusCode))
	}

	var res MicmonsterListVoicesDTO
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println("Error decoding response", err)
		return MicmonsterVoiceDTO{}, err
	}

	if len(res.Voices) > 0 {
		return res.Voices[0], nil
	}

	log.Println("Error getting last generated voice id, no voices found")
	return MicmonsterVoiceDTO{}, nil
}

type VoiceDownloadType string

const (
	VoiceDownloadTypeMP3 VoiceDownloadType = "mp3"
	VoiceDownloadTypeWAV VoiceDownloadType = "wav"
)

func attachCookiesToRequest(req *http.Request) error {
	cookies, err := getCookies()
	if err != nil {
		return err
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	return nil
}

func getCookies() ([]*http.Cookie, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return nil, err
	}

	return []*http.Cookie{
		{Name: "ci_session", Value: env.MicmonsterCISession},
		{Name: "cid", Value: env.MicmonsterCID},
		{Name: "cpass", Value: env.MicmonsterCPASS},
	}, nil
}

func DownloadAudio(audioPath, filename string) error {
	env, err := config.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	host := "micmonsterlive.s3.us-east-2.amazonaws.com"
	region := "us-east-2"
	service := "s3"
	expires := 1200

	t := time.Now().UTC()
	date := t.Format("20060102")
	amzDate := t.Format("20060102T150405Z")

	canonicalURI := "/uploads/audio/" + audioPath
	canonicalQueryString := "X-Amz-Algorithm=AWS4-HMAC-SHA256"
	canonicalQueryString += "&X-Amz-Credential=" + url.QueryEscape(env.MicmonsterAWSAccessKey+"/"+date+"/"+region+"/"+service+"/aws4_request")
	canonicalQueryString += "&X-Amz-Date=" + amzDate
	canonicalQueryString += "&X-Amz-Expires=" + fmt.Sprintf("%d", expires)
	canonicalQueryString += "&X-Amz-SignedHeaders=host"

	canonicalHeaders := "host:" + host + "\n"
	signedHeaders := "host"
	payloadHash := "UNSIGNED-PAYLOAD"

	canonicalRequest := strings.Join([]string{
		"GET",
		canonicalURI,
		canonicalQueryString,
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	}, "\n")

	// Generate the string to sign
	algorithm := "AWS4-HMAC-SHA256"
	credentialScope := date + "/" + region + "/" + service + "/" + "aws4_request"
	stringToSign := strings.Join([]string{
		algorithm,
		amzDate,
		credentialScope,
		hex.EncodeToString(hashSHA256([]byte(canonicalRequest))),
	}, "\n")

	// Generate the signing key
	signingKey := getSignatureKey(env.MicmonsterAWSSecretKey, date, region, service)

	// Generate the signature
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))

	// Create the final URL
	finalURL := fmt.Sprintf("https://%s%s?%s&X-Amz-Signature=%s", host, canonicalURI, canonicalQueryString, signature)
	fmt.Println("Generated URL:", finalURL)

	resp, err := http.Get(finalURL)
	if err != nil {
		log.Println("Error getting file", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error downloading file", resp.Status)
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	out, err := os.Create(filename)
	if err != nil {
		log.Println("Error creating file", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func hashSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func getSignatureKey(secretKey, date, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secretKey), date)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	kSigning := hmacSHA256(kService, "aws4_request")
	return kSigning
}
