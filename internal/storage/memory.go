package storage

import "github.com/Adit0507/wiki-search-engine/internal/models"

type MemoryStorage struct {
	documents map[uint32]*models.Document
	termIndex map[string][]uint32
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		documents: make(map[uint32]*models.Document),
		termIndex: make(map[string][]uint32),
	}
}

func (ms *MemoryStorage) AddDocument(doc *models.Document) {
	ms.documents[doc.ID] = doc

	for term := range doc.Terms {
		ms.termIndex[term] = append(ms.termIndex[term], doc.ID)
	}
}

func (ms *MemoryStorage) GetDocument(id uint32) *models.Document {
	return ms.documents[id]
}

func (ms *MemoryStorage) GetDocumentsForTerm(term string) []uint32 {
	return ms.termIndex[term]
}
