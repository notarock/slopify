package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/notarock/sludger/pkg/audio"
	"github.com/notarock/sludger/pkg/reddit"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

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

	audio.GetVoiceFile(thread.Title, "audio/title.mp3")

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
			audio.GetVoiceFile(comment, fmt.Sprintf("audio/%d.mp3", len(files)))
			files = append(files, filename)
		}
	}

	fmt.Printf("Processing audio...")
	audio.Concatenate("audio/title.mp3", files, "output.mp3")

	fmt.Println("Audio saved to output.mp3")

	videoFile := "source.webm"
	audioFile := "output.mp3"

	ffmpeg.Input(videoFile)

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	randomNumber := rand.Intn(56)    // Generate a number between 0 and 55 (inclusive)
	fmt.Println(randomNumber)

	video := ffmpeg.Input(videoFile, ffmpeg.KwArgs{"ss": fmt.Sprintf("00:%d:55", randomNumber)}).Video()
	audio := ffmpeg.Input(audioFile).Audio()

	outputFile := "slop-" + strconv.FormatInt(time.Now().Unix(), 10) + ".mp4"

	args := []ffmpeg.KwArgs{
		ffmpeg.KwArgs{"shortest": ""},
	}

	crop() //#TODO: Fix this

	out := ffmpeg.
		Output(
			[]*ffmpeg.Stream{video, audio},
			outputFile,
			args...,
		)

	out.Run()

	fmt.Println(out.GetArgs())
}

func crop() {
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Crop video for phone? (y/n): ")
	// input, err := reader.ReadString('\n')

	// if err != nil {
	// 	fmt.Println("Error reading input:", err)
	// 	return
	// }

	// // Trim space and convert to lowercase for comparison
	// input = strings.TrimSpace(strings.ToLower(input))

	// if input == "y" {
	// 	fmt.Println("You chose YES!")
	// 	// Add any logic you want to execute for "y" here
	// 	args = append(args, ffmpeg.KwArgs{"aspect": "1080:1920"})

	// 	// When source is not 16:9, crop to 16:9
	// 	// args = append(args, ffmpeg.KwArgs{"vf": "crop=800:480:x:y"})
	// }
}
