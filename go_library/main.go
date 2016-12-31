package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/codegangsta/negroni"
	gmux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yosssi/ace"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Book represents a book object
type Book struct {
	PK             int64  `db:"pk"`
	Title          string `db:"title"`
	Author         string `db:"author"`
	Classification string `db:"classification"`
	ID             string `db:"id"`
	User           string `db:"user"`
}

// ClassifySearchResponse represents the
// book collection returned from a search.
type ClassifySearchResponse struct {
	Results []SearchResult `xml:"works>work"`
}

// ClassifyBookResponse represented nested structs
// to parse XML node hierarchy.
type ClassifyBookResponse struct {
	BookData struct {
		Title  string `xml:"title,attr"`
		Author string `xml:"author,attr"`
		ID     string `xml:"owi,attr"`
	} `xml:"work"`
	Classification struct {
		MostPopular string `xml:"sfa,attr"`
	} `xml:"recommendations>ddc>mostPopular"`
}

// LoginPage represents the data reported to the Login Page.
type LoginPage struct {
	Error string
}

// A Page represents a book collection.
type Page struct {
	Books  []Book
	Filter string
	User   string
}

// A SearchResult represents a book-related query result.
type SearchResult struct {
	Title  string `xml:"title,attr"`
	Author string `xml:"author,attr"`
	Year   string `xml:"hyr,attr"`
	ID     string `xml:"owi,attr"`
}

// A User represents a unique app user with authentication credentials.
type User struct {
	Username string `db:"username"`
	Secret   []byte `db:"secret"`
}

// Db stores a global sqlite3 database.
var db *sql.DB

// page stores a newly instantiated Page.
var page = new(Page)

// store initializes a new CookieStore with a secret authentication key.
var cookieStore = sessions.NewCookieStore([]byte("password123"))

const store = "go_library"

var bindVar string

func initDB() {
	// Os.Getenv instructs the app to behave a certain way depending on
	// whether the app is in development or production mode.
	if os.Getenv("ENV") != "production" {
		// DB opens a sqlite3 database connection.
		db, _ = sql.Open("sqlite3", "dev.db")
		// &bindVar = "?"
	} else {
		db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		// &bindVar = fmt.Sprintf("$%d", 1)
	}

	// If the Books database table does not exist, create one from Book struct.
	db.Exec(`CREATE TABLE IF NOT EXISTS books(
    pk integer primary key autoincrement,
    title text,
    author text,
    id text,
    classification text)`)

	// If the Users table does not exist, create one from User struct.
	db.Exec(`CREATE TABLE IF NOT EXISTS users(
    username text primary key,
    secret blob
    )`)
}

// PingDB checks DB connectivity then calls next HandlerFunc upon connectivity.
func pingDB(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if err := db.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	next(w, r)
}

// setWhereClause returns a WHERE SQL clause using a filter and username;
func setWhereClause(f, u string) string {
	// wc stores the where clause or defaults to no clause.
	wc := fmt.Sprintf(" WHERE user=%q", u)
	if f == "fiction" {
		wc += " AND classification BETWEEN '800' AND '900' "
	} else if f == "nonfiction" {
		wc += " AND classification NOT BETWEEN '800' AND '900' "
	} else {
		wc += " "
	}

	return wc
}

// setCollection maps each DB row to a book and appends each book to a collection.
func setCollection(w http.ResponseWriter, rs *sql.Rows) {
	defer rs.Close()

	page.Books = []Book{}

	for rs.Next() {
		var b Book
		rs.Scan(&b.PK, &b.Title, &b.Author, &b.Classification)

		page.Books = append(page.Books, b)
	}

}

// queryBooks queries for books and returns the results sorted by name.
func queryBooks(where, sort string, w http.ResponseWriter) (*sql.Rows, error) {
	rows, err := db.Query("SELECT pk, title, author, classification FROM books" + where + "ORDER BY " + sort)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)

	}

	return rows, err
}

// getBookCollection returns a filtered and sorted book collection.
func getBookCollection(sort, filter, username string, w http.ResponseWriter, r *http.Request) bool {
	page.Filter = getSessionString(r, "filter")

	// Sort by "pk" when sort is not specified.
	if sort == "" {
		sort = "pk"
	}

	where := setWhereClause(filter, username)
	rows, _ := queryBooks(where, sort, w)

	setCollection(w, rows)

	return true
}

func startSession(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	cookieStore.Get(r, store)
	next(w, r)
}

// writeToSession gets a pre-existing session or creates a new session
// and sets key/value pair data.
func writeToSession(w http.ResponseWriter, r *http.Request, k string, v string) {
	s, _ := cookieStore.Get(r, store)
	s.Values[k] = v
	s.Save(r, w)
}

// GetSessionString returns a value from the session store in string format
func getSessionString(r *http.Request, key string) string {
	var strVal string
	var s, _ = cookieStore.Get(r, store)
	if val := s.Values[key]; val != nil {
		strVal = val.(string)
	}
	return strVal
}

// verifyUser is middleware handler that ensures user's session value is set before entering main page
func verifyUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// if currently on "/login", return from the func to avoid a http.Redirect Loop
	if r.URL.Path == "/login" {
		next(w, r)
		return
	}

	// verify username held in session store exists in DB
	// Username stores the username value stored in the session string
	// if the username is not empty will query the db for it
	// gets session string, if its blank, query username from form value to see if in DB
	// redirect to login if not logged in
	if username := getSessionString(r, "user"); username != "" {
		// db.Query searches the Database for a username.
		if user, _ := db.Query("SELECT username FROM users WHERE username=" + "'" + username + "'"); user != nil {
		} else {
			// http.Redirect redirects unregistered users to the login page
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}
	next(w, r)

}

func main() {
	initDB()

	// Mux replaces the default ServeMux with Gorilla/Mux Router.
	mux := gmux.NewRouter()

	// Mux.HandleFunc("/login") handles user authentication
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// Lp stores the login page to hold an error for the login page to display
		lp := new(LoginPage)

		// Creates and stores a new user upon registration, then redirects on successful
		if r.FormValue("register") != "" {
			// Secret stores a bcrypt password hash
			secret, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)

			// User stores a username and password hash
			user := User{r.FormValue("username"), secret}

			// db.Exec inserts new user into a new database row, or fails and return to login page.
			if _, err := db.Exec("insert into users (username, secret) values(?, ?)", &user.Username, &user.Secret); err != nil {
				lp.Error = err.Error()
			} else {
				// Session stores the user in the session store
				writeToSession(w, r, "user", user.Username)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			// Handles user authentication for registered users
		} else if r.FormValue("login") != "" {
			// User will store a username and password hash
			user := new(User)

			// dbQueryRow searches the Database for a user based on primary key,
			// then scans the query results, placing each column into the user object.
			err := db.QueryRow("SELECT * FROM users WHERE username="+"'"+r.FormValue("username")+"'").Scan(&user.Username, &user.Secret)

			// sql.ErrNoRows indicates the queried user does not exist in the database
			if err == sql.ErrNoRows {
				lp.Error = "No user found matching Username: " + r.FormValue("username")
			} else {
				// bcrypt.CompareHashAndPassword verifies the hashed password matches the plaintext password.
				if err = bcrypt.CompareHashAndPassword(user.Secret, []byte(r.FormValue("password"))); err != nil {
					lp.Error = "Invalid Password"
				} else {
					// Session stores the user in the session store
					writeToSession(w, r, "user", user.Username)
					http.Redirect(w, r, "/", http.StatusFound)
					return
				}
			}
		}

		// Template stores and caches an Ace Template loaded.
		template, err := ace.Load("templates/login", "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// writes HTTP response using template and displays p or error
		// Also sends an error to the login page when present.
		if err := template.Execute(w, lp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		// remove from session store by setting user value to nil
		writeToSession(w, r, "user", "")
		// remove default filter setting
		writeToSession(w, r, "filter", "")
		//redirects user back to login page
		http.Redirect(w, r, "/login", http.StatusFound)

	})

	mux.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {

		// Sb stores user's "sortBy" preference then validates the query value
		// to prevent SQL injection before storing it to the session.
		sb := r.FormValue("sortBy")
		if sb != "title" && sb != "author" && sb != "classification" {
			http.Error(w, "Invalid Column Name", http.StatusBadRequest)
			return
		}

		// Session stores the current URL's query's key value into the store.
		writeToSession(w, r, "sortBy", sb)

		// GetBookCollection returns a sorted book collection, then encodes it in JSON.
		getBookCollection(getSessionString(r, "sortBy"), getSessionString(r, "filter"), getSessionString(r, "user"), w, r)
		if err := json.NewEncoder(w).Encode(page.Books); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// Methods.Queries matches for "?key=" URL queries
		// that match the regexp patter for "key=value1, key=value2, key=valueN".
	}).Methods("GET").Queries("sortBy", "{sortBy:title|author|classification}")

	mux.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {

		// Session stores the current URL's query's key value into the store.
		writeToSession(w, r, "filter", r.FormValue("filter"))

		// GetBookCollection returns a sorted book collection, then encodes it in JSON.
		getBookCollection(getSessionString(r, "sortBy"), getSessionString(r, "filter"), getSessionString(r, "user"), w, r)

		if err := json.NewEncoder(w).Encode(page.Books); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// Methods.Queries matches for "?key=" URL queries
		// that match the regexp patter for "key=value1, key=value2, key=valueN".
	}).Methods("GET").Queries("filter", "{filter:all|fiction|nonfiction}")

	// Mux.HandleFunc("/") registers a handler function for requests on "/"
	// http.ResponseWriter writes an HTTP response
	// http.Request receives an HTTP server request or a client request
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// sets cookie expiration date and user info
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "username", Value: "gopher", Expires: expiration}
		http.SetCookie(w, &cookie)

		// sets sort value to sort value stored in session
		// sortBy, _ := cookieStore.Get(r, "go_library")
		// sortColumn := sortBy.Values["sortBy"].(string)

		getBookCollection(getSessionString(r, "sortBy"), getSessionString(r, "filter"), getSessionString(r, "user"), w, r)

		// page.User stores the logged in username
		page.User = getSessionString(r, "user")

		// Template stores and caches an Ace Template loaded.
		template, err := ace.Load("templates/index", "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// writes HTTP response using template and displays p or error
		if err := template.Execute(w, page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("GET")

	// HandleFunc() registers handler function for requests on "/search"
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		var results []SearchResult
		var err error

		// queries OCLC Book API using searchbar text or writes HTTP error
		if results, err = search(r.FormValue("search")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// encodes SearchResult data into JSON format
		encoder := json.NewEncoder(w)
		// writes output to the response?
		if err := encoder.Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("POST")

	// HandleFunc("/books") saves search results with URL of /books/add uses find() to search
	// OCLC Book API for a book's OCLC Work Identifier (OWI)
	mux.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		var book ClassifyBookResponse
		var err error

		// Book stores a object extracted from the requested URL's query string id value.
		// Book stores a {Book.Title, Book.Author, Book.Classifiction}.
		if book, err = find(r.FormValue("id")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// B uses query values and session values to store the lastest book,
		// prior to it's database insertion.
		b := Book{
			PK:             -1,
			Title:          book.BookData.Title,
			Author:         book.BookData.Author,
			Classification: book.Classification.MostPopular,
			ID:             r.FormValue("id"),
			User:           getSessionString(r, "user"),
		}

		// Db.exec inserts a book into the database, and
		// passes a nil primary key to allow DB to auto increment the book records.
		if _, err := db.Exec("insert into books (pk, title, author, id, classification, user) values(?, ?, ?, ?, ?, ?)",
			nil, &b.Title, &b.Author, &b.ID, &b.Classification, &b.User); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// returns JSON encoded book in http response
		if err := json.NewEncoder(w).Encode(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("PUT")

	mux.HandleFunc("/books/{pk}", func(w http.ResponseWriter, r *http.Request) {
		pk, _ := strconv.ParseInt(gmux.Vars(r)["pk"], 10, 64)
		// ensure book belongs to user
		// var b Book

		////////////

		// dbQueryRow searches the Database for a user based on primary key,
		// then scans the query results, placing each column into the user object.
		if b, err := db.Exec("SELECT * FROM books WHERE pk=? and user=?", pk, getSessionString(r, "user")); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println("416\terr", b)

			return
		}
		// fmt.Println("416\terr", err)

		////////////
		// performs external OS command to delete a DB book
		// based on query string's "pk" value
		if _, err := db.Exec("DELETE from books WHERE pk = ?", pk); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// writes 200 status on success
		w.WriteHeader(http.StatusOK)
	}).Methods("DELETE")

	// store Negroni object with default middleware for
	// Panic Recovery, Logging and Static File Serving
	n := negroni.Classic()

	// add DB verification to middleware stack
	n.Use(negroni.HandlerFunc(pingDB))
	n.Use(negroni.HandlerFunc(startSession))

	// adds user session verification
	n.Use(negroni.HandlerFunc(verifyUser))

	// PLACE ROUTER AT END OF MIDDLEWARE PIPELINE
	// Negroni use Gmux Router
	n.UseHandler(mux)

	// start paort 8080 webserver
	n.Run(":8080")
}

// searches OCLC Book API and returns a book's OCLC Work Identifier (OWI) ID
func find(id string) (ClassifyBookResponse, error) {
	var c ClassifyBookResponse

	// escapes query string id for proper HTTP request
	body, err := classifyAPI("http://classify.oclc.org/classify2/Classify?&summary=true&owi=" + url.QueryEscape(id))

	// returns empty ClassifyBookResponse struct and error
	if err != nil {
		return ClassifyBookResponse{}, err
	}

	// transforms book object into ClassifyBookResponse object
	err = xml.Unmarshal(body, &c)

	return c, err
}

// search() uses search query string to return search results list
// after sending HTTP request to OCLC Book API
func search(query string) ([]SearchResult, error) {
	var c ClassifySearchResponse

	// escapes query string query for proper HTTP request
	body, err := classifyAPI("http://classify.oclc.org/classify2/Classify?&summary=true&title=" + url.QueryEscape(query))

	if err != nil {
		return []SearchResult{}, err
	}

	// Parses XML data and stores result in ClassifySearchResponse
	// object using reflection
	err = xml.Unmarshal(body, &c)

	return c.Results, err
}

// sends HTTP request to OCLC Book API amd returns byte slice
func classifyAPI(url string) ([]byte, error) {
	var resp *http.Response
	var err error

	// queries book collection titles and returns books with a title
	// matching the search query, or returns empty []byte{} and error
	if resp, err = http.Get(url); err != nil {
		return []byte{}, err
	}

	// closes response body at end of the func
	defer resp.Body.Close()

	// reads and returns all bytes in response body (until an error or EOF)
	return ioutil.ReadAll(resp.Body)
}
