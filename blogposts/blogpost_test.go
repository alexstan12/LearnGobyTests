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

func (s StubFailingFS) Open(name string) (fs.File, error){
	return nil, errors.New("oh no I'm always failing")
}

func TestBlogPosts(t *testing.T) {
	//Given
	fs := fstest.MapFS{
		"hello-world.md":  {Data: []byte("Title: Post 1")},
		"hello-twitch.md": {Data: []byte("Title: Post 2")},
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

	expectedPost := blogposts.Post{Title: "Post 1"}
	assertPost(t, expectedPost, posts[0])
}

func assertPost(t *testing.T, got, want blogposts.Post) {
	if reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
