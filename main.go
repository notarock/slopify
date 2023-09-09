package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/notarock/sludger/pkg/reddit"
)

func main() {
	argsWithProg := os.Args
	if len(argsWithProg) < 2 {
		fmt.Println("Please provide a URL or file to scrape")
		os.Exit(1)
	}

	arg := argsWithProg[1]

	_, err := url.ParseRequestURI(arg)
	// If it's a URL, scrape it
	var thread reddit.Thread
	if err == nil {
		thread = reddit.ScrapeThread(arg)
	} else {
		thread = reddit.ScrapeFromFile(arg)
	}

	fmt.Printf("%+v\n", thread)
	fmt.Println(thread.TotalComments())
}
