package indexer

import (
	"fmt"
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

func (idx *Indexer) ProcessFile(filename string) error{
	docChan := make(chan *Document, 1000)

	// worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < idx.workers; i++ {
		wg.Add(1)
		go func ()  {
			defer wg.Done()
			
			for doc := range docChan{
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

func (idx *Indexer) addDocument(doc *Document){
	idx.mutex.Lock()
	defer idx.mutex.Unlock()

	idx.documents[doc.ID]= doc
	idx.docCount++

	for term := range doc.Terms{
		if _, exists := idx.termIndex[term]; !exists{
			idx.termIndex[term]= make([]uint32, 0)
		}

		idx.termIndex[term] = append(idx.termIndex[term], doc.ID)
	}

	if idx.docCount%1000 == 0{
		fmt.Printf("Processed %d documents... \n", idx.docCount)
	}
}