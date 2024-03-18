package cmd

import (
	"fmt"
	"log"

	"github.com/notarock/slopify/pkg/chatgpt"
	"github.com/notarock/slopify/pkg/config"
	"github.com/notarock/slopify/pkg/google"
	"github.com/notarock/slopify/pkg/video"
	"github.com/spf13/cobra"
)

var footageDir string
var workingDir string

func init() {
	rootCmd.PersistentFlags().StringVar(&footageDir, "footage", "", "Directory to search for footage. If not provided, the default footage will be used.")
	rootCmd.PersistentFlags().StringVar(&tokenFile, "tokenFile", "token.json", "Path to your token file where you want to store your YouTube token. If not provided, the default token file will be used.")
	rootCmd.PersistentFlags().StringVar(&oauthConfigFile, "oauthConfigFile", "youtubeConfig.json", "Path to your 0AUTH client secret file.")

	redditCmd.PersistentFlags().StringVar(&workingDir, "workingDir", "/tmp/slopify", "Directory to store temporary files. If not provided, the default working directory will be used.")

	rootCmd.AddCommand(redditCmd)
}

var redditCmd = &cobra.Command{
	Use:   "reddit https://old.reddit.com/r/...",
	Short: "Generate a video from a reddit thread",
	Long:  `Generate a video from a reddit URL or a file containing the HTML of the thread.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := redditVideo(cfg, args); err != nil {
			log.Fatal(err)
		}
	},
}

func redditVideo(cfg config.Config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("No URL or file were provided. Please provide a URL or file containing the HTML of the thread.")
	}
	arg := args[0]

	input := video.RedditThreadVideoInput{
		ThreadURL:        arg,
		WorkingDirectory: workingDir,
		FootageDirectory: footageDir,
	}

	slop, err := video.BuildFromThread(input)
	if err != nil {
		return fmt.Errorf("Error building slop from thread: %v", err)
	}

	// Ask if the user wants to upload the video to YouTube
	upload, err := askForConfirmation("Do you want to upload the video to YouTube?")
	if err != nil {
		return fmt.Errorf("Error asking for confirmation: %v", err)
	}

	if !upload {
		fmt.Println("Video not uploaded to YouTube.")
		return nil
	}

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

	fmt.Println("Prompting for title, description and tags")

	client := chatgpt.NewGPTClient(cfg.OpenaiKey)
	res, err := client.PromptFromContent(slop.Transcript)

	log.Printf("Uploading video: %s\n", slop.Path)
	log.Printf("Title: %s\n", res.Title)
	log.Printf("Description: %s\n", res.Description)
	uploadedVideo, err := uploader.UploadVideo(slop.Path, res.Title, res.Description, "private")

	if err != nil {
		return fmt.Errorf("Error uploading video: %v", err)
	}

	log.Printf("Video uploaded successfully: %s\n", uploadedVideo.Id)

	return nil
}
