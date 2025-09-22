package main

import (
    "fmt"
    "github.com/PRITISH-TOMAR/byted/internal/config"
)

func main() {
    // Load client config
    cfg, err := config.LoadConfig("test_config.json")
    if err != nil {
        fmt.Println("Error loading config:", err)
        return
    }

    // Print loaded config to verify
    fmt.Println("Client Config Loaded Successfully!")
    fmt.Printf("ClientID: %s\n", cfg.ClientID)
    fmt.Printf("Storage Mode: %s\n", cfg.StorageMode)
    fmt.Printf("WAL Path: %s\n", cfg.WALPath)
    fmt.Printf("Snapshot Path: %s\n", cfg.SnapshotPath)
    fmt.Printf("Max Memory (MB): %d\n", cfg.MaxMemoryMB)
    fmt.Printf("Fsync Every: %d\n", cfg.FsyncEvery)
}
