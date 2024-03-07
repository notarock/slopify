package cmd

import (
	"fmt"
	"log"

	"github.com/notarock/slopify/pkg/chatgpt"
	"github.com/notarock/slopify/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(openaiCmd)
}

var openaiCmd = &cobra.Command{
	Use:   "openai",
	Short: "Send a test prompt to OpenAI's chatgpt",
	Long:  `Send a test prompt to OpenAI's chatgpt and get a response.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := doOpenai(cfg, args); err != nil {
			log.Fatalf("Error: %v", err)
		}
	},
}

func doOpenai(cfg config.Config, args []string) error {
	client := chatgpt.NewGPTClient(cfg.OpenaiKey)
	content := "This is a test prompt. Be creative! Create a story that will make people want to watch the video."
	res, err := client.PromptFromContent(content)
	fmt.Println(res)
	return err
}
