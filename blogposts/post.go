package blogposts

import (
	"bufio"
	"io"
)

type Post struct {
	Title, Description string
}

const (
	titleSeparator       = "Title: "
	descriptionSeparator = "Description: "
)

func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)
	readLine := func() string {
		scanner.Scan()
		return scanner.Text()
	}

	titleLine := readLine()
	descriptionLine := readLine()

	post := Post{
		Title:       titleLine,
		Description: descriptionLine,
	}
	return post, nil
}
