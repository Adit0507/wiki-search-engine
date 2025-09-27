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
