package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"brfactorybackend/internal/config"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
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

	// lastGeneratedVoiceID, err := GetLastGeneratedVoiceID()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("lastGeneratedVoiceID", lastGeneratedVoiceID)

	// DownloadVoice(lastGeneratedVoiceID, "audio.mp3", VoiceDownloadTypeMP3)

	test()

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

func GetLastGeneratedVoiceID() (string, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return "", err
	}

	cookies, err := getCookies()
	if err != nil {
		log.Println("Error getting last generated voice id, getting cookies", err)
		return "", err
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
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error getting last generated voice id, sending request", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("invalid status code", resp.StatusCode, "response:", resp.Status)
		return "", errors.New("invalid status code " + strconv.Itoa(resp.StatusCode))
	}

	var res MicmonsterListVoicesDTO
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println("Error decoding response", err)
		return "", err
	}

	if len(res.Voices) > 0 {
		return res.Voices[0].ID, nil
	}

	log.Println("Error getting last generated voice id, no voices found")
	return "", nil
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

func test() {
	url := "https://app.micmonster.com/downloads/5473291?downloadType=wav"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("sec-ch-ua", `"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("Referer", "https://app.micmonster.com/project-detail/289da212-ce9d454c-b9c1dcad-ecee23a1")
	req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")

	attachCookiesToRequest(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			return
		}
		log.Println("Response body:", string(body))
		log.Println("Failed to download file, status code:", resp.StatusCode)
		return
	}

	outFile, err := os.Create("downloaded_audio.wav")
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		log.Println("Error copying response body to file:", err)
		return
	}

	log.Println("Audio file downloaded successfully")
}

// func DownloadVoice(voiceID, savePath string, voiceDownloadType VoiceDownloadType) error {
// 	env, err := config.ParseEnv()
// 	if err != nil {
// 		return err
// 	}

// 	req, err := http.NewRequest(
// 		"GET",
// 		fmt.Sprintf("%s/downloads/%s?downloadType=%s", env.MicmonsterApiURL, voiceID, string(voiceDownloadType)),
// 		nil,
// 	)

// 	if err != nil {
// 		log.Println("Error creating request:", err)
// 		return err
// 	}

// 	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
// 	req.Header.Set("accept-language", "en-US,en;q=0.9")
// 	req.Header.Set("sec-ch-ua", `"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"`)
// 	req.Header.Set("sec-ch-ua-mobile", "?0")
// 	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
// 	req.Header.Set("sec-fetch-dest", "document")
// 	req.Header.Set("sec-fetch-mode", "navigate")
// 	req.Header.Set("sec-fetch-site", "same-origin")
// 	req.Header.Set("sec-fetch-user", "?1")
// 	req.Header.Set("upgrade-insecure-requests", "1")
// 	req.Header.Set("cookie", )
// 	req.Header.Set("Referer", "https://app.micmonster.com/project-detail/289da212-ce9d454c-b9c1dcad-ecee23a1")
// 	req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.Println("Error sending request:", err)
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		body, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			log.Println("Error reading response body:", err)
// 			return err
// 		}
// 		log.Println("Response body:", string(body))
// 		log.Println("Failed to download file, status code:", resp.StatusCode)
// 		return err
// 	}

// 	outFile, err := os.Create(savePath)
// 	if err != nil {
// 		log.Println("Error creating file:", err)
// 		return err
// 	}
// 	defer outFile.Close()

// 	_, err = io.Copy(outFile, resp.Body)
// 	if err != nil {
// 		log.Println("Error copying response body to file:", err)
// 		return err
// 	}

// 	log.Println("Audio file downloaded successfully")

// 	return nil
// }
