package main

import (
	"GoCMS/users"
	"fmt"
	"net/http"
)

func main() {
/username, password := os.Getenv("GMAIL_USERNAME"), os.Getenv("GMAIL_PASSWORD")

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

	// send reset email
	err = users.SendPasswordResetEmail(username)
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/reset", users.ResetPassword)

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
