package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Adit0507/wiki-search-engine/internal/indexer"
)

func main() {
	var (
		dataPath  = flag.String("data", "./data/wikipedia", "Path to wikipedia data")
		indexPath = flag.String("index", "./indexes", "Path to store indexes")
		workers   = flag.Int("workers", 4, "No. of worker goroutines")
	)
	flag.Parse()

	if err := os.MkdirAll(*indexPath, 0755); err != nil {
		log.Fatal("Failed to create index directory: ", err)
	}

	fmt.Println("Starting wikipedia parser...")
	fmt.Printf("Data path: %s\n", *dataPath)
	fmt.Printf("Index path: %s\n", *indexPath)
	fmt.Printf("Workers: %d\n", *workers)

	idx := indexer.NewIndexer(*indexPath, *workers)

	err := filepath.Walk(*dataPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (filepath.Ext(path) == ".xml" || filepath.Ext(path) == ".bz2") {
			fmt.Printf("Processing file: %s\n", path)
			return idx.ProcessFile(path)
		}

		return nil
	})

	if err != nil {
		log.Fatal("Error processing files: ", err)
	}

	fmt.Println("Building index...")
	if err := idx.BuildIndex(); err != nil {
		log.Fatal("error building index: ", err)
	}

	fmt.Println("Saving index to disk")
	if err := idx.SaveToDisk(); err != nil {
		log.Fatal("error saving index: ", err)
	}

	fmt.Println("Indexing completed")
}
