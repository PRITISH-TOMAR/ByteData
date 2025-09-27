package main

import (
	"fmt"

	"github.com/PRITISH-TOMAR/byted/internal/cli"
	"github.com/PRITISH-TOMAR/byted/internal/config"
	"github.com/PRITISH-TOMAR/byted/internal/kv"
)

func main() {
	// 1. Load client config
	cfg, err := config.LoadConfig("/tmp/test_config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	fmt.Println("Client Config Loaded Successfully!")
	fmt.Printf("ClientID: %s | StorageMode: %s | WALPath: %s\n", cfg.ClientID, cfg.StorageMode, cfg.WALPath)

	// 2. Start KV engine
	engine, err := kv.NewKVEngine(cfg.WALPath, 4)
	if err != nil {
		fmt.Println("Failed to start engine:", err)
		return
	}
	defer engine.Close()
	fmt.Println("ByteData Engine started successfully!")

	// Parse CLI login args
	username, password := cli.ParseLoginArgs()

	// Start CLI shell
	cli.StartCLI(username, password, engine)
}
