package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/PRITISH-TOMAR/byted/internal/kv"
)

func ExecuteCommmand(input string, engine *kv.KVEngine) error {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}
	parts := strings.Fields(input)
	command := parts[0]

	switch command {
	case "put":
		if len(parts) < 3 {
			return errors.New("usage: put <key> <value>")
		}
		key := parts[1]
		value := strings.Join(parts[2:], " ")

		lsn, err := engine.Put([]byte(key), []byte(value))

		if err != nil {
			return fmt.Errorf("put failed: %v", err)
		}
		fmt.Printf("Put successful. LSN: %d\n", lsn)

	case "get":
		if len(parts) != 2 {
			return errors.New("usage: get <key>")
		}
		key := parts[1]
		value, err := engine.Get([]byte(key))
		if err != nil {
			return fmt.Errorf("get failed: %v", err)
		}
		if value == nil {
			fmt.Println("Key not found")
		} else {
			fmt.Printf("Value: %s\n", string(value))
		}

	case "del":
		if len(parts) != 2 {
			return errors.New("usage: delete <key>")
		}
		key := parts[1]
		lsn, err := engine.Delete([]byte(key))
		if err != nil {
			return fmt.Errorf("delete failed: %v", err)
		}
		fmt.Printf("Delete successful. LSN: %d\n", lsn)

	case "range":
		if len(parts) != 3 {
			return errors.New("usage: range <startKey> <endKey>")
		}
		startKey := parts[1]
		endKey := parts[2]
		results := engine.Range([]byte(startKey), []byte(endKey))
		if len(results) == 0 {
			fmt.Println("No keys found in the specified range")
		} else {
			for _, v := range results {
				fmt.Printf("Key: %s, Value: %s\n", v.Key, v.Value)
			}
		}

	case "exit":
		fmt.Println("Exiting...")
		os.Exit(0)
		return errors.New("exit") // special error to signal exit

	case "help":
		fmt.Println("Available commands:")
		fmt.Println("  put <key> <value>    - Add or update a key-value pair")
		fmt.Println("  get <key>            - Retrieve the value for a given key")
		fmt.Println("  del <key>            - Delete a key-value pair")
		fmt.Println("  range <start> <end>  - Retrieve all key-value pairs in the specified key range")
		fmt.Println("  exit                 - Exit the CLI")
		fmt.Println("  help                 - Show this help message")
		return nil

	default:
		return errors.New("unknown command. Type 'help' for available commands")
	}
	return nil

}
