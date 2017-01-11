package cms

import (
	"database/sql"
	"fmt"

	//Postgres Driver
	_ "github.com/lib/pq"
)

var store = newDB()

// PgStore represents database connection
type PgStore struct {
	DB *sql.DB
}

func newDB() *PgStore {
	db, err := sql.Open("postgres", "user=gocms dbname=gocms sslmode=disable")
	if err != nil {
		panic(err)
	}

	return &PgStore{
		DB: db,
	}
}

// CreatePage creates and inserts a Page into the database table
func CreatePage(p *Page) (int, error) {
	var id int
	err := store.DB.QueryRow("INSERT INTO pages(title, content) VALUES($1, $2) RETURNING id", p.Title, p.Content).Scan(&id)
	fmt.Println(id)
	return id, err
}

// GetPage returns a Page queried from a database table record
func GetPage(id string) (*Page, error) {
	var p Page
	err := store.DB.QueryRow("SELECT * FROM pages WHERE id = $1", id).Scan(&p.ID, &p.Title, &p.Content)
	return &p, err
}

// GetPages returns all Page objects from a database table query
func GetPages() ([]*Page, error) {
	rows, err := store.DB.Query("SELECT * FROM pages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pages := []*Page{}
	for rows.Next() {
		var p Page
		err = rows.Scan(&p.ID, &p.Title, &p.Content)

		if err != nil {
			return nil, err
		}

		pages = append(pages, &p)
	}

	return pages, nil
}
