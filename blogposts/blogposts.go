package blogposts

import (
	"io"
	"io/fs"
	"strings"
)

type Post struct {
	Title, Description, Body string
	Tags []string
}

func PostFromFS(filesystem fs.FS) ([]Post, error) {
	dir, err := fs.ReadDir(filesystem, ".")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for _,f := range dir {
		posts = append(posts, makePostFromFile(filesystem, f))
	}
	return posts, nil
}

func makePostFromFile(filesystem fs.FS, f fs.DirEntry) Post {
	blogFile, _ := filesystem.Open(f.Name())
	fileContents, _ := io.ReadAll(blogFile)
	title := strings.TrimPrefix(string(fileContents),  "Title: ")
	return Post{
		Title: title,
	}
}


