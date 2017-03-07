package users

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
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
	verifyKey   *rsa.PublicKey
)

// New creates a new Oauth2 configuration.
func New() *oauth2.Config {
	return &oauth2.Config{
		// ClientID:     os.Getenv("GOOGLE_KEY"),
		// ClientSecret: os.Getenv("GOOGLE_SECRET"),

		ClientID:     os.Getenv("GMAIL_USR"),
		ClientSecret: os.Getenv("GMAIL_PWD"),
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
	//genToken(w, email)
	genToken(w, []byte(email))

}

// genToken generates a Web Token and writes a JSON reponse.
//func genToken(w http.ResponseWriter, user string) {
func genToken(w http.ResponseWriter, user []byte) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Reserved claims: Some of them are: iss (issuer), exp (expiration time), sub (subject), aud (audience), and others.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// claims := make(jwt.MapClaims)
		claims["sub"] = user
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		claims["iat"] = time.Now().Unix()

		fmt.Println(claims["sub"], claims["exp"], claims["iat"])
	} else {
		// fmt.Println(err)
	}
	// claims := make(jwt.MapClaims)
	// claims["sub"] = user
	// claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// claims["iat"] = time.Now().Unix()

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
	//var claims jwt.Claims
	//claims := jwt.MapClaims{}
	//func VerifyToken(r *http.Request) (*rsa.PublicKey, error) {
	// Extract and parses a JWT token from an HTTP request.
	// Accepts a request and an extractor interface to define
	// the token extraction logic.

	//type Keyfunc func(*Token) (interface{}, error)
	// Parse methods use this callback function to supply
	// the key for verification.  The function receives the parsed,
	// but unverified Token.  This allows you to use properties in the
	// Header of the token (such as `kid`) to identify which key to use.

	// MapClaims is an alias for map[string]interface{} with built in validation behavior. Must type cast the claims property.
	// !! REQUEST IS CURRENTLY NIL !!
	fmt.Println("REQUEST", r)
	//token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &CustomClaimsExample{}, func(token *jwt.Token) (interface{}, error) {
	// Read the token out of the response body
	//buf := new(bytes.Buffer)
	//io.Copy(buf, r.Body)
	//r.Body.Close()
	//tokenString := strings.TrimSpace(buf.String())

	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, keyLookupFunc)
	//token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, keyLookupFuncear

	if err == nil {
		fmt.Println("PARSED TOKEN FROM REQUEST")
		//fmt.Printf("Token for SUB %v EXP %v IAT %v", claims["sub"], claims["exp"], claims["iat"])
	} else {
		fmt.Println(err)
	}

	// !! TOKEN IS CURRENTLY NIL !!
	fmt.Println(token)
	// Lookup and unpack key from PEM encoded PKCS8
	//	_, err := jwt.ParseRSAPublicKeyFromPEM(claims["sub"].([]byte))
	// check(err)

	// !! CLAIMS IS CURRENTLY NIL !!
	fmt.Println("TOKEN", token)
	return "TOKEN TEST", err
	//return claims["sub"].(string), err
}

// Validate the algorithm is the expected algorithm
/* Parse methods use this callback function to supply the key for verification. The function receives the parsed, but unverified Token. This allows you to use properties in the Header of the token (such as `kid`) to identify which key to use. */
func keyLookupFunc(t *jwt.Token) (interface{}, error) {
	fmt.Println("CALLING THE KEY LOOKUP FUNCTION")
	fmt.Println("TOKEN ARGUMENT: ", t)
	// // Look up key
	// key, err := lookupPublicKey(t.Header["sub"])
	// if err != nil {
	// 	return nil, err
	// }

	// Unpack key from PEM encoded PKCS8
	// return jwt.ParseRSAPublicKeyFromPEM(key)
	//return jwt.ParseRSAPublicKeyFromPEM(t.Header["sub"].([]byte))
	return verifyKey, nil
	// if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
	// 	return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
	// }
	// return t.Method.(*jwt.SigningMethodHMAC), nil

}
