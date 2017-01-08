package cms

import "html/template"

// Tmpl returns an initialized template parsed from the templates in the templates directory. Panics on non-nil errors
var Tmpl = template.Must(template.ParseGlob("../templates/*"))

// A Page represents HTML page content
type Page struct {
	Title   string
	Content string
}
