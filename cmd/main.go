package main

import (
	"fmt"
	"github.com/PRITISH-TOMAR/byted/cmd/cli"
	"github.com/PRITISH-TOMAR/byted/internal/bucket"
	"github.com/PRITISH-TOMAR/byted/constants"
)

func main() {

	// 1 Authenticating USer via CLI
	username,err :=cli.AuthCLI()
	if err != nil{
		fmt.Println("Authentication failed:", err)
		return
	}

	// 2. Initializing Bucket Manager
	bucketManager, err := bucket.NewBucketManager(constants.DBBUCKETSPATH)
	bucketManager.LoadMetaData()
	if err != nil {
		fmt.Println("Failed to initialize Bucket Manager:", err)
		return
	}
	// 3. Entering to CLI parent - Bytedata Shell
	cli.StartCLI( username, bucketManager)
}
