package google

import (
	"context"
	"fmt"

	video "cloud.google.com/go/videointelligence/apiv1"
	videopb "cloud.google.com/go/videointelligence/apiv1/videointelligencepb"
)

func SpeechTranscriptionURI(file string) ([]*videopb.SpeechTranscription, error) {
	ctx := context.Background()

	client, err := video.NewClient(ctx)
	if err != nil {
		return []*videopb.SpeechTranscription{}, err
	}
	defer client.Close()

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		Features: []videopb.Feature{
			videopb.Feature_SPEECH_TRANSCRIPTION,
		},
		VideoContext: &videopb.VideoContext{
			SpeechTranscriptionConfig: &videopb.SpeechTranscriptionConfig{
				LanguageCode:               "en-US",
				EnableAutomaticPunctuation: true,
			},
		},
		InputUri: file,
	})
	if err != nil {
		return []*videopb.SpeechTranscription{}, err
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return []*videopb.SpeechTranscription{}, err
	}

	fmt.Printf("ici %+v\n", resp)

	// A single video was processed. Get the first result.
	result := resp.AnnotationResults[0]

	fmt.Printf("%+v\n", result)

	for _, transcription := range result.SpeechTranscriptions {
		// The number of alternatives for each transcription is limited by
		// SpeechTranscriptionConfig.MaxAlternatives.
		// Each alternative is a different possible transcription
		// and has its own confidence score.
		for _, alternative := range transcription.GetAlternatives() {
			fmt.Printf("Alternative level information:\n")
			fmt.Printf("\tTranscript: %v\n", alternative.GetTranscript())
			fmt.Printf("\tConfidence: %v\n", alternative.GetConfidence())

			fmt.Printf("Word level information:\n")
			for _, wordInfo := range alternative.GetWords() {
				startTime := wordInfo.GetStartTime()
				endTime := wordInfo.GetEndTime()
				fmt.Printf("\t%4.1f - %4.1f: %v (speaker %v)\n",
					float64(startTime.GetSeconds())+float64(startTime.GetNanos())*1e-9, // start as seconds
					float64(endTime.GetSeconds())+float64(endTime.GetNanos())*1e-9,     // end as seconds
					wordInfo.GetWord(),
					wordInfo.GetSpeakerTag())
			}
		}
	}

	return result.SpeechTranscriptions, nil
}
