package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use:   "slopify [old.reddit.com URL]",
	Short: "Slopify is a short-form content generator",
	Long: `Slopify is a short-form content generator that takes a reddit post and
	generates a short-form video with the post's content.`,
	Run: func(cmd *cobra.Command, args []string) {
		full()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
