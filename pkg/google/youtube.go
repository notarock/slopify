package google

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

type YouTubeUploader struct {
	Service *youtube.Service
}

func NewYouTubeUploader(clientSecretFile string, tokenFile string) (*YouTubeUploader, error) {
	ctx := context.Background()
	b, err := os.ReadFile(clientSecretFile)
	if err != nil {
		return nil, err
	}

	// If modifying these scopes, delete the token.json file.
	config, err := google.ConfigFromJSON(
		b,
		youtube.YoutubeUploadScope,
		youtube.YoutubeReadonlyScope,
		youtube.YoutubeForceSslScope,
	)
	if err != nil {
		return nil, err
	}
	client := getClient(ctx, config, tokenFile)

	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}

	return &YouTubeUploader{Service: service}, nil
}

func getClient(ctx context.Context, config *oauth2.Config, tokenFile string) *http.Client {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(ctx, tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func (uploader *YouTubeUploader) UploadVideo(videoPath string, title string, description string, privacyStatus string) (*youtube.Video, error) {
	file, err := os.Open(videoPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	video := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: description,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: privacyStatus,
		},
	}

	call := uploader.Service.Videos.Insert([]string{"snippet", "status"}, video)
	response, err := call.Media(file).Do()
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetUserInfo retrieves information about the authenticated user's YouTube channel.
func (uploader *YouTubeUploader) GetUserInfo() (*youtube.ChannelListResponse, error) {
	channelsListCall := uploader.Service.Channels.List([]string{"snippet"}).Mine(true)
	response, err := channelsListCall.Do()
	if err != nil {
		return nil, err
	}
	return response, nil
}
