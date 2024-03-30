package video

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"time"
)

// Updated FFProbeOutput struct to include the Duration field.
type FFProbeOutput struct {
	Streams []struct {
		Width    int     `json:"width"`
		Height   int     `json:"height"`
		Duration float64 `json:"duration,string"` // Use string tag to handle both numeric and string durations.
	} `json:"streams"`
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

func pickFootageFromDirectory(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	randomWithSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := randomWithSeed.Intn(len(files))

	return dir + "/" + files[randomNumber].Name(), nil
}

func probeVideoResolution(filepath string) (FFProbeOutput, error) {
	var probeOutput FFProbeOutput

	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "format=duration:stream=width,height", "-of", "json", filepath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return probeOutput, fmt.Errorf("error executing ffprobe: %v", err)

	}

	if err := json.Unmarshal(output, &probeOutput); err != nil {
		return probeOutput, fmt.Errorf("error parsing ffprobe output: %v", err)
	}

	return probeOutput, nil
}
