package sitemap

import (
	"net/http"
	"strconv"

	"github.com/ReanSn0w/gew/v2/pkg/view"
	"github.com/ReanSn0w/sitemap/pkg/xml"
)

func (sm *sitemap) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != "/sitemap.xml" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		page := sm.pageFromRequest(r)
		var value view.View
		if page.Number == 0 {
			value = sm.Index(uint(page.Count))
		} else {
			value = sm.Page(uint(page.Count), uint(page.Number))
		}

		xml.Build(w, value)
	})
}

type Page struct {
	Number int
	Count  int
}

func (sm *sitemap) pageFromRequest(r *http.Request) Page {
	value, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil || value <= 0 {
		// if value is not set or not a number
		// return default value
		value = 1000
	}

	number, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if number <= 0 {
		number = 0
	}

	return Page{
		Number: number,
		Count:  value,
	}
}
