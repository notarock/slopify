package reddit

import (
	"encoding/json"
	"os"
)

type Thread struct {
	Title          string
	Url            string
	CommentThreads []Comment
}

type Comment struct {
	Comments []string
}

func (t Thread) ToString() string {
	output := t.Title + "\n"
	for _, comment := range t.CommentThreads {
		for _, c := range comment.Comments {
			output += c + "\n"
		}
	}
	return output
}

func (t Thread) TotalComments() int {
	total := 0
	for _, comment := range t.CommentThreads {
		total += len(comment.Comments)
	}
	return total
}

func (t Thread) DumpToFile(filename string) error {
	os.Remove(filename)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(t)
	if err != nil {
		return err
	}
	return nil
}
