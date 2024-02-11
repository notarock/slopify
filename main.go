package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/notarock/slopify/pkg/google"
	"github.com/notarock/slopify/pkg/reddit"
	"github.com/notarock/slopify/pkg/subs"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const SUBTITLE_COMMAND = "subtitles=output.srt:force_style='FontSize=24,Alignment=10'"
const BUCKET_NAME = "sludger-temp"
const SUBTITLES_FILE = "output.srt"

func main() {
	argsWithProg := os.Args
	if len(argsWithProg) < 2 {
		fmt.Println("Please provide a URL or file to scrape")
		os.Exit(1)
	}

	arg := argsWithProg[1]

	_, err := url.ParseRequestURI(arg)
	// If it's a URL, scrape it
	var thread reddit.Thread
	if err == nil {
		thread = reddit.ScrapeThread(arg)
	} else {
		thread = reddit.ScrapeFromFile(arg)
	}

	fmt.Printf("%+v\n", thread)
	fmt.Println(thread.TotalComments())

	if thread.Title == "" {
		fmt.Println("No title found ??")
		os.Exit(1)
	}

	fmt.Printf("Processing audio...")

	google.GetVoiceFile(thread.Title, "audio/title.mp3")

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

	fmt.Printf("Processing audio...")
	google.Concatenate("audio/title.mp3", files, "output.mp3")

	fmt.Println("Audio saved to output.mp3")

	videoFile := "source.webm"
	audioFile := "output.mp3"

	// Nécessaire?
	// ffmpeg.Input(videoFile)

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	randomNumber := rand.Intn(56)    // Generate a number between 0 and 55 (inclusive)
	fmt.Println(randomNumber)

	video := ffmpeg.Input(videoFile, ffmpeg.KwArgs{"ss": fmt.Sprintf("00:%d:55", randomNumber)}).Video()
	audio := ffmpeg.Input(audioFile).Audio()

	outputFile := "slop-" + strconv.FormatInt(time.Now().Unix(), 10)
	outputFileWithSubs := outputFile + "-subs.mp4"
	outputFile = outputFile + ".mp4"

	args := []ffmpeg.KwArgs{
		ffmpeg.KwArgs{"shortest": ""},
	}

	out := ffmpeg.
		Output(
			[]*ffmpeg.Stream{video, audio},
			outputFile,
			args...,
		)

	out.Run()

	fmt.Println(out.GetArgs())

	fmt.Println("Let's upload this to Google Cloud Storage")

	uri, err := google.UploadFile(BUCKET_NAME, outputFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("Uploaded to " + uri)
	defer google.DeleteFile(BUCKET_NAME, outputFile)

	transcript, err := google.SpeechTranscriptionURI(uri)
	if err != nil {
		panic(err)
	}

	fullTranscript := subs.BuildSubtitlesFromGoogle(transcript)

	// Convert to SRT format
	srtData := subs.ConvertToSRT(fullTranscript)

	// Write to SRT file
	err = subs.WriteSRT(SUBTITLES_FILE, srtData)
	if err != nil {
		panic(err)
	}

	err = ffmpeg.Input(outputFile).Output(outputFileWithSubs, ffmpeg.KwArgs{"vf": SUBTITLE_COMMAND}).Run()
	if err != nil {
        log.Fatal(err)
	}

	fmt.Println("All done! Your video is ready at " + outputFileWithSubs)

}
