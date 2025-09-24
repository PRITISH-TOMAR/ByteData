package tests

import (
	"fmt"
	"testing"
	"github.com/PRITISH-TOMAR/byted/internal/config"
)

func TestConfigLoad(t *testing.T) {
	// Load client config
	cfg, err := config.LoadConfig("../test_config.json")
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	// Print loaded config to verify (useful for debug)
	t.Log("Client Config Loaded Successfully!")
	fmt.Printf("ClientID: %s\n", cfg.ClientID)
	fmt.Printf("Storage Mode: %s\n", cfg.StorageMode)
	fmt.Printf("WAL Path: %s\n", cfg.WALPath)
	fmt.Printf("Snapshot Path: %s\n", cfg.SnapshotPath)
	fmt.Printf("Max Memory (MB): %d\n", cfg.MaxMemoryMB)
	fmt.Printf("Fsync Every: %d\n", cfg.FsyncEvery)

	// Basic checks
	if cfg.ClientID == "" {
		t.Error("ClientID should not be empty")
	}
	if cfg.StorageMode == "" {
		t.Error("StorageMode should not be empty")
	}
}
