package search

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Adit0507/wiki-search-engine/internal/storage"
)

type Engine struct {
	bm25 *BM25
}

func (e *Engine) Search(query string, limit int) ([]Result, error) {
	return e.bm25.Search(query, limit)
}

func NewEngine(indexPath string) (*Engine, error) {
	storage := storage.NewDiskStorage(indexPath)

	// load metadata
	metaFile := filepath.Join(indexPath, "metadata.json")
	metaData, err := os.ReadFile(metaFile)
	if err != nil {
		return nil, err
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(metaData, &metadata); err != nil {
		return nil, err
	}

	docCount := int(metadata["doc_count"].(float64))
	avgDocLen := metadata["avg_doc_len"].(float64)

	// loading documents
	documents, err := storage.LoadDocuments()
	if err != nil {
		return nil, err
	}
	// loaidng term index
	termIndex, err := storage.LoadTermIndex()
	if err != nil {
		return nil, err
	}

	bm25 := NewBM25(documents, termIndex, docCount, avgDocLen)

	return &Engine{bm25: bm25}, nil
}
