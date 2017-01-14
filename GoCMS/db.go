package cms

import (
	"database/sql"

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

	return id, err
}

// CreatePost creates and inserts a Post into the database table
func CreatePost(p *Post) (int, error) {
	var id int
	err := store.DB.QueryRow("INSERT INTO posts(title, content, date_created) VALUES($1, $2, $3) RETURNING id", p.Title, p.Content, p.DatePublished).Scan(&id)
	store.DB.Exec("INSERT INTO comments(author, content, date_created) VALUES($1, $2, $3)", p.Comments[0].Author, p.Comments[0].Comment, p.Comments[0].DatePublished)

	return id, err
}

// GetPost returns a Post record queried from a database table
func GetPost(id string) (*Post, error) {
	post := Post{
		Comments: []*Comment{
			&Comment{},
		},
	}

	err := store.DB.QueryRow("SELECT posts.id, posts.title, posts.content, posts.date_created, comments.author, comments.content, comments.date_created FROM posts, comments WHERE posts.id = $1 AND comments.id = $1", id).Scan(&post.ID, &post.Title, &post.Content, &post.DatePublished, &post.Comments[0].Author, &post.Comments[0].Comment, &post.Comments[0].DatePublished)

	return &post, err
}

// GetPage returns a Page record queried from a database table
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
