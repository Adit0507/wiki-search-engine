package indexer

import (
	"sync"

	"github.com/Adit0507/wiki-search-engine/internal/storage"
)

type Indexer struct {
	indexPath string
	workers   int
	documents map[uint32]*Document
	termIndex map[string][]uint32
	docCount  int
	avgDocLen float64
	storage   *storage.DiskStorage
	mutex     sync.RWMutex
}

func NewIndexer(indexPath string, workers int) *Indexer {
	return &Indexer{
		indexPath: indexPath,
		workers:   workers,
		documents: make(map[uint32]*Document),
		termIndex: make(map[string][]uint32),
		storage:   storage.NewDiskStorage(indexPath),
	}
}
