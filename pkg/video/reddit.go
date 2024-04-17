package video

import (
	"fmt"
	"math/rand"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/notarock/slopify/pkg/google"
	"github.com/notarock/slopify/pkg/reddit"
	"github.com/notarock/slopify/pkg/subs"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const SUBTITLE_COMMAND = "subtitles=%s:force_style='FontSize=24,Alignment=10'"
const BUCKET_NAME = "slopify-transcription-buffer"
const SUBTITLES_FILE = "output.srt"

type RedditThreadVideoInput struct {
	ThreadURL        string
	WorkingDirectory string
	FootageDirectory string
}

func BuildFromThread(input RedditThreadVideoInput) (vid SlopVideo, err error) {
	fmt.Println("Building video from thread...")
	fmt.Println("Working directory: " + input.WorkingDirectory)
	fmt.Println("Footage directory: " + input.FootageDirectory)
	fmt.Println("Thread URL: " + input.ThreadURL)

	err = createWorkspace(input.WorkingDirectory)
	if err != nil {
		return vid, fmt.Errorf("Error creating workspace: %v", err)
	}

	_, err = url.ParseRequestURI(input.ThreadURL)

	// If it's a URL, scrape it
	// If it's a file, scrape it
	var thread reddit.Thread
	if err == nil {
		thread = reddit.ScrapeThread(input.ThreadURL)
	} else {
		err = nil
		thread = reddit.ScrapeFromFile(input.ThreadURL)
	}

	fmt.Println("Title: " + thread.Title)
	fmt.Println("Total comments: " + strconv.Itoa(thread.TotalComments()))
	if thread.Title == "" {
		return vid, fmt.Errorf("No title found in the thread, please provide a valid URL or file containing the HTML of the thread.")
	}

	titlePath := input.WorkingDirectory + "/audio/title.mp3"
	google.GetVoiceFile(thread.Title, titlePath)

	// TODO: Function to get the top comments and limit depth
	files := []string{}
	for _, value := range thread.CommentThreads {
		if len(files) > 0 {
			break
		}
		for _, comment := range value.Comments {
			if len(files) > 0 {
				break
			}

			filename := fmt.Sprintf(input.WorkingDirectory+"/"+"audio/%d.mp3", len(files))
			google.GetVoiceFile(comment, filename)
			files = append(files, filename)
		}
	}

	fullAudioPath := input.WorkingDirectory + "/output.mp3"

	titleAudioFile := input.WorkingDirectory + "/audio/title.mp3"
	google.Concatenate(titleAudioFile, files, fullAudioPath)

	if err != nil {
		return vid, fmt.Errorf("Error writing thread to file: %v", err)
	}

	var videoFile string

	if input.FootageDirectory != "" {
		path, err := pickFootageFromDirectory(input.FootageDirectory)
		if err != nil {
			return vid, fmt.Errorf("Error picking a video from the directory: %v", err)
		}
		videoFile = path
	} else {
		videoFile = "source.webm"
	}

	// TODO: Randomize the start time
	// This works when using an hour long footage of subway surfer

	// Get the resolution of the video

	probeOutput, err := probeVideoResolution(videoFile)
	if err != nil {
		return vid, fmt.Errorf("Error probing video resolution: %v", err)
	}

	height := probeOutput.Streams[0].Height
	width := probeOutput.Streams[0].Height / 16 * 9
	duration, err := strconv.ParseFloat(probeOutput.Format.Duration, 64)

	if err != nil {
		fmt.Printf("Error converting string to float64: %s\n", err)
		return
	}
	// Convert float64 to int (truncating the decimal part)
	minutes := int(duration) / 60

	randomWithSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := randomWithSeed.Intn(minutes - 1) // Generate a number between 0 and 55 (inclusive)
	fmt.Printf("Random is %d\n", randomNumber)

	kwArgs := []ffmpeg.KwArgs{
		ffmpeg.KwArgs{"shortest": ""},
		ffmpeg.KwArgs{"vf": fmt.Sprintf("crop=%d:%d", width, height)},
	}

	ffmpegVideo := ffmpeg.Input(videoFile, ffmpeg.KwArgs{"ss": fmt.Sprintf("00:%d:00", randomNumber)}).Video()
	ffmpegAudio := ffmpeg.Input(fullAudioPath).Audio()

	outputFile := input.WorkingDirectory + "/" + "slop-" + strconv.FormatInt(time.Now().Unix(), 10)
	outputFileWithSubs := outputFile + "-subs.mp4"
	outputFile = outputFile + ".mp4"

	out := ffmpeg.
		Output(
			[]*ffmpeg.Stream{ffmpegVideo, ffmpegAudio},
			outputFile,
			kwArgs...,
		)

	err = out.Run()
	if err != nil {
		return vid, fmt.Errorf("Error generating video: %v", err)
	}

	fmt.Println("Let's upload this to Google Cloud Storage")

	uri, err := google.UploadFile(BUCKET_NAME, outputFile)
	defer google.DeleteFile(BUCKET_NAME, filepath.Base(outputFile))
	if err != nil {
		return vid, fmt.Errorf("Error uploading video to Google Cloud Storage: %v", err)
	}

	fmt.Println("Uploaded to " + uri)

	fmt.Println("Transcribing video...")
	transcript, err := google.SpeechTranscriptionURI(uri)
	if err != nil {
		return vid, fmt.Errorf("Error transcribing video: %v", err)
	}

	titleDuration, err := subs.AudioDuration(titleAudioFile)
	if err != nil {
		return vid, fmt.Errorf("Error getting audio duration: %v", err)
	}
	fullTranscript := subs.BuildSubtitlesFromGoogle(transcript, titleDuration)

	fmt.Println("Converting to SRT format...")

	// Convert to SRT format
	srtData := subs.ConvertToSRT(fullTranscript)

	// Write to SRT file
	srtPath := input.WorkingDirectory + "/" + SUBTITLES_FILE
	err = subs.WriteSRT(srtPath, srtData)
	if err != nil {
		return vid, fmt.Errorf("Error writing SRT file: %v", err)
	}

	fmt.Println("Adding subtitles to video...")

	err = ffmpeg.
		Input(outputFile).
		Output(outputFileWithSubs, ffmpeg.KwArgs{"vf": fmt.Sprintf(SUBTITLE_COMMAND, srtPath)}).
		Run()

	if err != nil {
		return vid, fmt.Errorf("Error adding subtitles to video: %v", err)
	}

	// Create title image
	titleImage := CreateTitleCard(thread.Title, input.WorkingDirectory+"/title_image.png", width)
	completeFilePath := strings.Split(outputFile, ".")[0] + "-complete.mp4"

	err = ffmpeg.Input(outputFileWithSubs).Output(
		completeFilePath,
		ffmpeg.KwArgs{"i": titleImage},
		ffmpeg.KwArgs{"filter_complex": fmt.Sprintf("[0:v][1:v] overlay=(W-w)/2:(H-h)/2:enable='between(t,0,%d)'", int(titleDuration))},
	).Run()

	if err != nil {
		return vid, fmt.Errorf("Error overlaying title image onto video: %v", err)
	}

	fmt.Println("Video with overlay created successfully")

	fmt.Println("All done! Your video is ready at " + completeFilePath)
	fmt.Println("output path: " + completeFilePath)

	vid = SlopVideo{
		Title:      thread.Title,
		Transcript: thread.ToString(),
		Path:       completeFilePath,
		AudioPath:  fullAudioPath,
	}

	return vid, nil
}
