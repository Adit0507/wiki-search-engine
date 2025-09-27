package storage

import (
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/Adit0507/wiki-search-engine/internal/models"
)

type DiskStorage struct {
	indexPath string
}

func NewDiskStorage(indexPath string) *DiskStorage {
	return &DiskStorage{indexPath: indexPath}
}

func (ds *DiskStorage) SaveDocuments(documents map[uint32] *models.Document) error {
	file, err := os.Create(filepath.Join(ds.indexPath, "documents.gob"))
    if err != nil {
        return err
    }
    defer file.Close()
    
    encoder := gob.NewEncoder(file)
    return encoder.Encode(documents)
}

func (ds *DiskStorage) SaveTermIndex(termIndex map[string][] uint32) error {
	file, err := os.Create(filepath.Join(ds.indexPath, "terms.gob"))
	if err != nil{
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(termIndex)
}