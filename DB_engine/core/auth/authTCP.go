package auth

import (
	"bufio"
	"fmt"
	"net"
	"strings"

)

// AuthTCP reads username/password from a TCP connection
func AuthTCP(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	if !AuthExists() {
		username, password, err := FirstTimeSetup( conn, reader) 
		
		if err != nil {
			return "", err
		}

		if err := CreateUser(username, password); err != nil {
			return "", err
		}
		return username, nil
	}

	// Ask for username/password over TCP
	fmt.Fprint(conn, "Username: ")
	username, _ := reader.ReadString('\n')
	fmt.Fprint(conn, "Password: ")
	username = strings.TrimSpace(username)

	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// 3️⃣ Validate credentials
	if err := ValidateUser(username, password); err != nil {
		return "", err
	}

	return username, nil
}
