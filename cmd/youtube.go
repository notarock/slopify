package cmd

import (
	"fmt"
	"log"

	"github.com/notarock/slopify/pkg/config"
	"github.com/notarock/slopify/pkg/google"
	"github.com/spf13/cobra"
)

var tokenFile string
var oauthConfigFile string

var videoPath string
var title string
var description string
var privacyStatus string

func init() {
	youtubeCmd.PersistentFlags().StringVar(&tokenFile, "tokenFile", "token.json", "Path to your token file where you want to store your YouTube token. If not provided, the default token file will be used.")
	youtubeCmd.PersistentFlags().StringVar(&oauthConfigFile, "oauthConfigFile", "youtubeConfig.json", "Path to your 0AUTH client secret file.")

	youtubeCmd.PersistentFlags().StringVar(&videoPath, "videoPath", "", "Path to the video you want to upload")
	youtubeCmd.PersistentFlags().StringVar(&title, "title", "My Uploaded Video", "Title of the video")
	youtubeCmd.PersistentFlags().StringVar(&description, "description", "Description of my video", "Description of the video")
	youtubeCmd.PersistentFlags().StringVar(&privacyStatus, "privacyStatus", "private", "Privacy status of the video")

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
	uploader, err := google.NewYouTubeUploader(oauthConfigFile, tokenFile)
	if err != nil {
		return fmt.Errorf("Error creating YouTube uploader: %v", err)
	}

	response, err := uploader.GetUserInfo()
	if err != nil {
		return fmt.Errorf("Error getting user info: %v", err)
	}

	for _, item := range response.Items {
		log.Printf("This is the info for your YouTube account: \n")
		log.Printf("Channel ID: %s\n", item.Id)
		log.Printf("Title: %s\n", item.Snippet.Title)
		log.Printf("Description: %s\n", item.Snippet.Description)
		log.Printf("Custom URL: %s\n", item.Snippet.CustomUrl)
	}

	if videoPath == "" {
		return fmt.Errorf("No video path provided. Please provide a valid video file path.")
	}

	log.Printf("Uploading video: %s\n", videoPath)
	video, err := uploader.UploadVideo(videoPath, title, description, privacyStatus)
	if err != nil {
		return fmt.Errorf("Error uploading video: %v", err)
	}

	log.Printf("Video uploaded successfully: %s\n", video.Id)
	return nil
}
