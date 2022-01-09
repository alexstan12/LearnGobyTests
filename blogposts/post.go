package blogposts

import (
	"io"
	"strings"
)

type Post struct {
	Title, Description, Body string
	Tags []string
}

func newPost(postFile io.Reader) (Post, error) {
	fileContents, err := io.ReadAll(postFile)
	if err != nil {
		return Post{}, err
	}
	title := strings.TrimPrefix(string(fileContents), "Title: ")
	post := Post{
		Title: title,
	}
	return post, nil
}
