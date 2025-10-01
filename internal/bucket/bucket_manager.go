package bucket

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"byted/constants"
	"byted/internal/kv"
)

type BucketManager struct {
	BaseDir  string
	Buckets  map[string]*Bucket
	mutex    sync.RWMutex
	isActive *Bucket
}

type MetaData struct {
	Buckets      []string `json:"buckets"`
	ActiveBucket string   `json:"active_bucket"`
}

func (bm *BucketManager) LoadMetaData() error {

	metaPath := constants.GLOBALMETAPATH
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		// No metadata file exists, nothing to load
		return nil
	}

	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil
	}

	var meta MetaData
	if err := json.Unmarshal(data, &meta); err != nil {
		return err
	}
	bm.mutex.Lock()
	for _, bucketName := range meta.Buckets {
		bucketDir := filepath.Join(bm.BaseDir, bucketName)
		walPath := filepath.Join(bucketDir, bucketName+constants.WALFILENAME)

		kvEngine, err := kv.NewKVEngine(walPath, constants.DEFAULTREEORDER)
		if err != nil {
			fmt.Printf("failed to load bucket %s: %v\n", bucketName, err)
			continue
		}

		kvEngine.ReplayWAL()

		bucket := &Bucket{
			Name:     bucketName,
			KvEngine: kvEngine,
		}
		bm.Buckets[bucketName] = bucket

	}
	bm.mutex.Unlock()
	return nil
}

func (bm *BucketManager) SaveMetaData() error {
	meta := MetaData{
		Buckets:      make([]string, 0, len(bm.Buckets)),
		ActiveBucket: "",
	}

	for bucketName := range bm.Buckets {
		meta.Buckets = append(meta.Buckets, bucketName)
	}

	if bm.isActive != nil {
		meta.ActiveBucket = bm.isActive.Name
	}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}

	metaPath := constants.GLOBALMETAPATH
	return os.WriteFile(metaPath, data, constants.OWNERPERMISSION)
}

func NewBucketManager(baseDir string) (*BucketManager, error) {

	if err := os.MkdirAll(baseDir, constants.OWNERPERMISSION); err != nil {
		return nil, err
	}

	bm := &BucketManager{
		BaseDir: baseDir,
		Buckets: make(map[string]*Bucket),
	}

	bm.LoadMetaData()
	return bm, nil
}
