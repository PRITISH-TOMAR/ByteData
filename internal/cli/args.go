package cli

import "os"

// ParseLoginArgs extracts username and password from CLI args
func ParseLoginArgs() (username, password string) {
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-u":
			if i+1 < len(args) {
				username = args[i+1]
				i++
			}
		case "-p":
			if i+1 < len(args) {
				password = args[i+1]
				i++
			}
		}
	}
	return
}
