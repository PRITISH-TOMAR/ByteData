package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/PRITISH-TOMAR/byted/internal/bucket"
	"github.com/PRITISH-TOMAR/byted/constants"
)

func ExecuteGlobalCommmand(input string, bucketManager *bucket.BucketManager) error {
	// Parse the command line input

	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}
	parts := strings.Fields(input)
	command := parts[0]

	switch command {
	case "exit", "quit":
		return handleExit()

	case "help":
		return handleHelp(0)

	case "list":
		return handleListBuckets("", bucketManager)

	case "use":
		return handleUseBucket(parts, bucketManager)

	case "create":
		return handleCreateBucket(parts, bucketManager)

	case "drop":
		return handleDropBucket(parts, bucketManager)
	
	case "pwb":
		return showActiveBucket(bucketManager)

	default:
		 fmt.Printf("unknown command: %s", command)
		 return handleHelp(0)

	}

	return nil
}

func handleListBuckets(input string, bucketManager *bucket.BucketManager) error {
	fmt.Println("Buckets:")

	for _, bucket := range bucketManager.ListBuckets(input) {
		{
			fmt.Printf(" - %s\n", bucket)
		}
	}
	return nil

}

func handleUseBucket(parts []string, bucketManager *bucket.BucketManager) error {
	if len(parts) != 2 {
		return fmt.Errorf("usage: use <bucket_name>")
	}
	if _, err := bucketManager.UseBucket(parts[1]); err != nil {
		return err
	}
	
	StartBucketCLI(bucketManager)
	return nil
}

func handleCreateBucket(parts []string, bucketManager *bucket.BucketManager) error {
	if len(parts) != 2 {
		return fmt.Errorf("usage: create <bucket_name>")
	}
	if err := bucketManager.CreateBucket(parts[1], constants.DEFAULTREEORDER); err != nil {
		return err
	}
	fmt.Printf("Bucket '%s' created successfully.\n", parts[1])
	return nil
}

func handleDropBucket(parts []string, bucketManager *bucket.BucketManager) error {
	if len(parts) != 2 {
		return fmt.Errorf("usage: drop <bucket_name>")
	}
	if err := bucketManager.DropBucket(parts[1]); err != nil {
		return err
	}
	fmt.Printf("Bucket '%s' dropped successfully.\n", parts[1])
	return nil
}


func showActiveBucket(bucketManager *bucket.BucketManager) error {
	activeBucket, err := bucketManager.GetActiveBucket()
	if err != nil {
		return err
	}
	fmt.Printf("Active Bucket: %s\n", activeBucket.Name)
	return nil
}
func handleExit() error {
	fmt.Println("Exiting...")
	os.Exit(0)
	return nil
}

func handleHelp(typeHelp int) error {

	if typeHelp == 0 {
		printHelpGlobal()
	} else {
		printHelpBucket()
	}
	return nil
}

func printHelpGlobal() {
	fmt.Println("Available commands:")
	fmt.Println("  create <bucket_name>         - Create a new bucket")
	fmt.Println("  list                         - List all buckets")
	fmt.Println("  use <bucket_name>            - Switch to the specified bucket")
	fmt.Println("  drop <bucket_name>           - Delete the specified bucket")
	fmt.Println("  pwb                          - Print the active bucket")
	fmt.Println("  exit / quit                  - Exit the CLI")
	fmt.Println("  help                         - Show this help message")
}
