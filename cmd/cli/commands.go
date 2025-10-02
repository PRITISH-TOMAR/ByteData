package cli

import (
	"byted/internal/bucket"
	"errors"
	"fmt"
	"strings"
)

func ExecuteCommand(input string, bucket *bucket.Bucket, bm * bucket.BucketManager) ([]string, error) {

	fmt.Printf("Bytedata: [%s]> ", bucket.Name)
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, nil
	}

	parts := strings.Fields(input)
	command := parts[0]

	switch command {
	case "put":
		return handlePut(parts, bucket)
	case "get":
		return handleGet(parts, bucket)
	case "del", "delete":
		return handleDelete(parts, bucket)
	case "range":
		return handleRange(parts, bucket)
	case "exit", "quit":
		return handleExitForBucket(bm)
	case "help":
		return handleHelp(1)
	default:
		return handleUnknown(command)
	}
	return nil, nil

}

func handlePut(parts []string, bucket *bucket.Bucket) ([]string, error) {
	if len(parts) < 3 {
		return nil, errors.New("usage: put <key> <value>")
	}

	key := parts[1]
	value := strings.Join(parts[2:], " ")

	lsn, err := bucket.KvEngine.Put([]byte(key), []byte(value))
	if err != nil {
		return nil, fmt.Errorf("put failed: %v", err)
	}

	return []string{fmt.Sprint("Put successful. LSN: %d\n", lsn)}, nil
}

func handleGet(parts []string, bucket *bucket.Bucket) ([]string, error) {
	if len(parts) != 2 {
		return nil, errors.New("usage: get <key>")
	}

	key := parts[1]
	value, err := bucket.KvEngine.Get([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("get failed: %v", err)
	}

	if value == nil {
		return nil, fmt.Errorf("Key not found")
	} 

	return []string{fmt.Sprintf("Value: %s\n", string(value))}, nil
}

func handleDelete(parts []string, bucket *bucket.Bucket) ([]string, error) {
	if len(parts) != 2 {
		return nil, errors.New("usage: del <key>")
	}

	key := parts[1]
	lsn, err := bucket.KvEngine.Delete([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("delete failed: %v", err)
	}

	return []string{ fmt.Sprintf("Delete successful. LSN: %d\n", lsn)}, nil
}

func handleRange(parts []string, bucket *bucket.Bucket) ([]string, error) {
	if len(parts) != 3 {
		return nil, errors.New("usage: range <startKey> <endKey>")
	}

	startKey := parts[1]
	endKey := parts[2]
	results := bucket.KvEngine.Range([]byte(startKey), []byte(endKey))

	if len(results) == 0 {
		return nil, errors.New("no keys found in the specified range")
	}

	// Build a slice of strings to return
	response := []string{fmt.Sprintf("Found %d key(s):", len(results))}
	for _, v := range results {
		response = append(response, fmt.Sprintf("  Key: %s, Value: %s", v.Key, v.Value))
	}

	return response, nil
}


func handleUnknown(command string) ([]string, error) {
	return nil, fmt.Errorf("unknown command: '%s'. Type 'help' for available commands", command)
}

func printHelpBucket() ([]string, error) {
	help := []string{`Available commands:
  put <key> <value>    - Add or update a key-value pair
  get <key>            - Retrieve the value for a given key
  del <key>            - Delete a key-value pair
  range <start> <end>  - Retrieve all key-value pairs in the specified key range
  exit                 - Exit the CLI
  help                 - Show this help message`}
	return help, nil
}

func handleExitForBucket(bm * bucket.BucketManager) ([]string, error) {
	bm.ExitBucket()
	return []string{fmt.Sprintf("Exiting.")}, nil
}
