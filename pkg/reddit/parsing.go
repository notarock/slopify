package reddit

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeFromFile(filename string) Thread {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	return ParseHtml(file, "")
}

var agents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Windows NT 1]]0.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.864.48 Safari/537.36 Edg/91.0.864.48",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36 OPR/77.0.4054.64",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
}

func ScrapeThread(threadURI string) Thread {
	reqURL := strings.Replace(threadURI, "www.reddit.com", "old.reddit.com", 1)
	randomIndex := rand.Intn(len(agents))
	// Pick an item randomly
	randomUserAgent := agents[randomIndex]

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", randomUserAgent)

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	fmt.Println(res.Body)

	return ParseHtml(res.Body, threadURI)
}

func ParseHtml(body io.ReadCloser, threadURI string) Thread {
	thread := Thread{
		Url: threadURI,
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	a := doc.Find(".title .may-blank")
	title := a.Text()

	thread.Title = title

	commentArea := doc.Find(".nestedlisting")

	s := commentArea.Find(".thing")
	// If ignoring the next comment, skip it

	Comment := Comment{}

	// Code dégueulasse
	// TODO: Refactor this
	// TODO: Add depth to the comment struct
	// TODO: Add a function to get the top comments and limit depth
	depth := 2
	comments := s.Find(".md")
	comments.Each(func(i int, s *goquery.Selection) {
		if depth == 0 {
			return
		}

		commentText := s.Text()
		cleanedText := strings.Replace(commentText, "\n", "", -1)
		Comment.Comments = append(Comment.Comments, cleanedText)
		depth--
	})

	thread.CommentThreads = append(thread.CommentThreads, Comment)

	return thread
}
