package bucket

import (
	"path/filepath"
	"sync"
	"github.com/PRITISH-TOMAR/byted/constants"
	"github.com/PRITISH-TOMAR/byted/internal/btree"
	"github.com/PRITISH-TOMAR/byted/internal/wal"
)

type Bucket struct{
	Name string
	Tree *btree.BPlusTree
	WAL *wal.WAL
	mutex sync.RWMutex
}

func NewBucket(name, baseDir string, btreeOrder int) (*Bucket, error){
	walPath := filepath.Join(baseDir, name + constants.WALFILENAME)
	w, err := wal.New(walPath)
	if err != nil{
		return nil, err
	}
	tree := btree.New(btreeOrder)

	bucket := &Bucket{
		Name: name,
		Tree: tree,
		WAL: w,
	}
	return bucket, nil
}


// Close closes the WAL file associated with the bucket.
func (b* Bucket) Close() error{
	return b.WAL.Close()
}

