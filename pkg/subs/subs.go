package subs

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	videopb "cloud.google.com/go/videointelligence/apiv1/videointelligencepb"
)

type Transcription struct {
	Results []Subtitle `json:"results"`
}

type Subtitle struct {
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Transcript string `json:"transcript"`
}

func BuildSubtitlesFromGoogle(transcript []*videopb.SpeechTranscription) Transcription {
	var transcription Transcription

	for _, t := range transcript {
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
	return transcription
}

func ConvertToSRT(transcription Transcription) string {
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

func WriteSRT(file, srtData string) error {
	return ioutil.WriteFile(file, []byte(srtData), 0644)
}
