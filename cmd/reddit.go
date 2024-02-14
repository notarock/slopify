package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/notarock/slopify/pkg/config"
	"github.com/notarock/slopify/pkg/google"
	"github.com/notarock/slopify/pkg/reddit"
	"github.com/notarock/slopify/pkg/subs"
	"github.com/spf13/cobra"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const SUBTITLE_COMMAND = "subtitles=output.srt:force_style='FontSize=24,Alignment=10'"
const BUCKET_NAME = "sludger-temp"
const SUBTITLES_FILE = "output.srt"

func init() {
	rootCmd.AddCommand(redditCmd)
}

var redditCmd = &cobra.Command{
	Use:   "reddit https://old.reddit.com/r/...",
	Short: "Generate a video from a reddit thread",
	Long: `Generate a video from a reddit URL or a file containing the HTML of the thread.`,
	Run: func(cmd *cobra.Command, args []string) {
		if	err := redditVideo(cfg,args); err != nil {
			log.Fatalf("Error generating slop from reddit thread: %v", err)
		}
	},
}

func redditVideo(cfg config.Config, args []string) (error) {
	if len(args) < 1 {
		return fmt.Errorf("No URL or file were provided. Please provide a URL or file containing the HTML of the thread.")
	}

	arg := args[0]

	_, err := url.ParseRequestURI(arg)

	// If it's a URL, scrape it
	// If it's a file, scrape it
	var thread reddit.Thread
	if err == nil {
		thread = reddit.ScrapeThread(arg)
	} else {
		err = nil
		thread = reddit.ScrapeFromFile(arg)
	}

	// TODO: Debug
	// fmt.Printf("%+v\n", thread)
	// fmt.Println(thread.TotalComments())

	if thread.Title == "" {
		return fmt.Errorf("No title found in the thread, please provide a valid URL or file containing the HTML of the thread.")
	}

	// TODO: Debug
	// fmt.Printf("Processing audio...")

	google.GetVoiceFile(thread.Title, "audio/title.mp3")

	// TODO: Function to get the top comments and limit depth
	files := []string{}
	for _, value := range thread.CommentThreads {
		if len(files) > 10 {
			break
		}
		for _, comment := range value.Comments {
			if len(files) > 10 {
				break
			}

			filename := fmt.Sprintf("audio/%d.mp3", len(files))
			google.GetVoiceFile(comment, fmt.Sprintf("audio/%d.mp3", len(files)))
			files = append(files, filename)
		}
	}

	google.Concatenate("audio/title.mp3", files, "output.mp3")

	videoFile := "source.webm"
	audioFile := "output.mp3"

	randomWithSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := randomWithSeed.Intn(56)    // Generate a number between 0 and 55 (inclusive)
	fmt.Println(randomNumber)

	video := ffmpeg.Input(videoFile, ffmpeg.KwArgs{"ss": fmt.Sprintf("00:%d:55", randomNumber)}).Video()
	audio := ffmpeg.Input(audioFile).Audio()

	outputFile := "slop-" + strconv.FormatInt(time.Now().Unix(), 10)
	outputFileWithSubs := outputFile + "-subs.mp4"
	outputFile = outputFile + ".mp4"

	kwArgs := []ffmpeg.KwArgs{
		ffmpeg.KwArgs{"shortest": ""},
	}

	out := ffmpeg.
		Output(
			[]*ffmpeg.Stream{video, audio},
			outputFile,
			kwArgs...,
		)

	err = out.Run()
	if err != nil {
		return fmt.Errorf("Error generating video: %v", err)
	}

	fmt.Println("Let's upload this to Google Cloud Storage")

	uri, err := google.UploadFile(BUCKET_NAME, outputFile)
	if err != nil {
		return fmt.Errorf("Error uploading video to Google Cloud Storage: %v", err)
	}

	fmt.Println("Uploaded to " + uri)
	defer google.DeleteFile(BUCKET_NAME, outputFile)

	transcript, err := google.SpeechTranscriptionURI(uri)
	if err != nil {
		return fmt.Errorf("Error transcribing video: %v", err)
	}

	fullTranscript := subs.BuildSubtitlesFromGoogle(transcript)

	// Convert to SRT format
	srtData := subs.ConvertToSRT(fullTranscript)

	// Write to SRT file
	err = subs.WriteSRT(SUBTITLES_FILE, srtData)
	if err != nil {
		return fmt.Errorf("Error writing SRT file: %v", err)
	}

	err = ffmpeg.Input(outputFile).Output(outputFileWithSubs, ffmpeg.KwArgs{"vf": SUBTITLE_COMMAND}).Run()
	if err != nil {
		return fmt.Errorf("Error adding subtitles to video: %v", err)
	}

	fmt.Println("All done! Your video is ready at " + outputFileWithSubs)

	return nil
}
