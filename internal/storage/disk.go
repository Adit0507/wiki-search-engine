package storage

type DiskStorage struct {
	indexPath string
}

func NewDiskStorage(indexPath string) *DiskStorage {
	return &DiskStorage{indexPath: indexPath}
}
