package audio

import (
	"fmt"
	"log"

	"github.com/Vernacular-ai/godub"
)

func Concatenate(title string, files []string, output string) {
	segment, _ := godub.NewLoader().Load(title)

	fmt.Print("Concatenating audio files...: ")
	fmt.Print(title)

	for _, file := range files {
		fmt.Print(", ", file)
		segmentToAdd, err := godub.NewLoader().Load(file)
		if err != nil {
			log.Fatal(err)
		}
		segment, err = segment.Append(segmentToAdd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Save the newly created audio segment as mp3 file.
	godub.NewExporter(output).WithDstFormat("mp3").WithBitRate(128000).Export(segment)
}
