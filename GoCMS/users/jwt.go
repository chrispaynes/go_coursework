package users

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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
	verifyKey   *rsa.PublicKey
)

// initClient initializes the client information for Oauth2
func initClient() {
	c, err := ioutil.ReadFile("../secret/client_config.txt")
	check(err)

	os.Setenv("ClientID", strings.Split(string(c), "\n")[0])
	os.Setenv("ClientSecret", strings.Split(string(c), "\n")[1])
}

// New creates a new Oauth2 configuration.
func New() *oauth2.Config {
	initClient()

	return &oauth2.Config{
		ClientID:     os.Getenv("ClientID"),
		ClientSecret: os.Getenv("ClientSecret"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:3000/auth/gplus/callback",
		Scopes:       []string{"email", "profile"},
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
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
	GenToken(w, []byte(email))

}

// genToken generates a Web Token and writes a JSON reponse.
func GenToken(w http.ResponseWriter, user []byte) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Reserved claims: Some of them are: iss (issuer), exp (expiration time), sub (subject), aud (audience), and others.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claims["sub"] = user
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		claims["iat"] = time.Now().Unix()
		token.Claims = claims

		fmt.Println(claims["sub"], claims["exp"], claims["iat"])
	} else {
		// fmt.Println(err)
	}

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
	fmt.Println("REQUEST", r)

	// Extract and parses a JWT token from an HTTP request.
	// Accepts a request and an extractor interface to define
	// the token extraction logic.
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, keyLookupFunc)

	if err == nil {
		fmt.Println("PARSED TOKEN FROM REQUEST")
		//fmt.Printf("Token for SUB %v EXP %v IAT %v", claims["sub"], claims["exp"], claims["iat"])
	} else {
		fmt.Println("DID NOT PARSE TOKEN", err)
	}

	fmt.Println(token)
	// Lookup and unpack key from PEM encoded PKCS8
	//	_, err := jwt.ParseRSAPublicKeyFromPEM(claims["sub"].([]byte))
	// check(err)

	fmt.Println("TOKEN", token)
	return "TOKEN", err
	//return claims["sub"].(string), err
}

// Validate the algorithm is the expected algorithm
// The function receives the parsed,
// but unverified Token.  This allows you to use properties in the
// Header of the token (such as `kid`) to identify which key to use.
func keyLookupFunc(t *jwt.Token) (interface{}, error) {
	fmt.Println("CALLING THE KEY LOOKUP FUNCTION")
	fmt.Println("TOKEN ARGUMENT: ", t)

	return t.Method.(*jwt.SigningMethodHMAC), nil

}
