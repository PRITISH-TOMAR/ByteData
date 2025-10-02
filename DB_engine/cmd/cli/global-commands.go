package cli

import (
	"byted/DB_engine/constants"
	"byted/DB_engine/core/bucket"
	"fmt"
	"net"
	"strings"
)


func ExecuteGlobalCommmand(input string, bucketManager *bucket.BucketManager, conn net.Conn) ([]string, error) {
	// Parse the command line input

	input = strings.TrimSpace(input)
	if input == "" {
		return nil, nil
	}
	parts := strings.Fields(input)
	command := parts[0]

	switch command {
	case "exit", "quit":
		return handleExit(conn)

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

	

	default:
		return handleHelp(0)

	}
}

func handleListBuckets(input string, bucketManager *bucket.BucketManager) ([]string, error) {
	return bucketManager.ListBuckets(input), nil
}

func handleUseBucket(parts []string, bucketManager *bucket.BucketManager) ([]string, error) {
	if len(parts) != 2 {
		return nil, fmt.Errorf("usage: use <bucket_name>")
	}

	if _, err := bucketManager.UseBucket(parts[1]); err != nil {
		return nil, err
	}

	return []string{fmt.Sprintf("Switched to bucket '%s'", parts[1])}, nil
}

func handleCreateBucket(parts []string, bucketManager *bucket.BucketManager) ([]string, error) {
	// enc := comm.Enc
	if len(parts) != 2 {
		// enc.Encode(structs.Message{Type: "error", Message:"usage: create <bucket_name>"})
		return nil, fmt.Errorf("usage: create <bucket_name>")
	}
	if err := bucketManager.CreateBucket(parts[1], constants.DEFAULTREEORDER); err != nil {
		return nil, err
	}
	return []string{"Bucket created successfully."}, nil
}

func handleDropBucket(parts []string, bucketManager *bucket.BucketManager) ([]string, error) {
	if len(parts) != 2 {
		return nil, fmt.Errorf("usage: drop <bucket_name>")
	}
	if err := bucketManager.DropBucket(parts[1]); err != nil {
		return nil, err
	}

	return []string{fmt.Sprintf("Bucket '%s' dropped successfully.\n", parts[1])}, nil
}

// func showActiveBucket(bucketManager *bucket.BucketManager) ([]string, error) {
// 	activeBucket, err := bucketManager.GetActiveBucket()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return []string{activeBucket.Name}, nil
// }

func handleExit(conn net.Conn) ([]string, error) {
	conn.Close()
	return []string{"Exiting..."}, nil
}

func handleHelp(typeHelp int) ([]string, error) {

	if typeHelp == 0 {
		return printHelpGlobal()
	}
	return printHelpBucket()

}

func printHelpGlobal() ([]string, error) {
	help := []string{`Available commands:
  create <bucket_name>         - Create a new bucket
  list                         - List all buckets
  use <bucket_name>            - Switch to the specified bucket
  drop <bucket_name>           - Delete the specified bucket
  pwb                          - Print the active bucket
  exit / quit                  - Exit the CLI
  help                         - Show this help message`}
	return help, nil
}
