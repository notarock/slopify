package reddit

import (
	"fmt"
	"io"
	"log"
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

func ScrapeThread(threadURI string) Thread {
	res, err := http.Get(threadURI)
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

	ignoreNext := 0

	commentArea.Find(".thing").Each(func(i int, s *goquery.Selection) {
		// If ignoring the next comment, skip it
		if ignoreNext != 0 {
			ignoreNext = ignoreNext - 1
			return
		}

		comments := s.Find(".md")
		Comment := Comment{}
		ignoreNext = comments.Length() - 1

		fmt.Println("Ignoring comments: ", ignoreNext)

		comments.Each(func(i int, s *goquery.Selection) {
			commentText := s.Text()
			cleanedText := strings.Replace(commentText, "\n", "", -1)
			Comment.Comments = append(Comment.Comments, cleanedText)
		})

		thread.CommentThreads = append(thread.CommentThreads, Comment)
	})

	// doc.Find(".details").Each(func(i int, s *goquery.Selection) {
	// 	// For each item found, get the title, link, tags and author.
	// 	var taglist []string
	// 	if tags, ok := s.Find(".tag").Attr("title"); ok {
	// 		taglist = strings.Split(tags, TAGS_SEPARATOR)
	// 	}

	// 	article := articles.Article{
	// 		ID:     articles.LinkToID(link),
	// 		Title:  title,
	// 		Link:   link,
	// 		Tags:   taglist,
	// 		Author: author,
	// 		Source: SOURCE_NAME,
	// 	}

	// 	lobsterArticles = append(lobsterArticles, article)
	// })

	return thread

}
