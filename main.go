package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/notarock/sludger/pkg/audio"
	"github.com/notarock/sludger/pkg/reddit"
)

MERGE_COMMAND = "ffmpeg -i %s -i %s -c copy -map 0:v  -map 1:a %s"

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

}
