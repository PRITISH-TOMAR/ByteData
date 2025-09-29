package bucket

import(
	"fmt"
	"path/filepath"
	"sync"
	"os"
	"strings"
	"github.com/PRITISH-TOMAR/byted/constants"
)

type BucketManager struct{
	BaseDir string
	Buckets map[string]*Bucket
	mutex sync.RWMutex
	isActive *Bucket
}

func NewBucketManager(baseDir string) (*BucketManager, error){
	
	if err:= os.MkdirAll(baseDir, constants.OWNERPERMISSION); err != nil{
		return nil, err
	}

	return &BucketManager{
		BaseDir: baseDir,
		Buckets: make(map[string]*Bucket),
		
	}, nil
}

func (bm *BucketManager) CreateBucket(name string, btreeOrder int) error{
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	if _, exists := bm.Buckets[name]; exists{
		return fmt.Errorf("bucket %s already exists", name)
	}
	bucketDir := filepath.Join(bm.BaseDir, name)
	if err := os.MkdirAll(bucketDir, constants.BUCKETPERMISSION); err != nil{
		return err
	}

	bucket, err := NewBucket(name, bm.BaseDir, btreeOrder)
	if err != nil{
		return err
	}

	bm.Buckets[name] = bucket
	return nil
}

func (bm * BucketManager) UseBucket(name string) (*Bucket, error){
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	bucket, exists := bm.Buckets[name]
	if !exists{
		return nil, fmt.Errorf("bucket %s does not exist", name)
	}
	bm.isActive = bucket
	return bucket, nil
}

func( bm * BucketManager) GetActiveBucket() (*Bucket, error){
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	if bm.isActive == nil{
		return nil, fmt.Errorf("no active bucket selected")
	}
	return bm.isActive, nil
}

func (bm * BucketManager) ListBuckets(input string ) []string{
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()
	
	bucketNames := make([]string, 0, len(bm.Buckets))
	for name := range bm.Buckets{
		if input == "" || strings.Contains(name, input){
			bucketNames = append(bucketNames, name)
		}
	}
	return bucketNames
}

func (bm * BucketManager) DropBucket(name string) error{
	bm.mutex.Lock()
	defer bm.mutex.Unlock()
	
	bucket, exists := bm.Buckets[name]
	if !exists{
		return fmt.Errorf("bucket %s does not exist", name)
	}

	if err := bucket.Close(); err != nil{
		return err
	}
	
	bucketDir := filepath.Join(bm.BaseDir, name)
	if err := os.RemoveAll(bucketDir); err != nil{
		return err
	}

	delete(bm.Buckets, name)
	return nil
}