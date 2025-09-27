package kv

import (
	"errors"
	"fmt"

	"github.com/PRITISH-TOMAR/byted/internal/wal"
	"github.com/PRITISH-TOMAR/byted/internal/btree"
)

// valueMeta holds the value and its last associated LSN.
type valueMeta struct {
	value []byte
	lsn   uint64
}

type KVEngine struct {
	wal *wal.WAL  // write ahead logs for durability
	pointIndex map[string] * valueMeta // in-memory point index
	index *btree.BPlusTree // on-disk b+tree for range queries
}

// NewKVEngine initializes the key-value engine with WAL and B+ tree.
func NewKVEngine(walPath string, btreeOrder int) (*KVEngine, error) {
	

	// opens wal (create if not exists) and recovers last LSN
	w, err := wal.New(walPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize WAL: %w", err)
	}

	// create B Tree index
	tree:= btree.New(btreeOrder)
	
	engine := &KVEngine{
		wal: w,
		pointIndex: make(map[string]*valueMeta),
		index: tree,
	}

	// replay Wal to rebuild memory state
	err = engine.ReplayWAL()
	if err != nil {
		_ = w.Close()
		return nil, fmt.Errorf("WAL replay failed: %w", err)
	}
	fmt.Println("WAL replayed successfully, in-memory state restored.")
	return engine, nil
}

// Close gracefully closes the KV engine, ensuring all data is flushed.
func (kv *KVEngine) Close() error {
	if kv.wal != nil {
		if err := kv.wal.Close(); err != nil {
			return fmt.Errorf("failed to close WAL: %w", err)
		}
	}
	// Add any additional cleanup if necessary
	return nil
}


// reads wal from start and applies each record to in-memory structure.
func (kv * KVEngine) ReplayWAL() error {

	if kv.wal == nil {
		return errors.New("WAL is not initialized")
	}

	// handler to process each WAL record
	handler := func(lsn uint64, recordType uint8, key, value []byte) error {
		switch recordType {

		case wal.RecordPut:
			// create valueMeta and update point index

			vm := &valueMeta{ value: make([]byte, len(value)), lsn: lsn }
			copy(vm.value, value)
			kv.pointIndex[string(key)] = vm
			kv.index.Insert(string(key), value) // also insert into B+ tree

		case wal.RecordDelete:
			// delete from point index
			delete(kv.pointIndex, string(key))
			kv.index.Delete(string(key)) // also delete from B+ tree
			

		default:
			return fmt.Errorf("unknown record type: %d", recordType)
		}
		return nil
	}

	// replay WAL using the handler
	if err := kv.wal.Replay(handler); err != nil {
		return fmt.Errorf("WAL replay error: %w", err)
	}
	return nil
}

// Put adds or updates a key-value pair in the KV engine.
func (kv *KVEngine) Put(key, value []byte) (uint64, error) {
	if kv.wal == nil {
		return 0, errors.New("WAL is not initialized")
	}
	// append to WAL
	lsn, err := kv.wal.AppendPut(key, value)
	if err != nil {
		return 0, fmt.Errorf("failed to append PUT to WAL: %w", err)
	}

	// update in-memory point index
	vm := &valueMeta{ value: make([]byte, len(value)), lsn: lsn }
	copy(vm.value, value)
	kv.pointIndex[string(key)] = vm

	// insert into B+ tree for range queries
	kv.index.Insert(string(key), value)

	return lsn, nil
}

// Get retrieves the value for a given key.
func (kv *KVEngine) Get(key []byte) ([]byte, error){
	vm, ok := kv.pointIndex[string(key)]
	if !ok {
		return nil, errors.New("key not found")
	}
	// return a copy to prevent external modification
	valueCopy := make([]byte, len(vm.value))
	copy(valueCopy, vm.value)
	return valueCopy, nil
}

// Delete removes a key-value pair from the KV engine.
func (kv *KVEngine) Delete(key []byte) (uint64, error) {
	if kv.wal == nil {
		return 0, errors.New("WAL is not initialized")
	}
	// append delete record to WAL
	lsn, err := kv.wal.AppendDelete(key)
	if err != nil {
		return 0, fmt.Errorf("failed to append DELETE to WAL: %w", err)
	}

	// remove from in-memory point index
	delete(kv.pointIndex, string(key))

	// remove from B+ tree
	kv.index.Delete(string(key))

	return lsn, nil
}

// Range retrieves all key-value pairs within the specified key range [startKey, endKey].
func (kv *KVEngine) Range(startKey, endKey []byte) ([]btree.KVPair) {
	return kv.index.RangeQuery(string(startKey), string(endKey))
}