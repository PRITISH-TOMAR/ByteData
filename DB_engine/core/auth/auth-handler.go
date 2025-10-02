package auth

import (
	"byted/DB_engine/structs"
	"fmt"
)

func HandleAuthenticatedConnection(comm *structs.Communicators) bool {
	dec := comm.Dec
	enc := comm.Enc
	enc.Encode(structs.Message{Type: "request", Field: "username"})

	var msg structs.Message
	
	dec.Decode(&msg)

	if !AuthExists() {
		enc.Encode(structs.Message{Type: "info", Message: fmt.Sprintf("Hey %s, it looks like your first login, let's setup account", msg.Username)})

		enc.Encode(structs.Message{Type: "request", Field: "password", Message: fmt.Sprintf("\nEnter %s's password: ", msg.Username)})

		dec.Decode(&msg)
		if msg.Username== "" || msg.Password == ""{
			enc.Encode(structs.Message{Type: "error", Message: "\nUsername & password must not be empty!"})
		}

		if CreateUser(msg.Username, msg.Password) != nil {
			enc.Encode(structs.Message{Type: "error", Message: "\nNew user creation failure!"})
		}
		enc.Encode(structs.Message{Type: "success", Message: fmt.Sprintf("\nWelcome %s", msg.Username)});
		return true
	}

	if !UserExists(msg.Username) {
		enc.Encode(structs.Message{Type: "error", Message: "\nUser not found"})
		return false
	}
	enc.Encode(structs.Message{Type: "request", Field: "password", Message: fmt.Sprintf("Enter %s's password: ", msg.Username)})

	if err := dec.Decode(&msg); err != nil || msg.Password == "" {
		enc.Encode(structs.Message{Type: "error", Message: "\nInvalid password"})
		return true
	}

	if err := ValidateUser(msg.Username, msg.Password); err != nil {
		enc.Encode(structs.Message{Type: "error", Message: "\nAuthentication failed"})
		return false
	}

	enc.Encode(structs.Message{Type: "success", Message: "\nLogging in...."})
	return true
}
