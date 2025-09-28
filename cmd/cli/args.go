package cli

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"golang.org/x/term"
	"errors"
)

// ParseLoginArgs extracts username from CLI args and prompts for password
func ParseLoginArgs() (username, password string, err error) {
	args := os.Args[1:]

	// Parse username from arguments only
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-u":
			if i+1 < len(args) {
				username = args[i+1]
				i++
			}

		default:
			if strings.HasPrefix(args[i], "-") {
				err = fmt.Errorf("unsupported flag: %s. Only -u flag is supported", args[i])
				return "", "", err
			}

		}
	}

	// Check if username was provided
	if username == "" {
		err = errors.New("username is required. Usage: bytedata -u <username>")
		return "", "", err
	}

	// Prompt for password
	fmt.Printf("Password for user '%s': ", username)
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", fmt.Errorf("failed to read password: %v", err)
	}

	password = strings.TrimSpace(string(passwordBytes))

	// Validate password is not empty
	if password == "" {
		err = errors.New("password cannot be empty")
		return "", "", err
	}

	return username, password, nil
}
