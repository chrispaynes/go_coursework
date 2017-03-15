package main

import (
	"GoCMS/users"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
	http.HandleFunc("/", authHandler)
	http.HandleFunc("/reset", users.ResetPassword)
	http.HandleFunc("/auth/gplus/authorize", users.AuthURLHandler)
	http.HandleFunc("/auth/gplus/callback", users.CallbackURLHandler)
	http.HandleFunc("/oauth", oauthRestrictedHandler)
	http.HandleFunc("/restricted", restrictedHandler)

	http.ListenAndServe(":3000", nil)
}

func initUser() {
	f, err := ioutil.ReadFile("../secret/user_config.txt")
	check(err)

	os.Setenv("GMAIL_USR", strings.Split(string(f), "\n")[0])
	os.Setenv("GMAIL_PWD", strings.Split(string(f), "\n")[1])
	usr, pwd := os.Getenv("GMAIL_USR"), os.Getenv("GMAIL_PWD")

	fmt.Printf("Attempting to authenticating user: %s using password: %s\n", usr, pwd)

	err = users.NewUser(usr, pwd)
	if err != nil {
		fmt.Printf("Couldn't create user %s\n%s", usr, err.Error())
		return
	}

	fmt.Printf("Successfully created and authenticated user %s\n", usr)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func oauthRestrictedHandler(w http.ResponseWriter, r *http.Request) {
	user, err := users.VerifyToken(r)
	fmt.Println("oauthRestrictedHandler ERROR:", err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	users.GenToken(w, []byte(user))

	w.Write([]byte(user))
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
