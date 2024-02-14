package cmd

import (
	"fmt"
	"os"

	"github.com/notarock/slopify/pkg/config"
	"github.com/spf13/cobra"
)

var cfg config.Config


func init() {
	var err error
	if	cfg, err = config.GetEnvConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(os.Args) == 1 {
		rootCmd.Help()
	}
}


var rootCmd = &cobra.Command{
	Use:   "slopify [old.reddit.com URL]",
	Short: "Slopify is a short-form content generator",
	Long: `Slopify is a short-form content generator that takes a reddit post and
	generates a short-form video with the post's content.`,
	Run: func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
