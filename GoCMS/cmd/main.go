package main

import (
	cms "GoCMS"
	"os"
)

func main() {
	p := &cms.Page{
		Title:   "Hello World Template Title",
		Content: "FooBaz Template Content",
	}

	cms.Tmpl.ExecuteTemplate(os.Stdout, "index", p)
}
