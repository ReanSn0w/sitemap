package main

import (
	"fmt"
	"time"

	"github.com/ReanSn0w/sitemap/pkg/service"
	"github.com/ReanSn0w/sitemap/pkg/sitemap"
)

func main() {
	service.New(
		"http://localhost:8080", true,
		exampleSMBuilder,
	).Run()
}

func exampleSMBuilder() ([]sitemap.Element, error) {
	res := []sitemap.Element{}

	for i := 0; i < 10000; i++ {
		res = append(res, *sitemap.NewElement(fmt.Sprintf("https://example.com/blog/%v", i)))
	}

	res = append(res, sitemap.Element{
		Location:     "https://example.com",
		LastModified: time.Now(),
		Changefreq:   "daily",
		Priority:     1.0,
	},
		sitemap.Element{
			Location:     "https://example.com/about",
			LastModified: time.Now(),
			Changefreq:   "daily",
			Priority:     0.8,
		},
		sitemap.Element{
			Location:     "https://example.com/contact",
			LastModified: time.Now(),
			Changefreq:   "daily",
			Priority:     0.8,
		},
		sitemap.Element{
			Location:     "https://example.com/blog",
			LastModified: time.Now(),
			Changefreq:   "daily",
			Priority:     0.8,
		},
	)

	return res, nil
}
