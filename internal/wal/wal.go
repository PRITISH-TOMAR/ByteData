package wal

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

// wal record types
const (
	RecordPut    = 1
	RecordDelete = 2
)

type WAL struct {
	f       *os.File // underlying file
	lastLSN uint64   // last log sequence number( monotonically increasing)
}

// opens/ created a WAL fule at path. If file exits, it is opened in append mode.
func New(path string) (*WAL, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &WAL{f: f}, nil
}

// Close closes the underlying file.
func (w *WAL) Close() error {
	if w.f == nil {
		return errors.New("WAL file is not open")
	}
	return w.f.Close()
}

// Put writes a put record to the WAL.
func (W *WAL) AppendPut(key, value []byte) (uint64, error) {
	return w.appendRecord(RecordPut, key, value)
}

// Delete writes a delete record to the WAL.
func (w *WAL) AppendDelete(key []byte) (uint64, error) {
	return w.appendRecord(RecordDelete, key, nil)
}

// appendRecord is the low-level writer.
// Format: | uint32 totalLen | uint64 LSN(8) | unit8 Type(1) | uint32 KeySize | uint32 ValueSize | Key | Value |

func (w *WAL) appendRecord(recordType uint8, key string, value []byte) (uint64, error) {

	if w.f == nil {
		return 0, errors.New("WAL file is not open")
	}
	// increment LSN
	w.lastLSN++
	lsn := w.lastLSN

	KeySize := uint32(len(key))
	ValueSize := uint32(0)

	if value != nil {
		ValueSize = uint32(len(value))
	}
	// total length of the record- ( 8 lsn + 1 type + 4 key size + 4 value size + key + value) //bytes
	totalLen := 8 + 1 + 4 + 4 + KeySize + ValueSize


	// buffer writer in sequence
	buf := make([]byte, 4+totalLen) // 4 bytes for totalSize too.

	curr := 0

	// write total length
	binary.LittleEndian.PutUint32(( buf[curr : ] ), totalLen); curr +=4 
	// littleEndian - least significant byte first , particularly used for bytes

	// write LSN
	binary.LittleEndian.PutUint64(( buf[curr : ] ), lsn); curr += 8
	// write record type
	buf[curr] = recordType; curr += 1
	// write key size
	binary.LittleEndian.PutUint32(( buf[curr : ] ), KeySize); curr += 4
	// write value size
	binary.LittleEndian.PutUint32(( buf[curr : ] ), ValueSize); curr += 4
	// write key

	// used copy instead of littleEndian for keySIze to be a variable length
	// particularly used in case of strings or byte slices
	copy(buf[curr:curr+int(KeySize)], key); curr += int(KeySize)
	// write value if present
	if ValueSize > 0 {
		copy(buf[curr:curr+int(ValueSize)], value)
	}

	// writing to file
	if _, err := w.f.Write(buf); err != nil {
		return 0, err
	}
	
	// disk persistence for durability
	if err := w.f.Sync(); err != nil {
		return 0, err
	}

	return lsn, nil
}
