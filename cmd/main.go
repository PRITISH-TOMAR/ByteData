package main

import (
	"fmt"

	"github.com/PRITISH-TOMAR/byted/cmd/cli"
	"github.com/PRITISH-TOMAR/byted/constants"
	"github.com/PRITISH-TOMAR/byted/internal/auth"
	"github.com/PRITISH-TOMAR/byted/internal/config"
	"github.com/PRITISH-TOMAR/byted/internal/kv"
)

func main() {

	// 1 First time setup
	var username, password string

	if !auth.AuthExists() {
		var err error
		username, password, err = auth.FirstTimeSetup()
		if err != nil {
			fmt.Println("Error during setup:", err)
			return
		}
	} else {
		var err error
		username, password, err = cli.ParseLoginArgs()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	//  2. Validate user credentials
	if err := auth.ValidateUser(username, password); err != nil {
		return
	}

	// ---------------------------------------------------------------------------------------------------
	// 3. Load client config
	cfg, err := config.LoadConfig(constants.CONFIGFILEPATH)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	fmt.Println("Client Config Loaded Successfully!")
	fmt.Printf("ClientID: %s | StorageMode: %s | WALPath: %s\n", cfg.ClientID, cfg.StorageMode, cfg.WALPath)

	// 4.  Start KV engine
	engine, err := kv.NewKVEngine(constants.DEFAULTWALPATH, 4)
	if err != nil {
		fmt.Println("Failed to start engine:", err)
		return
	}
	defer engine.Close()
	fmt.Println("ByteData Engine started successfully!")

	// 5. Start CLI shell
	cli.StartCLI(username, password, engine)
}
