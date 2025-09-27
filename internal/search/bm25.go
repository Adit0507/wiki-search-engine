package search

import "github.com/Adit0507/wiki-search-engine/internal/models"

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
		docCount: docCount,
		avgDocLen: avgDocLen,
	}
}
