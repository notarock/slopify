package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
