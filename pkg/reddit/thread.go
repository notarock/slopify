package reddit

type Thread struct {
	Title          string
	Url            string
	CommentThreads []Comment
}

type Comment struct {
	Comments []string
}

func (t Thread) TotalComments() int {
	total := 0
	for _, comment := range t.CommentThreads {
		total += len(comment.Comments)
	}
	return total
}
