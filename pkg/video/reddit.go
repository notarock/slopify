package video

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
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

	google.Concatenate(input.WorkingDirectory+"/audio/title.mp3", files, fullAudioPath)

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
	//
	// randomWithSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	// randomNumber := randomWithSeed.Intn(56)    // Generate a number between 0 and 55 (inclusive)
	// fmt.Println(randomNumber)
	// video := ffmpeg.Input(videoFile, ffmpeg.KwArgs{"ss": fmt.Sprintf("00:%d:55", randomNumber)}).Video()

	video := ffmpeg.Input(videoFile).Video()
	audio := ffmpeg.Input(fullAudioPath).Audio()

	outputFile := input.WorkingDirectory + "/" + "slop-" + strconv.FormatInt(time.Now().Unix(), 10)
	outputFileWithSubs := outputFile + "-subs.mp4"
	outputFile = outputFile + ".mp4"

	// Get the resolution of the video
	probeOutput, err := probeVideoResolution(videoFile)
	if err != nil {
		return vid, fmt.Errorf("Error probing video resolution: %v", err)
	}

	height := probeOutput.Streams[0].Height
	width := probeOutput.Streams[0].Height / 16 * 9

	kwArgs := []ffmpeg.KwArgs{
		ffmpeg.KwArgs{"shortest": ""},
		ffmpeg.KwArgs{"vf": fmt.Sprintf("crop=%d:%d", width, height)},
	}

	out := ffmpeg.
		Output(
			[]*ffmpeg.Stream{video, audio},
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

	fullTranscript := subs.BuildSubtitlesFromGoogle(transcript)

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

	err = ffmpeg.Input(outputFile).Output(outputFileWithSubs, ffmpeg.KwArgs{"vf": fmt.Sprintf(SUBTITLE_COMMAND, srtPath)}).Run()
	if err != nil {
		return vid, fmt.Errorf("Error adding subtitles to video: %v", err)
	}

	fmt.Println("All done! Your video is ready at " + outputFileWithSubs)
	fmt.Println("output path: " + outputFileWithSubs)

	vid = SlopVideo{
		Title:      thread.Title,
		Transcript: thread.ToString(),
		Path:       outputFileWithSubs,
		AudioPath:  fullAudioPath,
	}

	return vid, nil
}
