package blogpost_test

import (
	blogpost "github.com/alexstan12/LearnGobyTests/reading-files-demo"
	"testing"
	"testing/fstest"
)

func TestBlogPosts(t *testing.T) {
	//'Given
	fs := fstest.MapFS{
		"hello-world.md":  {Data: []byte("hello, world")},
		"hello-twitch.md": {Data: []byte("hello, twitch")},
	}
	//When
	posts := blogpost.PostFromFS(fs)

	//Then
	if len(posts) != len(fs) {
		t.Errorf("expected %d posts, got %d posts", len(fs), len(posts))
	}
}
