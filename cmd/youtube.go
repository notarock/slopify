package cmd

import (
	"fmt"
	"log"

	"github.com/notarock/slopify/pkg/config"
	"github.com/notarock/slopify/pkg/google"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(youtubeCmd)
}

var youtubeCmd = &cobra.Command{
	Use:   "youtube",
	Short: "View and manage your slops.",
	Long:  `View and manage your slops on YouTube`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := doYoutube(cfg, args); err != nil {
			log.Fatalf("Error generating slop from reddit thread: %v", err)
		}
	},
}

func doYoutube(cfg config.Config, args []string) error {
	clientSecretFile := "youtubeConfig.json" // Path to your client secret file
	tokenFile := "token.json"                // Path to your token file

	uploader, err := google.NewYouTubeUploader(clientSecretFile, tokenFile)
	if err != nil {
		return fmt.Errorf("Error creating YouTube uploader: %v", err)
	}

	response, err := uploader.GetUserInfo()
	if err != nil {
		return fmt.Errorf("Error creating YouTube uploader: %v", err)
	}
	for _, item := range response.Items {
		fmt.Printf("Channel ID: %s\n", item.Id)
		fmt.Printf("Title: %s\n", item.Snippet.Title)
		fmt.Printf("Description: %s\n", item.Snippet.Description)
		fmt.Printf("Custom URL: %s\n", item.Snippet.CustomUrl)
	}

	videoPath := "slop-1708318402-subs.mp4" // Path to the video you want to upload
	title := "My Uploaded Video"
	description := "Description of my video"
	privacyStatus := "private" // or "private" or "unlisted"

	video, err := uploader.UploadVideo(videoPath, title, description, privacyStatus)
	if err != nil {
		return fmt.Errorf("Error uploading video: %v", err)
	}

	fmt.Printf("Video uploaded successfully: %s\n", video.Id)
	return nil
}
