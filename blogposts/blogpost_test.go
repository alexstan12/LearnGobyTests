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
Description: Description 1
Tags: tdd, go
---
Hello
World`
		secondBody = `Title: Post 2
Description: Description 2
Tags: rust, borrow-checker
---
B
L
M`
	)

	//Given
	fs := fstest.MapFS{
		"hello-world.md":  {Data: []byte(firstBody)},
		"hello world2.md": {Data: []byte(secondBody)},
	}
	//When
	posts, err := blogposts.PostFromFS(fs)
	t.Logf("the posts are %s", posts)
	if err != nil {
		t.Fatal(err)
	}

	//Then
	if len(posts) != len(fs) {
		t.Errorf("expected %d posts, got %d posts", len(fs), len(posts))
	}

	expectedPost := blogposts.Post{
		Title: "Post 1",
		Description: "Description 1",
		Tags:[]string{"tdd", "go"},
		Body: `Hello
World`,
	}
	assertPost(t, posts[1], expectedPost)
}

func assertPost(t *testing.T, got, want blogposts.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
