package cli

import (
	"bufio"
	"byted/internal/bucket"
	"fmt"
	"os"
)

func StartBucketCLI(bucketManager *bucket.BucketManager) {
	reader := bufio.NewScanner(os.Stdin)
	bucket, err := bucketManager.GetActiveBucket()
	if err != nil {
		fmt.Println("Error retrieving active bucket:", err)
		return
	}

	for {

		fmt.Printf("Bytedata: [%s]> ", bucket.Name)
		if !reader.Scan() {
			break
		}
		// cmdLine := reader.Text()
		// err := ExecuteCommand(cmdLine, bucket)

		if err != nil && err.Error() == "exit" {
			bucketManager.ExitBucket()
			break
		}
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
