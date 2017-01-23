package cms

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ServeIndex responses to HTTP requests
func ServeIndex(w http.ResponseWriter, r *http.Request) {

	posts, err := GetPosts()
	fmt.Println(posts, err)
	p := &Page{
		Title:   "GO Content Management System",
		Content: "Welcome",
		Posts:   posts,
	}

	Tmpl.ExecuteTemplate(w, "page", p)
}

// HandleNew handles preview logic
func HandleNew(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		Tmpl.ExecuteTemplate(w, "new", nil)
	case "POST":
		createPostTemplate(w, r)
	default:
		http.Error(w, "Method not supported: "+r.Method, http.StatusMethodNotAllowed)
	}
}

// ServePage serves a page base HTTP requests to "root/page"
func ServePage(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/page/")

	if path == "" {
		pages, err := GetPages()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		Tmpl.ExecuteTemplate(w, "pages", pages)
		return
	}

	page, err := GetPage(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	Tmpl.ExecuteTemplate(w, "page", page)
}

// ServePost serves a content post based on HTTP requests to "root/post"
func ServePost(w http.ResponseWriter, r *http.Request) {
	pathID := strings.TrimLeft(r.URL.Path, "/post/")

	if pathID == "" {
		http.NotFound(w, r)
		return
	}

	p, err := GetPost(pathID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	Tmpl.ExecuteTemplate(w, "post", p)
}

func createPostTemplate(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	contentType := r.FormValue("content-type")
	r.ParseForm()

	if contentType == "page" {
		p := &Page{
			Title:   title,
			Content: content,
		}

		_, err := CreatePage(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		Tmpl.ExecuteTemplate(w, "page", p)
		return
	}

	if contentType == "post" {
		p := &Post{
			Title:         title,
			Content:       content,
			DatePublished: time.Now().UTC(),
			Comments: []*Comment{
				&Comment{
					Author:        "Undefined User",
					Comment:       r.FormValue("comment"),
					DatePublished: time.Now().UTC().Add(-time.Hour / 2),
				},
			},
		}

		_, err := CreatePost(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		Tmpl.ExecuteTemplate(w, "post", p)

		return
	}

}
