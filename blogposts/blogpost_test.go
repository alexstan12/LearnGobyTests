package blogposts_test

import (
	"errors"
	blogposts "github.com/alexstan12/LearnGobyTests/blogposts"
	"io/fs"
	"testing"
	"testing/fstest"
)

type StubFailingFS struct {

}

func (s StubFailingFS) Open(name string) (fs.File, error){
	return nil, errors.New("oh no I'm always failing")
}

func TestBlogPosts(t *testing.T) {
	//Given
	fs := fstest.MapFS{
		"hello-world.md":  {Data: []byte("Title: Hello, TDD world!")},
		"hello-twitch.md": {Data: []byte("Title: Hello, twitchy world")},
	}
	//When
	posts, err := blogposts.PostFromFS(fs)

	if err != nil {
		t.Fatal(err)
	}

	//Then
	if len(posts) != len(fs) {
		t.Errorf("expected %d posts, got %d posts", len(fs), len(posts))
	}

	expectedPost := blogposts.Post{Title: "Hello, TDD world!"}
	if posts[0] != expectedPost {
		t.Errorf("got %v, want %v", posts[0], expectedPost)
	}
}
