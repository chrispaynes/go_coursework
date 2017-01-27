package main

import (
	"GoCMS/users"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

//var DB, err = bolt.Open("user.db", 0600, nil)

func main() {
	&DB{DB}, err := bolt.Open("user.db", 0600, nil)
	if err != nil {
	log.Fatal(err)
		}
	defer db.Close()
	username, password := os.Getenv("GMAIL_USERNAME"), os.Getenv("GMAIL_PASSWORD")

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

func authHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t, _ := template.New("login").Parse(loginTemplate)
		t.Execute(w, nil)
	case "POST":
		user := r.FormValue("user")
		pass := r.FormValue("password")
		err := users.AuthenticateUser(user, pass)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users.SetSession(w, user)
		w.Write([]byte("Signed in successfully"))
	}
}

func restrictedHandler(w http.ResponseWriter, r *http.Request) {
	user := users.GetSession(w, r)
	w.Writer([]byte(user))
}
