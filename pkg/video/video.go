package video

import (
	"fmt"
	"log"
	"os"
)

type SlopVideo struct {
	Title      string
	Transcript string
	Path       string
	AudioPath  string
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
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.Mkdir(directory, 0755) // 0755 is the default permission
		if err != nil {
			return fmt.Errorf("Error creating directory: %v", err)
		}
		log.Println("Directory", directory, "created successfully.")
	} else if err != nil {
		return err
	}
	return nil
}
