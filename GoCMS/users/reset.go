package users

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/http"
	"net/smtp"
	"os"
)

// validTokens stores currently valid tokens
var validTokens []string

const passwordForm = `
  <h1>Enter your new password</h1>
  <form action="/reset" method="POST">
  <input type="hidden" name="email" value="{{ . }}" required>

  <label for="password">Password</label>
  <input type="password" name="password" required>
  <input type="submit" value="Submit">
  </form>
`

// SendPasswordResetEmail sends a password reset email to an email address.
func SendPasswordResetEmail(email string) error {
	token := string(genRandBytes())
	validTokens = append(validTokens, token)
	resetLink := "http://localhost:3000/reset?user=" + email + "&token=" + token

	username := os.Getenv("GMAIL_USERNAME")
	password := os.Getenv("GMAIL_PASSWORD")

	auth := smtp.PlainAuth("smtp.gmail.com:587", username, password, "smtp.gmail.com")
	return smtp.SendMail("smtp.gmail.com:587", auth, username,
		[]string{email}, []byte("Click here to reset your password"+resetLink))
}

// ResetPassword is an HTTP handler to handler password resets
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var isValid bool
		var index int

		// get query string params
		email := r.URL.Query().Get("user")
		token := r.URL.Query().Get("token")

		// verify token is a valid token
		for i, tok := range validTokens {
			if tok == token {
				isValid = true
				index = i
			}
		}

		if isValid != true {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Delete token from valid tokens
		validTokens = append(validTokens[:index], "")

		// render the reset template
		t, _ := template.New("password").Parse(passwordForm)
		t.Execute(w, email)
		return

	case "POST":
		//Get new password from form
		password := r.FormValue("password")
		email := r.FormValue("email")
		r.ParseForm()

		// Reset password by overwriting old password
		err := OverrideOldPassword(email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Successfully reset password"))
	}

}

// genRandBytes creates 32 byte strong of random bytes
func genRandBytes() []byte {
	b := make([]byte, 24)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return []byte(base64.URLEncoding.EncodeToString(b))
}
