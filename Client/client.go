package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"golang.org/x/term"
)

// Message matches the JSON protocol used by the server
type Message struct {
	Type     string   `json:"type"`
	Bucket   string   `json:"bucket"`
	Field    string   `json:"field,omitempty"`
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	Command  string   `json:"command,omitempty"`
	Message  string   `json:"message,omitempty"`
	Data     []string `json:"data,omitempty"`
}

func main() {
	// CLI flags
	username, addr := getClientInfo()

	// Connect to server
	conn := getConnection(addr)
	defer conn.Close()
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	var msg Message

	// Handle entire authentication
	if !handleAuth(username, enc, dec, &msg) {
		return
	}

	// handle command mode
	CommandLoop(enc, dec, msg)
}

func getClientInfo() (*string, *string) {
	uname := flag.String("u", "", "Username for the session")
	addr := flag.String("addr", "localhost:8080", "Server address")
	flag.Parse()

	if *uname == "" {
		fmt.Println("Usage: bytedata -u <username>")
		os.Exit(1)
	}
	return uname, addr
}

func getConnection(addr *string) net.Conn {
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Println("Failed to connect:", err)
		os.Exit(1)
	}

	return conn
}

func handleAuth(username *string, enc *json.Encoder, dec *json.Decoder, msg *Message) bool {
	if err := dec.Decode(&msg); err != nil {
		fmt.Println("Connection closed by server!", err)
		return false
	}
	enc.Encode(Message{Type: "auth", Username: *username})

	// Step 3: Read server reply (user exists or error)
	if err := dec.Decode(&msg); err != nil {
		fmt.Println("Connection closed by server!")
		return false
	}
	if msg.Type == "error" {
		fmt.Println(msg.Message)
		return false
	}

	// Step 4: Server asks for password
	if msg.Type == "info" {
		fmt.Print(msg.Message)
		dec.Decode(&msg)
	}

	if msg.Type == "request" && msg.Field == "password" {
		fmt.Print(msg.Message)
		passBytes, _ := term.ReadPassword(int(os.Stdin.Fd()))
		password := strings.TrimSpace(string(passBytes))

		enc.Encode(Message{Type: "auth", Password: password})

	}
	// Step 5: Read authentication result
	if err := dec.Decode(&msg); err != nil {
		fmt.Println("Connection closed by server!")
		return false
	}
	if msg.Type == "error" {
		fmt.Println(msg.Message)
		return false
	}
	fmt.Println(msg.Message) // Welcome message
	return true
}

func CommandLoop(enc *json.Encoder, dec *json.Decoder, msg Message) {
	reader := bufio.NewReader(os.Stdin)
	active := ""
	for {
		fmt.Print("ByteData> " + active)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == ""{
			continue
		}
		// Send command
		enc.Encode(Message{Type: "command", Command: line})

		// Read server response
		if err := dec.Decode(&msg); err != nil {
			fmt.Println("Exiting...")
			return
		}
		
		active = msg.Bucket
		if active != "" {
			active = "["+active+"]:"
		}

		if msg.Type == "error" {
			fmt.Println("ByteData> " + active + msg.Message)
		} else {
			fmt.Println("ByteData> " + active + strings.Join(msg.Data, "\n"))
		}
	}
}
