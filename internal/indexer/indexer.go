package indexer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Adit0507/wiki-search-engine/internal/models"
	"github.com/Adit0507/wiki-search-engine/internal/storage"
)

type Indexer struct {
	indexPath string
	workers   int
	documents map[uint32]*models.Document
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
		documents: make(map[uint32]*models.Document),
		termIndex: make(map[string][]uint32),
		storage:   storage.NewDiskStorage(indexPath),
	}
}

func (idx *Indexer) ProcessFile(filename string) error {
	docChan := make(chan *models.Document, 1000)

	// worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < idx.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for doc := range docChan {
				idx.addDocument(doc)
			}
		}()
	}

	// pain file
	parser := NewParser(docChan)
	err := parser.ParseFile(filename)
	close(docChan)
	wg.Wait()

	return err
}

func (idx *Indexer) addDocument(doc *models.Document) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()

	idx.documents[doc.ID] = doc
	idx.docCount++

	for term := range doc.Terms {
		if _, exists := idx.termIndex[term]; !exists {
			idx.termIndex[term] = make([]uint32, 0)
		}

		idx.termIndex[term] = append(idx.termIndex[term], doc.ID)
	}

	if idx.docCount%1000 == 0 {
		fmt.Printf("Processed %d documents... \n", idx.docCount)
	}
}

func (idx *Indexer) BuildIndex() error {
	fmt.Println("building index structures")
	// doc lenth
	totalLen := 0
	for _, doc := range idx.documents {
		totalLen += doc.Length
	}
	if idx.docCount > 0 {
		idx.avgDocLen = float64(totalLen) / float64(idx.docCount)
	}

	fmt.Printf("Total documents: %d\n", idx.docCount)
	fmt.Printf("Total terms: %d\n", len(idx.termIndex))
	fmt.Printf("Average document length: %.2f\n", idx.avgDocLen)

	return nil
}

func (idx *Indexer) SaveToDisk() error {
	metadata := map[string]interface{}{
        "doc_count":    idx.docCount,
        "avg_doc_len":  idx.avgDocLen,
        "total_terms":  len(idx.termIndex),
    }

	metaFile := filepath.Join(idx.indexPath, "metadata.json")
	metaData, _ := json.Marshal(metadata)
	if err := os.WriteFile(metaFile, metaData, 0644); err != nil{
		return err
	}

	// savin documents
	fmt.Println("saving docs...")
	if err := idx.storage.SaveDocuments(idx.documents); err != nil {
		return err
	}

	// savin term index
	fmt.Println("saving term index")
	if err := idx.storage.SaveTermIndex(idx.termIndex); err != nil{
		return err
	}

	return  nil
}