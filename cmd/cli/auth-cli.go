package cli

import (
	"fmt"
	"github.com/PRITISH-TOMAR/byted/internal/auth"
)


func AuthCLI() (string, error) {
		var username, password string
   
	if !auth.AuthExists() {
		var err error
		username, password, err = auth.FirstTimeSetup()
		if err != nil {
			fmt.Println("Error during setup:", err)
			return "", err
		}
	} else {
		var err error
		username, password, err = ParseLoginArgs()
		if err != nil {
			fmt.Println(err)
			return "", err
		}
	}

	//  2. Validate user credentials
	if err := auth.ValidateUser(username, password); err != nil {
		return "", err
	}

	return username, nil
}