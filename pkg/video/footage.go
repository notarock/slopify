package video

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"time"
)

type FFProbeOutput struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
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

	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "json", filepath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return probeOutput, fmt.Errorf("error executing ffprobe: %v", err)

	}

	if err := json.Unmarshal(output, &probeOutput); err != nil {
		return probeOutput, fmt.Errorf("error parsing ffprobe output: %v", err)
	}

	return probeOutput, nil
}
