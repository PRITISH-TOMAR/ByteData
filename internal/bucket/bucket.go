package bucket

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"github.com/PRITISH-TOMAR/byted/constants"
	"github.com/PRITISH-TOMAR/byted/internal/kv"
)

type Bucket struct{
	Name string
	kvEngine *kv.KVEngine
}

func NewBucket(name, baseDir string, btreeOrder int) (*Bucket, error){
	
	walPath := filepath.Join(baseDir, name + constants.WALFILENAME)
	kvEngine, err := kv.NewKVEngine(walPath, btreeOrder)
	if err != nil{
		return nil, err
	}
	bucket := &Bucket{
		Name: name,
		kvEngine: kvEngine,
	}
	return bucket, nil
}

// Close closes the WAL file associated with the bucket.
func (b* Bucket) Close() error{
	return b.kvEngine.Close()
}


func (bm *BucketManager) CreateBucket(name string, btreeOrder int) error {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	if _, exists := bm.Buckets[name]; exists {
		return fmt.Errorf("bucket %s already exists", name)
	}
	bucketDir := filepath.Join(bm.BaseDir, name)
	if err := os.MkdirAll(bucketDir, constants.BUCKETPERMISSION); err != nil {
		return err
	}

	bucket, err := NewBucket(name, bm.BaseDir, btreeOrder)
	if err != nil {
		return err
	}

	bm.Buckets[name] = bucket
	bm.SaveMetaData()
	return nil
}

func (bm *BucketManager) UseBucket(name string) (*Bucket, error) {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	bucket, exists := bm.Buckets[name]
	if !exists {
		return nil, fmt.Errorf("bucket %s does not exist", name)
	}
	bm.isActive = bucket
	fmt.Printf("Switched to bucket: %s\n", name)
	return bucket, nil
}

func (bm *BucketManager) GetActiveBucket() (*Bucket, error) {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	if bm.isActive == nil {
		return nil, fmt.Errorf("no active bucket selected")
	}
	return bm.isActive, nil
}

func (bm *BucketManager) ListBuckets(input string) []string {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	bucketNames := make([]string, 0, len(bm.Buckets))
	for name := range bm.Buckets {
		if input == "" || strings.Contains(name, input) {
			bucketNames = append(bucketNames, name)
		}
	}
	return bucketNames
}

func (bm *BucketManager) ExitBucket() error {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	if bm.isActive == nil {
		return fmt.Errorf("no active bucket to exit")
	}
	bm.isActive = nil
	return nil
}

func (bm *BucketManager) DropBucket(name string) error {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	bucket, exists := bm.Buckets[name]
	if !exists {
		return fmt.Errorf("bucket %s does not exist", name)
	}

	if err := bucket.Close(); err != nil {
		return err
	}

	bucketDir := filepath.Join(bm.BaseDir, name)
	if err := os.RemoveAll(bucketDir); err != nil {
		return err
	}

	delete(bm.Buckets, name)
	bm.SaveMetaData()
	return nil
}


