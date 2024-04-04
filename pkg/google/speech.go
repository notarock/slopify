package google

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Vernacular-ai/godub"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
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
			Name:         "en-US-Wavenet-J",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_MALE,
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
