package blogposts_test

import (
	"errors"
	blogposts "github.com/alexstan12/LearnGobyTests/blogposts"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

type StubFailingFS struct {
}

func (s StubFailingFS) Open(name string) (fs.File, error) {
	return nil, errors.New("oh no I'm always failing")
}

func TestBlogPosts(t *testing.T) {
	const (
		firstBody = `Title: Post 1
Description: Description 1`
		secondBody = `Title: Post 2
Description: Description 2`
	)

	//Given
	fs := fstest.MapFS{
		"hello-world.md":  {Data: []byte(firstBody)},
		"hello-twitch.md": {Data: []byte(secondBody)},
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

	expectedPost := blogposts.Post{Title: "Post 1", Description: "Description 1"}
	assertPost(t, posts[0], expectedPost)
}

func assertPost(t *testing.T, got, want blogposts.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
