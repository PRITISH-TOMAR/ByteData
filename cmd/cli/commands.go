package cli

import (
	"byted/internal/bucket"
	"errors"
	"fmt"
	"net"
	"strings"
)

func ExecuteCommand(input string, bucket *bucket.Bucket, conn net.Conn) error {

	fmt.Fprintf(conn, "Bytedata: [%s]> ", bucket.Name)
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}
	parts := strings.Fields(input)
	command := parts[0]

	switch command {
	case "put":
		return handlePut(parts, bucket, conn)
	case "get":
		return handleGet(parts, bucket, conn)
	case "del", "delete":
		return handleDelete(parts, bucket, conn)
	case "range":
		return handleRange(parts, bucket, conn)
	case "exit", "quit":
		return handleExitForBucket(conn)
	case "help":
		return handleHelp(1, conn)
	default:
		return handleUnknown(command, conn)
	}

}

func handlePut(parts []string, bucket *bucket.Bucket, conn net.Conn) error {
	if len(parts) < 3 {
		return errors.New("usage: put <key> <value>")
	}

	key := parts[1]
	value := strings.Join(parts[2:], " ")

	lsn, err := bucket.KvEngine.Put([]byte(key), []byte(value))
	if err != nil {
		return fmt.Errorf("put failed: %v", err)
	}

	fmt.Fprintf(conn, "Put successful. LSN: %d\n", lsn)
	return nil
}

func handleGet(parts []string, bucket *bucket.Bucket, conn net.Conn) error {
	if len(parts) != 2 {
		return errors.New("usage: get <key>")
	}

	key := parts[1]
	value, err := bucket.KvEngine.Get([]byte(key))
	if err != nil {
		return fmt.Errorf("get failed: %v", err)
	}

	if value == nil {
		fmt.Println("Key not found")
	} else {
		fmt.Fprintf(conn, "Value: %s\n", string(value))
	}

	return nil
}

func handleDelete(parts []string, bucket *bucket.Bucket, conn net.Conn) error {
	if len(parts) != 2 {
		return errors.New("usage: del <key>")
	}

	key := parts[1]
	lsn, err := bucket.KvEngine.Delete([]byte(key))
	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}

	fmt.Fprintf(conn, "Delete successful. LSN: %d\n", lsn)
	return nil
}

func handleRange(parts []string, bucket *bucket.Bucket, conn net.Conn) error {
	if len(parts) != 3 {
		return errors.New("usage: range <startKey> <endKey>")
	}

	startKey := parts[1]
	endKey := parts[2]
	results := bucket.KvEngine.Range([]byte(startKey), []byte(endKey))

	if len(results) == 0 {
		fmt.Println("No keys found in the specified range")
	} else {
		fmt.Fprintf(conn, "Found %d key(s):\n", len(results))
		for _, v := range results {
			fmt.Fprintf(conn, "  Key: %s, Value: %s\n", v.Key, v.Value)
		}
	}

	return nil
}

func handleUnknown(command string, conn net.Conn) error {
	return fmt.Errorf("unknown command: '%s'. Type 'help' for available commands", command)
}

func printHelpBucket(conn net.Conn) {
	fmt.Println("Available commands:")
	fmt.Println("  put <key> <value>    - Add or update a key-value pair")
	fmt.Println("  get <key>            - Retrieve the value for a given key")
	fmt.Println("  del <key>            - Delete a key-value pair")
	fmt.Println("  range <start> <end>  - Retrieve all key-value pairs in the specified key range")
	fmt.Println("  exit                 - Exit the CLI")
	fmt.Println("  help                 - Show this help message")
}

func handleExitForBucket(conn net.Conn) error {
	fmt.Println("Exiting bucket CLI...")
	return errors.New("exit")
}
