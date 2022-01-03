package blogpost

import (
	"io/fs"
	"testing/fstest"
)

type Post struct {
}

func PostFromFS(filesystem fstest.MapFS) []Post {
	dir, _ := fs.ReadDir(filesystem, ".")
	var posts []Post
	for range dir {
		posts = append(posts, Post{})
	}
	return posts
}
