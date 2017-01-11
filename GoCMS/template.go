package cms

import (
	"html/template"
	"time"
)

// Tmpl returns an initialized template parsed from the templates in the templates directory. Panics on non-nil errors
var Tmpl = template.Must(template.ParseGlob("../templates/*"))

// A Page represents HTML page content
type Page struct {
	ID      int
	Title   string
	Content string
	Posts   []*Post
}

// A Post represents a blog post that belongs to a Page
type Post struct {
	ID            int
	Title         string
	Content       string
	DatePublished time.Time
	Comments      []*Comment
}

// A Comment represents a response that belongs to a Post
type Comment struct {
	ID            int
	Author        string
	Comment       string
	DatePublished time.Time
}
