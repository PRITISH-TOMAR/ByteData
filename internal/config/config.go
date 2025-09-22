package config

// Package declaration : file -> config package

// to import external packages we need
import (
	"encoding/json"
	"os"
)

// Dfeine Client Configuration structure
type ClientConfig struct {
	ClientID     string `json:"client_id"`     // unique client id
	StorageMode  string `json:"storage_mode"`  // memory/disk/hybrid
	WALPath      string `json:"wal_path"`      // WAL logs address
	SnapshotPath string `json:"snapshot_path"` // Snapshot address
	MaxMemoryMB  int    `json:"max_memory_mb"` // Max RAM requirement
	FsyncEvery   int    `json:"fsync_every"`   // how often to force sync to disk
}

// LoaddConfig function reads config set by client from JSON file

// function NAME (parameters) (return types)
func LoadConfig(filePath string) (*ClientConfig, error) {
	// Return : a pointer to clientconfig ( nil if any error)
	// 			an error value : if successfull- nil)

	// Step 1: open the file
	f, err := os.Open(filePath)

	if err != nil {
		return nil, err // error if file unfound
	}

	// Step 2: closing the file
	defer f.Close()

	// STep 3: Prepare to decode JSON
	// declaring a type ClientConfig variable named as cfg.
	var cfg ClientConfig

	// Step 4: Decoding the file content
	decoder := json.NewDecoder(f)

	// STep 6: If any error return {nill, error}
	if err := decoder.Decode((&cfg)); err != nil {
		return nil, err // error if json invalid
	}

	// Step 7: Return final config read from file.
	return &cfg, nil
}
