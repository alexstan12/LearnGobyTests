package main

import (
blogposts "github.com/quii/fstest-spike"
"log"
"os"
)

func main() {
posts, err := blogposts.New(os.DirFS("posts"))
if err != nil {
log.Fatal(err)
}
log.Println(posts)
}