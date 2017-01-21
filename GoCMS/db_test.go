package cms

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

var p *Page
var po *Post

func Test_CreatePage(t *testing.T) {
	p = &Page{
		Title:   "test title",
		Content: "test content",
	}

	id, err := CreatePage(p)
	if err != nil {
		t.Errorf("Failed to create page: %s\n", err.Error())
	}
	p.ID = id
}

func Test_GetPage(t *testing.T) {
	page, err := GetPage(strconv.Itoa(p.ID))
	if err != nil {
		t.Errorf("Failed to get page: %s\n", err.Error())
	}
	if page.ID != p.ID {
		t.Errorf("Page IDs do not match: expected %d\n but received %d\n", p.ID, page.ID)
	}
	if reflect.DeepEqual(page, p) != true {
		t.Errorf("Pages do not match: expected %+v\n but received %+v\n", p, page)
	}
}

func Test_CreatePost(t *testing.T) {
	po = &Post{
		Title:         "post title",
		Content:       "post content",
		DatePublished: time.Now().UTC(),
		Comments: []*Comment{
			&Comment{
				Author:        "test comment author",
				Comment:       "test comment comment",
				DatePublished: time.Now().UTC().Add(-time.Hour / 2),
			},
		},
	}

	id, err := CreatePost(po)
	if err != nil {
		t.Errorf("Failed to create post: %s\n", err.Error())
	}

	po.ID = id

}

func Test_ReadPost(t *testing.T) {
	post, err := GetPost(strconv.Itoa(po.ID))
	if err != nil {
		t.Errorf("Failed to get post: %s\n", err.Error())
	}
	if post.ID != po.ID {
		t.Errorf("Post IDs do not match: expected %d\n but received %d\n", po.ID, post.ID)
	}
	if post.Title != po.Title {
		t.Errorf("Post Titles do not match: expected %s\n but received %s\n", po.Title, post.Title)
	}
	if post.Content != po.Content {
		t.Errorf("Post Contents do not match: expected %s\n but received %s\n", po.Content, post.Content)
	}

	// TEST FAILING NEED BUG FIX: Postgres rounds time to Milliseconds and Go truncates time at Nanoseconds.
	// see https://github.com/lib/pq/issues/227
	//if post.DatePublished != po.DatePublished {
	//t.Errorf("Post DatePublished do not match: expected %s\n but received %s\n", po.DatePublished, post.DatePublished)
	//}

	if post.Comments[0].ID != po.Comments[0].ID {
		t.Errorf("Post Comments.ID do not match: expected %d\n but received %d\n", po.Comments[0].ID, post.Comments[0].ID)
	}

	if post.Comments[0].Author != po.Comments[0].Author {
		t.Errorf("Post Comments.Author do not match: expected %s\n but received %s\n", po.Comments[0].Author, post.Comments[0].Author)
	}

	if post.Comments[0].Comment != po.Comments[0].Comment {
		t.Errorf("Post Comments.Author do not match: expected %s\n but received %s\n", po.Comments[0].Comment, post.Comments[0].Comment)
	}

}
