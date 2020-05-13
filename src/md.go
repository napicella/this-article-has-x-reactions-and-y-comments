package main

import (
	"github.com/pkg/errors"
	"strings"
)

const frontMatterTitle = "title:"

type simpleFrontMatterEditor struct {
	markdown   string
	startIndex int
	endIndex   int
	title      string
}

func newFrontMatterEditor(markdown string) (*simpleFrontMatterEditor, error) {
	startIndex := strings.Index(markdown, frontMatterTitle)
	if startIndex == -1 {
		return nil, errors.New("Unable to find the article title")
	}
	var endIndex = 0
	for endIndex = startIndex; endIndex < len(markdown); endIndex++ {
		if markdown[endIndex] == '\n' {
			break
		}
	}
	if endIndex == len(markdown) {
		return nil, errors.New("Unable to find the article title")
	}
	return &simpleFrontMatterEditor{
		markdown:   markdown,
		title:      markdown[startIndex+len(frontMatterTitle) : endIndex],
		startIndex: startIndex,
		endIndex:   endIndex,
	}, nil
}

func (t *simpleFrontMatterEditor) shouldUpdateTitle(newTitle string) bool {
	return t.title != newTitle
}

func (t *simpleFrontMatterEditor) updateTitle(newTitle string) {
	if t.shouldUpdateTitle(newTitle) {
		t.title = newTitle
		t.markdown = t.markdown[:t.startIndex+len(frontMatterTitle)] +
			newTitle +
			t.markdown[t.endIndex:]
	}
}
