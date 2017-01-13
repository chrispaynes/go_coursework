package cms

import (
	"reflect"
	"strconv"
	"testing"
)

var p *Page

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
