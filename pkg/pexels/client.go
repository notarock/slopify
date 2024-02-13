package pexels

import (
	"context"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"strconv"
	"time"

	pexels "github.com/kosa3/pexels-go"
)

const SELECT_FROM_TOP_10 = 10

type PexelsClient struct {
	APIKey string
	api *pexels.Client
}

func NewPexelsClient(apiKey string) *PexelsClient {
	if apiKey == "" {
		fmt.Println("API Key is required")
		return nil
	}

	return &PexelsClient{
		APIKey: apiKey,
		api: pexels.NewClient(apiKey),
	}
}


func (c *PexelsClient) SearchVideos(query string) ([]*pexels.Video, error) {
	ps, err := c.api.VideoService.Search(context.TODO(), &pexels.VideoParams{
		Query: query,
	})
	if err != nil {
		return nil, fmt.Errorf("Error searching videos: %v", err)
	}

	videos := ps.Videos

	return videos, nil
}


func (c *PexelsClient) DownloadRandomVideo(query , resultDir string) (string, error) {
	resutls, err := c.SearchVideos(query)
	if err != nil {
		return "", fmt.Errorf("Error searching videos: %v", err)
	}

	if len(resutls) == 0 {
		return "", fmt.Errorf("No videos found for query: %s", query)
	}

	// Randomize the video selection
	rand.Seed(time.Now().UnixNano())
	// Select a random video from the top 10 results
	selection := rand.Intn(int(math.Min(float64(len(resutls)), SELECT_FROM_TOP_10)))
	link := resutls[selection].VideoFiles[0].Link
	name := strconv.FormatInt(time.Now().Unix(), 10)
	extension := GetVideoExtension(resutls[selection].VideoFiles[0].Link)

	resultPath := resultDir + "/" + name + extension


	fmt.Printf("Download url : %s", link)

	err = DownloadFile(resultPath, link)
	if err != nil {
		return "", fmt.Errorf("Error downloading video: %v", err)
	}

	return resultPath, nil
}

func GetVideoExtension(link string) string {
	extension := filepath.Ext(link)
    if idx := strings.Index(extension, "?"); idx != -1 {
        return extension[:idx]
    }
    return extension
}

// DownloadFile will download from a given url to a file. It will
// write as it downloads (useful for large files).
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
