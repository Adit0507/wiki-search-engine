package models

import (
	"strings"

	"github.com/Adit0507/wiki-search-engine/internal/utils"
)

type Document struct {
	ID      uint32         `json:"id"`
	Title   string         `json:"title"`
	Content string         `json:"content"`
	URL     string         `json:"url"`
	Terms   map[string]int `json:"terms"`
	Length  int            `json:"length"`
}

func NewDocument(id uint32, title, content, url string) *Document {
	doc := &Document{
		ID:      id,
		Title:   title,
		Content: content,
		URL:     url,
		Terms:   make(map[string]int),
		Length:  0,
	}

	doc.processText()

	return doc
}

func (d *Document) processText() {
	// combin title & content for indexig
	text := strings.ToLower(d.Title + " " + d.Title + " " + d.Content)

	tokens := utils.Tokenize(text)
	for _, token := range tokens {
		stemmed := utils.Stem(token)

		if len(stemmed) > 2 { //filtering out very short terms
			d.Terms[stemmed]++
			d.Length++
		}
	}
}

func (d *Document) GetTermFreq(term string) int {
	return d.Terms[term]
}

func (d *Document) GetLength() int {
	return d.Length
}
