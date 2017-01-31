package main

import (
	"GoCMS/users"
	"fmt"
	"html/template"
	"net/http"
)

const loginTemplate = `
<h1>Enter your username and password</h1>
<form action="/" method="POST">
	<input type="text" name="user" required>

	<label for="password">Password</label>
	<input type="password" name="password" required>

	<input type="submit" value="Submit">
</form>
`

func main() {
	username, password := "GMAIL_USERNAME", "GMAIL_PASSWORD"
	//username, password := os.Getenv("GMAIL_USERNAME"), os.Getenv("GMAIL_PASSWORD")

	err := users.NewUser(username, password)
	if err != nil {
		fmt.Printf("Couldn't create user %s\n", username, err.Error())
		return
	}

	fmt.Printf("Successfully created and authenticated user %s\n", username)

	http.HandleFunc("/", authHandler)
	http.HandleFunc("/restricted", restrictedHandler)
	http.HandleFunc("/reset", users.ResetPassword)
	http.ListenAndServe(":3000", nil)
}

// authHandler handles user authorization and session setting
// based on user form input.
func authHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t, _ := template.New("login").Parse(loginTemplate)
		t.Execute(w, nil)
	case "POST":
		u := r.FormValue("user")
		p := r.FormValue("password")
		err := users.AuthenticateUser(u, p)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		users.SetSession(w, u)
		w.Write([]byte("Signed in successfully"))
	}
}

// restrictedHandler verifies a session was set for a user to prevent
// users from accessing restricted pages.
func restrictedHandler(w http.ResponseWriter, r *http.Request) {
	user := users.GetSession(w, r)
	w.Write([]byte(user))
}
