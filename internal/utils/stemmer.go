package utils

import (
	"strings"

	"github.com/kljensen/snowball"
)

func Stem(word string) string {
	stemmed, err := snowball.Stem(word, "english", true)
	if err != nil{
		return strings.ToLower(word)
	}

	return stemmed
}
// porter stemming
func simpleStem(word string) string {
	word = strings.ToLower(word)

	suffixes := []string{"ing", "ly", "ed", "ies", "ied", "ies", "s"}

	for _, suffix := range suffixes{
		if strings.HasSuffix(word, suffix) && len(word) > len(suffix)+2 {
			word = word[: len(word)- len(suffix)]
			break
		}
	}

	return word
}