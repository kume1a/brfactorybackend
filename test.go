package main

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
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

	"crypto/hmac"
	"encoding/hex"

	"github.com/google/uuid"
)

type MicmonsterAudioDTO struct {
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

type MicmonsterListAudiosDTO struct {
	Voices       []MicmonsterAudioDTO `json:"voices"`
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

func splitIntoChunks(text string) []string {
	chunkLength := 50

	re := regexp.MustCompile(`(?m)([^.!?]+[.!?]*)`)
	sentences := re.FindAllString(text, -1)

	var result []string
	var currentChunk string

	for _, sentence := range sentences {
		if !isValidSentence(sentence) {
			continue
		}

		if len(currentChunk)+len(sentence)+1 > chunkLength {
			result = append(result, strings.TrimSpace(currentChunk))
			currentChunk = sentence
		} else {
			if currentChunk != "" {
				currentChunk += " "
			}
			currentChunk += sentence
		}
	}

	if currentChunk != "" {
		result = append(result, strings.TrimSpace(currentChunk))
	}

	return result
}

func isValidSentence(sentence string) bool {
	re := regexp.MustCompile(`[a-zA-Z]`)
	return re.MatchString(sentence)
}

func getAudioDuration(audioFile string) (float64, error) {
	cmd := exec.Command("ffmpeg", "-i", audioFile, "-f", "null", "-")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

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

func formatTimeForSRT(seconds float64) string {
	hrs := int(seconds / 3600)
	mins := int(seconds/60) % 60
	secs := int(seconds) % 60
	millis := int((seconds - float64(int(seconds))) * 1000)
	return fmt.Sprintf("%02d:%02d:%02d,%03d", hrs, mins, secs, millis)
}

func stitchVideoWithAudioAndSubtitles(videoFile string, audioFiles []string, srtFile string, outputFile string) error {
	audioInputs := []string{}
	for _, audioFile := range audioFiles {
		audioInputs = append(audioInputs, "-i", audioFile)
	}

	commandArgs := append(audioInputs, "-i", videoFile, "-vf", fmt.Sprintf("subtitles=%s", srtFile), outputFile)

	log.Println("ffmpeg command", commandArgs)

	cmd := exec.Command("ffmpeg", commandArgs...)
	return cmd.Run()
}

func ChangeAudioSpeed(inputPath, outputPath string, speed float64) error {
	if speed < 0.5 || speed > 2.5 {
		return fmt.Errorf("speed must be between 0.5 and 2.0")
	}

	cmd := exec.Command("ffmpeg", "-i", inputPath, "-filter:a", fmt.Sprintf("atempo=%f", speed), "-vn", outputPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to change audio speed: %v", err)
	}

	return nil
}

func GenerateDownloadAudios(texts []string) (files []string, err error) {
	// for i, sentence := range sentences {
	// 	filename := "data/audio_" + strconv.Itoa(i+1) + ".mp3"
	// 	audioFiles = append(audioFiles, filename)

	// 	log.Println("generating audio, filename=" + filename + ", text=" + sentence)

	// 	_, err := GenerateAudio(sentence)
	// 	if err != nil {
	// 		fmt.Println("Error generating audio:", err)
	// 		return
	// 	}

	// 	lastGeneratedAudio, err := GetLastGeneratedAudio()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	DownloadAudio(lastGeneratedAudio.Audio, filename)

	// 	duration, err := getAudioDuration(filename)
	// 	if err != nil {
	// 		fmt.Println("Error getting audio duration:", err)
	// 		return
	// 	}
	// 	durations = append(durations, duration)
	// }

	return []string{}, nil
}

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Fatal("Couldn't load env vars, returning")
	}

	// audio, err := GetLastGeneratedAudio()

	// if err != nil {
	// 	log.Fatal("Couldn't get last generated audio, returning")
	// }

	// log.Println(audio)
	// return

	story := "Nursing student. Born and raised on a farm. Twenty-eight years old. Slim. Defensive posture and a soft voice (low self-esteem). Sitting in a shitty bar at 8 PM on a Wednesday night. Physically, she looked a lot like the last one.\nShe was perfect.\nAfter some small talk and four cans of beer, her voice softens. Her cheeks flush, and I can sense the sexual tension building. I shift to the offensive, leaning closer and letting my hand brush her knee and shoulder. I wait for a reaction. It comes in the form of a shy glance and a slight openness to more physical contact.\nEverything was going according to plan. In fact, I hadn’t expected things to flow this smoothly. Two out of the last three had required more than one encounter to reach this level of intimacy.\nI invite her to leave the bar and grab a quick bite to eat. I tell her I know a great spot nearby where we can get something fast before calling it a night. I feel her hesitation for a moment, probably weighing the risks of saying yes. She mentions she has class early tomorrow, but I reassure her—it won’t take long. The place is just fifteen minutes away. Convinced, she follows me out and into my car.\nWe laugh and chat during the drive. She only realizes we’ve entered the park about ten minutes in and asks if we’re close to the destination. I assure her that we are—it’s just up ahead. I just took a shortcut.\nWe’re now deep inside the park, where the lights become sparse and then disappear completely. It’s the perfect place—one I know like the back of my hand. I had practiced this route several times to ensure everything would run smoothly and avoid any unexpected encounters. All the others ended up here too.\nI pull over at the pre-planned spot and ask her to step out of the car. Confusion spreads across her face as she senses something is wrong, and her body stiffens. I open the door and yank her out, and she collapses onto the grass.\nGrabbing her by the neck, I steer her along the path. She starts begging for mercy, sobbing uncontrollably now. I ignore her and continue down the short trail toward my usual location.\nOnce there, I throw her to the ground and tie her hands with the rope I’d left ready. As I reach for the knife I had buried nearby, a sharp, burning pain stabs my side, and I lose balance.\nOn the ground, I realize I’ve been shot in the thigh. A man steps out from the shadows with a shotgun in hand and unties the girl. They embrace, and I hear him say, \"This is what Catherine would have wanted. Now she can rest in peace—you were perfect.”. The resemblance hit me like a jolt, and in that moment, I remembered—Catherine, the last girl.\nHe reloads the shotgun and steps toward me. I try to reason with him, plead for calm, explaining that it’s all a misunderstanding. But the cold steel of the barrel presses against my forehead."
	sentences := splitIntoChunks(story)

	log.Println("sentences len", len(sentences))

	durations := []float64{}
	audioFiles := []string{
		"data/audio_1_1_5x.mp3", "data/audio_2_1_5x.mp3", "data/audio_3_1_5x.mp3", "data/audio_4_1_5x.mp3", "data/audio_5_1_5x.mp3",
		"data/audio_6_1_5x.mp3", "data/audio_7_1_5x.mp3", "data/audio_8_1_5x.mp3", "data/audio_9_1_5x.mp3", "data/audio_10_1_5x.mp3",
		"data/audio_11_1_5x.mp3", "data/audio_12_1_5x.mp3", "data/audio_13_1_5x.mp3", "data/audio_14_1_5x.mp3", "data/audio_15_1_5x.mp3",
		"data/audio_16_1_5x.mp3", "data/audio_17_1_5x.mp3", "data/audio_18_1_5x.mp3", "data/audio_19_1_5x.mp3", "data/audio_20_1_5x.mp3",
		"data/audio_21_1_5x.mp3", "data/audio_22_1_5x.mp3", "data/audio_23_1_5x.mp3", "data/audio_24_1_5x.mp3", "data/audio_25_1_5x.mp3",
		"data/audio_26_1_5x.mp3", "data/audio_27_1_5x.mp3", "data/audio_28_1_5x.mp3", "data/audio_29_1_5x.mp3", "data/audio_30_1_5x.mp3",
		"data/audio_31_1_5x.mp3", "data/audio_32_1_5x.mp3", "data/audio_33_1_5x.mp3", "data/audio_34_1_5x.mp3", "data/audio_35_1_5x.mp3",
		"data/audio_36_1_5x.mp3", "data/audio_37_1_5x.mp3", "data/audio_38_1_5x.mp3", "data/audio_39_1_5x.mp3", "data/audio_40_1_5x.mp3",
		"data/audio_41_1_5x.mp3", "data/audio_42_1_5x.mp3", "data/audio_43_1_5x.mp3", "data/audio_44_1_5x.mp3",
	}

	for _, audioFile := range audioFiles {
		// ext := filepath.Ext(audioFile)
		// base := audioFile[:len(audioFile)-len(ext)]
		// outputPath := fmt.Sprintf("%s_1_5x%s", base, ext)

		// log.Println("changing audio speed", audioFile, outputPath)

		// err := ChangeAudioSpeed(audioFile, outputPath, 1.5)
		// if err != nil {
		// 	log.Println("Error changing audio speed", err)
		// 	return
		// }

		duration, err := getAudioDuration(audioFile)
		if err != nil {
			log.Println("Error getting audio duration:", err)
			return
		}

		durations = append(durations, duration)
	}

	srtFile := "data/output.srt"
	if err := createSRTFile(sentences, durations, srtFile); err != nil {
		fmt.Println("Error creating SRT file:", err)
		return
	}

	backgroundVideo := "data/background.mp4"
	outputVideo := "data/final_output_" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".mp4"

	file, err := os.Open(srtFile)
	if err != nil {
		fmt.Println("Error opening subtitle file:", err)
		return
	}
	defer file.Close()

	timestampRegex := regexp.MustCompile(`([0-9]{2}):([0-9]{2}):([0-9]{2}),([0-9]{3})`)
	inputFiles := []string{"-i", backgroundVideo}
	filterComplex := ""
	audioIndex := 1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := timestampRegex.FindStringSubmatch(line); len(matches) == 5 {
			hours, _ := strconv.Atoi(matches[1])
			minutes, _ := strconv.Atoi(matches[2])
			seconds, _ := strconv.Atoi(matches[3])
			milliseconds, _ := strconv.Atoi(matches[4])
			delay := (hours*3600+minutes*60+seconds)*1000 + milliseconds

			filterComplex += fmt.Sprintf("[%d]adelay=%d|%d[a%d]; ", audioIndex, delay, delay, audioIndex)
			inputFiles = append(inputFiles, "-i", fmt.Sprintf("data/audio_%d_1_5x.mp3"))
			audioIndex++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading subtitle file:", err)
		return
	}

	amixInputs := ""
	for i := 1; i < audioIndex; i++ {
		amixInputs += fmt.Sprintf("[a%d]", i)
	}
	filterComplex += fmt.Sprintf("%samix=inputs=%d[audio]", amixInputs, audioIndex-1)

	subtitleFilter := fmt.Sprintf("subtitles='%s':force_style='FontSize=16,Bold=1,Alignment=10,MarginV=0,OutlineColour=&H000000&,Outline=1'", srtFile)
	videoFilter := "scale=1080:1920:force_original_aspect_ratio=increase,crop=1080:1920"
	finalFilter := fmt.Sprintf("%s,%s", videoFilter, subtitleFilter)
	ffmpegArgs := append(inputFiles, "-filter_complex", filterComplex, "-vf", finalFilter, "-map", "0:v", "-map", "[audio]", "-shortest", outputVideo)

	fmt.Println("Running ffmpeg command:")
	fmt.Println("ffmpeg", ffmpegArgs)

	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error executing ffmpeg:", err)
	}
}

// func test() {

// }

func GenerateAudio(text string) (bool, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return false, err
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	w.WriteField("text", text)
	w.WriteField("language", "English (US)")
	w.WriteField("language_code", "en-US")
	w.WriteField("voice", "en-US-ChristopherNeural")
	w.WriteField("voice_gender", "Male")
	w.WriteField("project_id", env.MicmonsterProjectID)
	w.WriteField("project_id_plain", env.MicmonsterProjectIDPlain)
	w.WriteField("humanname", "Christopher")
	w.WriteField("audio_quality", "medium")
	w.WriteField("tts_type", "MS")
	w.WriteField("voice_style", "")
	w.WriteField("msspeed", "1")
	w.WriteField("mspitch", "1")
	w.WriteField("csrf_test_name", "8fcdfeb2e90e21b0fe880093f257962d")
	w.WriteField("featureType", "generate")
	w.WriteField("user_voice_name", uuid.New().String())

	w.Close()

	req, err := http.NewRequest("POST", env.MicmonsterApiURL+"/generate-voice", &b)
	if err != nil {
		log.Println("Error creating request", err)
		return false, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Add("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("sec-ch-ua", `"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", `"Windows"`)
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	req.Header.Add("cookie", "cid="+env.MicmonsterCID+"; cpass="+env.MicmonsterCPASS+"; isQuickStartSeen=yes; ci_session=ct35g1lvbtif62llsqrs7jju4j33tnaj")
	req.Header.Add("Referer", env.MicmonsterApiURL+"/project-detail/"+env.MicmonsterProjectID)
	req.Header.Add("Referrer-Policy", "strict-origin-when-cross-origin")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending request", err)
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("invalid status code", res.StatusCode, "response:", res.Status)
		return false, errors.New("invalid status code " + strconv.Itoa(res.StatusCode))
	}

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	if status, ok := result["status"].(string); ok && status == "success" {
		log.Println("Success generating audio")
		return true, nil
	}

	log.Println("Error generating audio, status not success")
	return false, errors.New("error generating audio, status not success")
}

// func (mc *MicmonsterClient) DeleteAudio(audioID string) error {
// 	data := url.Values{}
// 	data.Set("id", audioID)

// 	req, err := http.NewRequest("POST", mc.environmentVariableProvider.MicMonsterApiUrl()+"/delete-voice", bytes.NewBufferString(data.Encode()))
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	mc.addCookies(req)

// 	_, err = mc.client.Do(req)
// 	return err
// }

type MicmonsterListAudiosBodyDTO struct {
	Start          string `json:"start"`
	Limit          string `json:"limit"`
	SortBy         string `json:"sortBy"`
	SortOrder      string `json:"sortOrder"`
	ProjectID      string `json:"project_id"`
	ProjectIDPlain string `json:"project_id_plain"`
	CsrfTestName   string `json:"csrf_test_name"`
	Timezone       string `json:"timezone"`
}

func GetLastGeneratedAudio() (MicmonsterAudioDTO, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return MicmonsterAudioDTO{}, err
	}

	cookies, err := getCookies()
	if err != nil {
		log.Println("Error getting last generated audio id, getting cookies", err)
		return MicmonsterAudioDTO{}, err
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
		log.Println("Error getting last generated audio id, creating request", err)
		return MicmonsterAudioDTO{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error getting last generated audio id, sending request", err)
		return MicmonsterAudioDTO{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("invalid status code", resp.StatusCode, "response:", resp.Status)
		return MicmonsterAudioDTO{}, errors.New("invalid status code " + strconv.Itoa(resp.StatusCode))
	}

	var res MicmonsterListAudiosDTO
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println("Error decoding response", err)
		return MicmonsterAudioDTO{}, err
	}

	if len(res.Voices) > 0 {
		return res.Voices[0], nil
	}

	log.Println("Error getting last generated audio id, no audios found")
	return MicmonsterAudioDTO{}, nil
}

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

	algorithm := "AWS4-HMAC-SHA256"
	credentialScope := date + "/" + region + "/" + service + "/" + "aws4_request"
	stringToSign := strings.Join([]string{
		algorithm,
		amzDate,
		credentialScope,
		hex.EncodeToString(hashSHA256([]byte(canonicalRequest))),
	}, "\n")

	signingKey := getSignatureKey(env.MicmonsterAWSSecretKey, date, region, service)
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))

	finalURL := fmt.Sprintf("https://%s%s?%s&X-Amz-Signature=%s", host, canonicalURI, canonicalQueryString, signature)

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
		log.Println("Error copying file", err)
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
