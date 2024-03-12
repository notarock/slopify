package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/notarock/slopify/pkg/chatgpt"
	"github.com/notarock/slopify/pkg/config"
	"github.com/notarock/slopify/pkg/google"
	"github.com/notarock/slopify/pkg/reddit"
	"github.com/notarock/slopify/pkg/subs"
	"github.com/spf13/cobra"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const SUBTITLE_COMMAND = "subtitles=%s:force_style='FontSize=24,Alignment=10'"
const BUCKET_NAME = "slopify-transcription-buffer"
const SUBTITLES_FILE = "output.srt"

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
			log.Fatalf("Error generating slop from reddit thread: %v", err)
		}
	},
}

func redditVideo(cfg config.Config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("No URL or file were provided. Please provide a URL or file containing the HTML of the thread.")
	}

	err := createWorkspace(workingDir)
	if err != nil {
		return fmt.Errorf("Error creating workspace: %v", err)
	}

	arg := args[0]

	_, err = url.ParseRequestURI(arg)

	// If it's a URL, scrape it
	// If it's a file, scrape it
	var thread reddit.Thread
	if err == nil {
		thread = reddit.ScrapeThread(arg)
	} else {
		err = nil
		thread = reddit.ScrapeFromFile(arg)
	}

	if thread.Title == "" {
		return fmt.Errorf("No title found in the thread, please provide a valid URL or file containing the HTML of the thread.")
	}

	titlePath := workingDir + "/audio/title.mp3"
	google.GetVoiceFile(thread.Title, titlePath)

	// TODO: Function to get the top comments and limit depth
	files := []string{}
	for _, value := range thread.CommentThreads {
		if len(files) > 0 {
			break
		}
		for _, comment := range value.Comments {
			if len(files) > 0 {
				break
			}

			filename := fmt.Sprintf(workingDir+"/"+"audio/%d.mp3", len(files))
			google.GetVoiceFile(comment, filename)
			files = append(files, filename)
		}
	}

	fullAudioPath := workingDir + "/output.mp3"

	google.Concatenate(workingDir+"/audio/title.mp3", files, fullAudioPath)

	if err != nil {
		return fmt.Errorf("Error writing thread to file: %v", err)
	}

	var videoFile string

	if footageDir != "" {
		path, err := pickFromDirectory(footageDir)
		if err != nil {
			return fmt.Errorf("Error picking a video from the directory: %v", err)
		}
		videoFile = path
	} else {
		videoFile = "source.webm"
	}

	// TODO: Randomize the start time
	// This works when using an hour long footage of subway surfer
	//
	// randomWithSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	// randomNumber := randomWithSeed.Intn(56)    // Generate a number between 0 and 55 (inclusive)
	// fmt.Println(randomNumber)
	// video := ffmpeg.Input(videoFile, ffmpeg.KwArgs{"ss": fmt.Sprintf("00:%d:55", randomNumber)}).Video()

	video := ffmpeg.Input(videoFile).Video()
	audio := ffmpeg.Input(fullAudioPath).Audio()

	outputFile := workingDir + "/" + "slop-" + strconv.FormatInt(time.Now().Unix(), 10)
	outputFileWithSubs := outputFile + "-subs.mp4"
	outputFile = outputFile + ".mp4"

	// Get the resolution of the video
	probeOutput, err := probeVideoResolution(videoFile)
	if err != nil {
		return fmt.Errorf("Error probing video resolution: %v", err)
	}

	height := probeOutput.Streams[0].Height
	width := probeOutput.Streams[0].Height / 16 * 9

	kwArgs := []ffmpeg.KwArgs{
		ffmpeg.KwArgs{"shortest": ""},
		ffmpeg.KwArgs{"vf": fmt.Sprintf("crop=%d:%d", width, height)},
	}

	out := ffmpeg.
		Output(
			[]*ffmpeg.Stream{video, audio},
			outputFile,
			kwArgs...,
		)

	err = out.Run()
	if err != nil {
		return fmt.Errorf("Error generating video: %v", err)
	}

	fmt.Println("Let's upload this to Google Cloud Storage")

	uri, err := google.UploadFile(BUCKET_NAME, outputFile)
	defer google.DeleteFile(BUCKET_NAME, filepath.Base(outputFile))
	if err != nil {
		return fmt.Errorf("Error uploading video to Google Cloud Storage: %v", err)
	}

	fmt.Println("Uploaded to " + uri)

	fmt.Println("Transcribing video...")
	transcript, err := google.SpeechTranscriptionURI(uri)
	if err != nil {
		return fmt.Errorf("Error transcribing video: %v", err)
	}

	fullTranscript := subs.BuildSubtitlesFromGoogle(transcript)

	fmt.Println("Converting to SRT format...")

	// Convert to SRT format
	srtData := subs.ConvertToSRT(fullTranscript)

	// Write to SRT file
	srtPath := workingDir + "/" + SUBTITLES_FILE
	err = subs.WriteSRT(srtPath, srtData)
	if err != nil {
		return fmt.Errorf("Error writing SRT file: %v", err)
	}

	fmt.Println("Adding subtitles to video...")

	err = ffmpeg.Input(outputFile).Output(outputFileWithSubs, ffmpeg.KwArgs{"vf": fmt.Sprintf(SUBTITLE_COMMAND, srtPath)}).Run()
	if err != nil {
		return fmt.Errorf("Error adding subtitles to video: %v", err)
	}

	fmt.Println("All done! Your video is ready at " + outputFileWithSubs)
	fmt.Println("output path: " + outputFileWithSubs)

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
	content := thread.ToString()
	res, err := client.PromptFromContent(content)

	log.Printf("Uploading video: %s\n", videoPath)
	log.Printf("Title: %s\n", res.Title)
	log.Printf("Description: %s\n", res.Description)
	uploadedVideo, err := uploader.UploadVideo(outputFileWithSubs, res.Title, res.Description, "private")

	if err != nil {
		return fmt.Errorf("Error uploading video: %v", err)
	}

	log.Printf("Video uploaded successfully: %s\n", uploadedVideo.Id)

	return nil
}
func askForConfirmation(message string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message + " (y/n): ")
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("Error reading input: %v", err)
	}

	response = strings.TrimSpace(response)
	return strings.ToLower(response) == "y", nil
}

func pickFromDirectory(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	randomWithSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := randomWithSeed.Intn(len(files)) // Generate a number between 0 and 55 (inclusive)

	return dir + "/" + files[randomNumber].Name(), nil
}

func createWorkspace(directory string) error {
	err := createDirectory(directory)
	if err != nil {
		return fmt.Errorf("Error creating workspace: %v", err)
	}
	err = createDirectory(directory + "/audio")
	if err != nil {
		return fmt.Errorf("Error creating audio workspace: %v", err)
	}
	err = createDirectory(directory + "/video")
	if err != nil {
		return fmt.Errorf("Error creating video workspace: %v", err)
	}
	return nil
}

func createDirectory(directory string) error {
	// Check if the directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.Mkdir(directory, 0755) // 0755 is the default permission
		if err != nil {
			return fmt.Errorf("Error creating directory: %v", err)
		}
		log.Println("Directory", directory, "created successfully.")
	} else if err != nil {
		// Some error occurred while checking directory existence
		return err
	}
	return nil
}

type FFProbeOutput struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

func probeVideoResolution(filepath string) (FFProbeOutput, error) {
	// Construct the ffprobe command to get video resolution
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "json", filepath)

	// Execute the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing ffprobe:", err)
		return FFProbeOutput{}, err
	}
	// Parse the output
	var probeOutput FFProbeOutput
	if err := json.Unmarshal(output, &probeOutput); err != nil {
		fmt.Println("Error parsing ffprobe output:", err)
		return FFProbeOutput{}, err
	}

	return probeOutput, nil
}
