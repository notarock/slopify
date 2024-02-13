package cmd

import (
	"fmt"
	"log"

	"github.com/notarock/slopify/pkg/config"
	"github.com/notarock/slopify/pkg/pexels"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(footageCmd)
}

var footageCmd = &cobra.Command{
	Use:   "footage [search query]",
	Short: "Search for footage on Pexels",
	Long: `Search for footage on Pexels and download it to use in your video.`,
	Run: func(cmd *cobra.Command, args []string) {
		if	err := footage(cfg,args); err != nil {
			log.Fatalf("Error searching for footage: %v", err)
		}
	},
}

func footage(cfg config.Config, args []string) (error) {
	fmt.Println("Searching for footage on Pexels")
	searchQuery := args[0]

	client := pexels.NewPexelsClient(cfg.PexelsAPIKey)
	outputDirectory := "./footage"
	path, err := client.DownloadRandomVideo(searchQuery, outputDirectory )
	fmt.Println("Downloaded footage to: ", path)
	return err
}
