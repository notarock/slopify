package main

import (
	"fmt"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func main() {
	// ffmpeg.Input("source.mp4").Output("slop.mp4").Run()
	videoFile := "source.mp4"
	audioFile := "output.mp3"
	outputFile := "output_with_audio-code.mp4"

	ffmpeg.Input(videoFile)

	video := ffmpeg.Input(videoFile).Video()
	audio := ffmpeg.Input(audioFile).Audio()

	out := ffmpeg.Output([]*ffmpeg.Stream{video, audio}, outputFile)

	out.Run()

	fmt.Println(out.GetArgs())

}
