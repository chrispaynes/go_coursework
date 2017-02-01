package users

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	identityURL = "http://www.googleapis.com/oauth2/v2/userinfo"
	provider    = New()
	signingKey  = genRandBytes()
)

// New creates a new Oauth2 configuration.
func New() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("Google_Key"),
		ClientSecret: os.Getenv("google_Secret"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:3000/auth/gplus/callback",
		Scopes:       []string{"email", "profile"},
	}
}

// AuthURLHandler redirects user to Oauth sign-in page for a given provider.
func AuthURLHandler(w http.ResponseWriter, r *http.Request) {
	authURL := provider.AuthCodeURL("", oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func CallbackURLHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := provider.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	client := provider.Client(oauth2.NoContext, token)
	resp, err := client.Get(identityURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	defer resp.Body.Close()

	user := make(map[string]string)
	json.NewDecoder(resp.Body).Decode(&user)

	email := user["email"]
	genToken(w, email)

}

// genToken generates a Web Token and writes a JSON reponse.
func genToken(w http.ResponseWriter, user string) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["sub"] = user
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	token.Claims["iat"] = time.Now().Unix()

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\n\ttoken: " + tokenString + "\n}"))
}

// VerifyToken receives token from HTTP request, verifies the token is valid
// and returns the username.
func VerifyToken(r *http.Request) (string, error) {
	token, err := jwt.Parse(r, func(token, *jwt.Token) (interface{}, error) {
		claims := token.Claims.(jwt.MapClaims)
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return signingKey, nil
	})
	if err != nil {
		return "", err
	}

	if token.Valid == false {
		return "", jwt.ErrInvalidKey
	}
	return token.Claims["sub"].(string), nil
}
