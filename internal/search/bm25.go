package search

import (
	"math"
	"strings"

	"github.com/Adit0507/wiki-search-engine/internal/models"
)

const (
	K1 = 1.2
	B  = 0.75
)

type BM25 struct {
	documents map[uint32]*models.Document
	termIndex map[string][]uint32
	docCount  int
	avgDocLen float64
}

func NewBM25(documents map[uint32]*models.Document, termIndex map[string][]uint32, docCount int, avgDocLen float64) *BM25 {
	return &BM25{
		documents: documents,
		termIndex: termIndex,
		docCount:  docCount,
		avgDocLen: avgDocLen,
	}
}

func (bm *BM25) generateSnippet(doc *models.Document, terms []string, maxLen int) string {
	content := doc.Content
	if len(content) <= maxLen {
		return content
	}

	// best positiomn to extract snippet
	lower := strings.ToLower(content)
	bestPos := 0
	bestScore := 0

	for i := 0; i < len(content)-maxLen; i += 50 {
		score := 0
		snippet := lower[i : i+maxLen]

		for _, term := range terms {
			score += strings.Count(snippet, term)
		}

		if score > bestScore {
			bestScore = score
			bestPos = i
		}
	}

	snippet := content[bestPos : bestPos+maxLen]

	// trim to word boundaries
	if bestPos > 0 {
		if spaceIdx := strings.Index(snippet, " "); spaceIdx > 0 {
			snippet = snippet[spaceIdx+1:]
		}
	}
	if lastSpace := strings.LastIndex(snippet, " "); lastSpace > 0 {
		snippet = snippet[:lastSpace]
	}

	return snippet + "..."
}

func (bm *BM25) getCandidateDocuments(terms []string) map[uint32]bool {
	candidates := make(map[uint32]bool)

	for _, term := range terms {
		if docList, exists := bm.termIndex[term]; exists {
			for _, docId := range docList {
				candidates[docId] = true
			}
		}
	}

	return candidates
}

func (bm *BM25) calculateBM25Score(terms []string, doc *models.Document) float64 {
	score := 0.0

	for _, term := range terms {
		// term frequency
		tf := float64(doc.GetTermFreq(term))
		if tf == 0 {
			continue
		}

		// document frequency
		df := float64(len(bm.termIndex[term]))
		if df == 0 {
			continue
		}

		// idf measures the importnce of term across the corpus
		idf := math.Log((float64(bm.docCount) - df + 0.5) / (df + 0.5))

		// doc length normalization
		docLen := float64(doc.GetLength())
		normalization := K1 * ((1 - B) + B*(docLen/bm.avgDocLen))

		// bm25 formula
		termScore := idf * (tf * (K1 + 1)) / (tf + normalization)
		score += termScore
	}

	return score
}
