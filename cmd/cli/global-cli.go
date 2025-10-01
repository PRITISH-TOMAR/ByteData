package cli

import (
	"bufio"
	"fmt"
	"os"

	"byted/internal/bucket"
)

// StartCLI starts the ByteData shell for a given KVEngine
func StartCLI(username string, bucketManager *bucket.BucketManager) {

	fmt.Printf("Welcome %s! Connected to ByteData Engine.\n", username)

	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Bytedata> ")
		if !reader.Scan() {
			break
		}
		// cmdLine := reader.Text()
		// err := ExecuteGlobalCommmand(cmdLine, bucketManager)
		// if err != nil {
		// 	fmt.Println("Error:", err)
		// }
	}
	fmt.Println("\nBye!")
}
