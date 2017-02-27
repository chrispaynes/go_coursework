package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
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
		ClientID:     os.Getenv("GOOGLE_KEY"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
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
	//genToken(w, email)
	genToken(w, []byte(email))

}

// genToken generates a Web Token and writes a JSON reponse.
//func genToken(w http.ResponseWriter, user string) {
func genToken(w http.ResponseWriter, user []byte) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["sub"] = user
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["iat"] = time.Now().Unix()
	//token.Claims["sub"] = user
	//token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	//token.Claims["iat"] = time.Now().Unix()

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
	//func VerifyToken(r *http.Request) (*rsa.PublicKey, error) {
	// Extract and parses a JWT token from an HTTP request.
	// Accepts a request and an extractor interface to define
	// the token extraction logic.

	//type Keyfunc func(*Token) (interface{}, error)
	// Parse methods use this callback function to supply
	// the key for verification.  The function receives the parsed,
	// but unverified Token.  This allows you to use properties in the
	// Header of the token (such as `kid`) to identify which key to use.

	token, _ := request.ParseFromRequest(r, request.OAuth2Extractor, keyLookupFunc)

	// MapClaims is an alias for map[string]interface{} with built in validation behavior. Must type cast the claims property.
	claims := token.Claims.(jwt.MapClaims)
	fmt.Printf("Token for user %v expires %v", claims["user"], claims["exp"])

	// Lookup and unpack key from PEM encoded PKCS8
	key, err := jwt.ParseRSAPublicKeyFromPEM(claims["sub"].([]byte))
	if err != nil {
		return "", err
	}

	return key.N.String(), err
}

func keyLookupFunc(t *jwt.Token) (interface{}, error) {
	// Validate the algorithm is the expected algorithm
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
	}

	return t, nil
}
