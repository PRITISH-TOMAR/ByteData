package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"golang.org/x/crypto/bcrypt"
)

type User struct{
	Username string `json:"username"`
	Password string `json:"password"`
}

const authFile = "/tmp/bytedatata/auth.json"

func InitRoot() error {
	
	if _, err := os.Stat(authFile); err ==nil{
		return nil
	}

	// create parent if missing
	if err:= os.MkdirAll("/tmp/bytedatata", 0700); err!=nil{
		return fmt.Errorf("failed to create auth parent dir: %v", err)
	}

	// hash default password
	hashed, err := bcrypt.GenerateFromPassword([]byte("root"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash default password: %v", err)
	}

	root := User{
		Username: "root",
		Password: string(hashed),
	}

	data, err := json.MarshalIndent([]User{root}, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal root user: %v", err)
	}

	// only owner can read/write
	if err := os.WriteFile(authFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write auth file: %v", err)
	}

	return nil
}