package cms

import (
	"net/http"
	"strings"
	"time"
)

// ServeIndex responses to HTTP requests
func ServeIndex(w http.ResponseWriter, r *http.Request) {
	p := &Page{
		Title:   "GO Content Management System",
		Content: "Welcome",
		Posts: []*Post{
			&Post{
				Title:         "Testing 1-2-3",
				Content:       "Bravo Good Chap. Stellar Test",
				DatePublished: time.Now(),
			},
			&Post{
				Title:         "Charlie Rose and the Chocolate Factory: Part I",
				Content:       "Very exciting read Charlie. Looking forward to Part II.",
				DatePublished: time.Now(),
				Comments: []*Comment{
					&Comment{
						Author:        "Dr. Foobazman",
						Comment:       "Excellent post!",
						DatePublished: time.Now().Add(-time.Hour / 2),
					},
				},
			},
		},
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
	path := strings.TrimLeft(r.URL.Path, "/post/")

	if path == "" {
		http.NotFound(w, r)
		return
	}

	p := &Post{
		Title:   strings.ToTitle(path),
		Content: "Welcome to the Post",
		Comments: []*Comment{
			&Comment{
				Author:        "Dr. Post-Author",
				Comment:       "Cool post.",
				DatePublished: time.Now().Add(-time.Hour / 2),
			},
		},
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
			DatePublished: time.Now(),
			Comments: []*Comment{
				&Comment{
					Author:        "Undefined User",
					Comment:       r.FormValue("comment"),
					DatePublished: time.Now().Add(-time.Hour / 2),
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
