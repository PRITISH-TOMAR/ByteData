package cli

import (
	"fmt"
	"net"
	"os"
	"strings"

	"byted/constants"
	"byted/internal/bucket"
)

func ExecuteGlobalCommmand(conn net.Conn, input string, bucketManager *bucket.BucketManager) error {
	// Parse the command line input

	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}
	parts := strings.Fields(input)
	command := parts[0]

	switch command {
	case "exit", "quit":
		return handleExit(conn)

	case "help":
		return handleHelp(0, conn)

	case "list":
		return handleListBuckets("", bucketManager, conn)

	case "use":
		return handleUseBucket(parts, bucketManager, conn)

	case "create":
		return handleCreateBucket(parts, bucketManager, conn)

	case "drop":
		return handleDropBucket(parts, bucketManager, conn)

	case "pwb":
		return showActiveBucket(bucketManager, conn)

	default:
		fmt.Fprintf(conn, "unknown command: %s", command)
		return handleHelp(0, conn)

	}

	return nil
}

func handleListBuckets(input string, bucketManager *bucket.BucketManager, conn net.Conn) error {
	fmt.Println("Buckets:")

	for _, bucket := range bucketManager.ListBuckets(input) {
		{
			fmt.Fprintf(conn, " - %s\n", bucket)
		}
	}
	return nil

}

func handleUseBucket(parts []string, bucketManager *bucket.BucketManager, conn net.Conn) error {
	if len(parts) != 2 {
		return fmt.Errorf("usage: use <bucket_name>")
	}
	if _, err := bucketManager.UseBucket(parts[1]); err != nil {
		return err
	}

	StartBucketCLI(bucketManager, conn)
	return nil
}

func handleCreateBucket(parts []string, bucketManager *bucket.BucketManager, conn net.Conn) error {
	if len(parts) != 2 {
		return fmt.Errorf("usage: create <bucket_name>")
	}
	if err := bucketManager.CreateBucket(parts[1], constants.DEFAULTREEORDER); err != nil {
		return err
	}
	fmt.Fprintf(conn, "Bucket '%s' created successfully.\n", parts[1])
	return nil
}

func handleDropBucket(parts []string, bucketManager *bucket.BucketManager, conn net.Conn) error {
	if len(parts) != 2 {
		return fmt.Errorf("usage: drop <bucket_name>")
	}
	if err := bucketManager.DropBucket(parts[1]); err != nil {
		return err
	}
	fmt.Fprintf(conn, "Bucket '%s' dropped successfully.\n", parts[1])
	return nil
}

func showActiveBucket(bucketManager *bucket.BucketManager, conn net.Conn) error {
	activeBucket, err := bucketManager.GetActiveBucket()
	if err != nil {
		return err
	}
	fmt.Fprintf(conn, "Active Bucket: %s\n", activeBucket.Name)
	return nil
}
func handleExit(conn net.Conn) error {
	fmt.Println("Exiting...")
	os.Exit(0)
	return nil
}

func handleHelp(typeHelp int, conn net.Conn) error {

	if typeHelp == 0 {
		printHelpGlobal(conn)
	} else {
		printHelpBucket(conn)
	}
	return nil
}

func printHelpGlobal(conn net.Conn) {
	fmt.Fprintln(conn, "Available commands:")
	fmt.Fprintln(conn, "  create <bucket_name>         - Create a new bucket")
	fmt.Fprintln(conn, "  list                         - List all buckets")
	fmt.Fprintln(conn, "  use <bucket_name>            - Switch to the specified bucket")
	fmt.Fprintln(conn, "  drop <bucket_name>           - Delete the specified bucket")
	fmt.Fprintln(conn, "  pwb                          - Print the active bucket")
	fmt.Fprintln(conn, "  exit / quit                  - Exit the CLI")
	fmt.Fprintln(conn, "  help                         - Show this help message")
}
