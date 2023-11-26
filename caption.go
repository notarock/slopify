package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/notarock/sludger/pkg/audio"
)

type Transcription struct {
	Results []Subtitle `json:"results"`
}

type Subtitle struct {
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Transcript string `json:"transcript"`
}

func main() {
	transcriptions, err := audio.SpeechTranscriptionURI("gs://sludger-temp/Snapinsta.app_video_404249692_2331325873721638_6752152422268909354_n.mp4")
	if err != nil {
		panic(err)
	}

	var transcription Transcription

	for _, t := range transcriptions {
		alternative := t.GetAlternatives()[0]

		fmt.Printf("Word level information:\n")
		for _, wordInfo := range alternative.GetWords() {
			startTime := wordInfo.GetStartTime()
			endTime := wordInfo.GetEndTime()
			transcription.Results = append(transcription.Results, Subtitle{
				StartTime:  fmt.Sprintf("%4.1f", float64(startTime.GetSeconds())+float64(startTime.GetNanos())*1e-9),
				EndTime:    fmt.Sprintf("%4.1f", float64(endTime.GetSeconds())+float64(endTime.GetNanos())*1e-9),
				Transcript: wordInfo.GetWord(),
			})
		}
	}

	fmt.Printf("%+v\n", transcription)

	// Convert to SRT format
	srtData := convertToSRT(transcription)

	// Write to SRT file
	err = ioutil.WriteFile("output.srt", []byte(srtData), 0644)
	if err != nil {
		panic(err)
	}
}

func convertToSRT(transcription Transcription) string {
	var srtBuilder strings.Builder

	for i, item := range transcription.Results {
		startTime := formatTimestamp(item.StartTime)
		endTime := formatTimestamp(item.EndTime)
		srtBuilder.WriteString(fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, startTime, endTime, item.Transcript))
	}

	return srtBuilder.String()
}

func formatTimestamp(timestamp string) string {
	// Assuming timestamp is in seconds and as a string, e.g., "1.0" or "3.5"

	seconds, err := strconv.ParseFloat(strings.Trim(timestamp, " "), 64)
	if err != nil {
		panic(err)
	}
	t := time.Duration(seconds * float64(time.Second))
	return fmt.Sprintf("%02d:%02d:%02d,%03d", int(t.Hours()), int(t.Minutes())%60, int(t.Seconds())%60, (t.Milliseconds())%1000)
}
