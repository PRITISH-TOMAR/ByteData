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