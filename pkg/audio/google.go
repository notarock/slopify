// Command quickstart generates an audio file with the content "Hello, World!".
package audio

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	video "cloud.google.com/go/videointelligence/apiv1"
	videopb "cloud.google.com/go/videointelligence/apiv1/videointelligencepb"
)

func GetVoiceFile(toSay, outputPath string) {
	// Instantiates a client.
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: toSay},
		},
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			Name:         "en-US-Wavenet-F", // This is one of the WaveNet female voices for en-US
			SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(outputPath, resp.AudioContent, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Audio content written to file: %v\n", outputPath)
}

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
