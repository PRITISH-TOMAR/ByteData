package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/PRITISH-TOMAR/byted/constants"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var authFile = constants.AUTHFILEPATH

func InitAuthFile() error {
	dir := filepath.Dir(authFile)

	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create auth directory: %v", err)
		}
	}

	return nil
}

func AuthExists() bool {
	_, err := os.Stat(authFile)
	return !os.IsNotExist(err)

}

func CreateUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user := User{
		Username: username,
		Password: string(hash),
	}

	

	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %v", err)
	}

	if err := InitAuthFile(); err != nil {
		return err
	}
	
	return os.WriteFile(authFile, data, 0600)
}

func ValidateUser(username, password string) error {
	data, err := os.ReadFile(authFile)
	if err != nil {
		return fmt.Errorf("failed to read auth file: %v", err)
	}

	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return err
	}
	if user.Username != username {
		fmt.Println("Invalid username")
		return errors.New("invalid username")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println("Invalid password")
		return errors.New("invalid password")
	}

	return nil
}

func FirstTimeSetup() (string, string, error) {

	var newUser, newPass string

	fmt.Println("Welcome! Please set up your username and password.");
	fmt.Print("Enter new username: ")
	fmt.Scanln(&newUser)
	fmt.Print("Enter new password: ")
	fmt.Scanln(&newPass)

	if newUser == "" || newPass == "" {
		return "", "", errors.New("username and password cannot be empty")
	}

	if err:= CreateUser(newUser, newPass); err != nil {
		return "", "", err
	}

	fmt.Println("User created successfully!")
	return newUser, newPass, nil
}