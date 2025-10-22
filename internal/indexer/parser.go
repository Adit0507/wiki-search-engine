package indexer

import (
	"compress/bzip2"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/Adit0507/wiki-search-engine/internal/models"
)

type Redirect struct {
	Title string `xml:"title,attr"`
}

type WikiPage struct {
	Title    string   `xml:"title"`
	ID       int64    `xml:"id"`
	Redirect Redirect `xml:"redirect"`
	Text     string   `xml:"revision>text"`
}

type Parser struct {
	docChan chan<- *models.Document
	docID   uint32
}

func NewParser(docChan chan<- *models.Document) *Parser {
	return &Parser{
		docChan: docChan,
		docID:   0,
	}
}

func (p *Parser) ParseFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var reader io.Reader = file

	// handling bzip2 compressed files
	if strings.HasSuffix(filename, ".bz2") {
		reader = bzip2.NewReader(file)
	}

	decoder := xml.NewDecoder(reader)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		
		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "page" {
				var page WikiPage

				if err := decoder.DecodeElement(&page, &se); err != nil {
					continue //skippin malformed pages
				}

				if p.shouldIndex(&page) {
					doc := p.createDocument(&page)
					if doc != nil {
						p.docChan <- doc
					}
				}
			}
		}
	}

	return nil
}

func (p *Parser) shouldIndex(page *WikiPage) bool {
	if page.Redirect.Title != "" { //skippin redirects
		return false
	}

	// skip special pages
	if strings.HasPrefix(page.Title, "File:") || strings.HasPrefix(page.Title, "Category:") || strings.HasPrefix(page.Title, "Template:") || strings.HasPrefix(page.Title, "Wikipedia:") || strings.HasPrefix(page.Title, "User:") || strings.HasPrefix(page.Title, "Talk:") {
		return false
	}

	// must've have content
	if len(strings.TrimSpace(page.Text)) < 100 {
		return false
	}

	return true
}

func (p *Parser) createDocument(page *WikiPage) *models.Document {
	p.docID++

	// clean the content
	content := p.cleanWikiText(page.Text)
	if len(content) < 50 {
		return nil
	}

	url := fmt.Sprintf("https://en.wikipedia.org/wiki/%s", strings.ReplaceAll(page.Title, " ", "_"))

	return models.NewDocument(p.docID, page.Title, content, url)
}

func (p *Parser) cleanWikiText(text string) string {
	re := regexp.MustCompile(`\{\{[^}]*\}\}`)
	text = re.ReplaceAllString(text, "")

	re = regexp.MustCompile(`\[\[[^|\]]*\|([^\]]*)\]\]`)
	text = re.ReplaceAllString(text, "$1")

	re = regexp.MustCompile(`\[\[([^\]]*)\]\]`)
	text = re.ReplaceAllString(text, "$1")

	re = regexp.MustCompile(`\[[^\]]*\]`)
	text = re.ReplaceAllString(text, "")

	re = regexp.MustCompile(`<[^>]*>`)
	text = re.ReplaceAllString(text, "")

	re = regexp.MustCompile(`&[a-zA-Z]+;`)
	text = re.ReplaceAllString(text, "")

	// removin extra whitespace
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}
