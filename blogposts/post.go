package blogposts

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Post struct {
	Title, Description string
	Tags []string
	Body string
}

const (
	titleSeparator       = "Title: "
	descriptionSeparator = "Description: "
	tagsSeparator = "Tags: "
)

func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)
	readMetaLine := func(tagName string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tagName)
	}

	titleLine := readMetaLine(titleSeparator)
	descriptionLine := readMetaLine(descriptionSeparator)
	tagsLine := readMetaLine(tagsSeparator)

	body := readBody(scanner)

	tags := strings.Split(tagsLine, ", ")
	post := Post{
		Title:       titleLine,
		Description: descriptionLine,
		Tags:        tags,
		Body:        body,
	}
	return post, nil
}

func readBody(scanner *bufio.Scanner) string {
	buff := bytes.Buffer{}
	scanner.Scan() //ignore a line
	for scanner.Scan() {
		fmt.Fprintln(&buff, scanner.Text())
	}
	body := strings.TrimSuffix(buff.String(), "\n")
	return body
}
