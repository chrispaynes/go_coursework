package main

import (
	"GoCMS/users"
	"fmt"
)

func main() {
	username, password := "gopher", "archLinux123"

	err := users.NewUser(username, password)
	if err != nil {
		fmt.Printf("Couldn't create user %s\n", username, err.Error())
		return
	}

	err = users.AuthenticateUser(username, password)
	if err != nil {
		fmt.Printf("Couldn't authenticate user %s\n %s", username, err.Error())
	}

	fmt.Printf("Successfully created and authenticated user %s\n", username)

	err = users.NewUser(username, password)
	if err != nil {
		fmt.Printf("Couldn't create user %s\n %s", username, err.Error())
		return
	}

	err = users.AuthenticateUser(username, password)
	if err != nil {
		fmt.Printf("Couldn't authenticate user %s\n %s", username, err.Error())
	}

}
