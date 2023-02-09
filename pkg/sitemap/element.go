package sitemap

import (
	"context"
	"time"

	"github.com/ReanSn0w/gew/v2/pkg/view"
	"github.com/ReanSn0w/sitemap/pkg/xml"
)

const (
	PriorityHigh   = 1
	PriorityMedium = 0.7
	PriorityLow    = 0.4
)

func NewElement(location string) *Element {
	return &Element{
		Location:     location,
		LastModified: time.Now(),
		Changefreq:   "weekly",
		Priority:     PriorityMedium,
	}
}

type Element struct {
	Location     string
	LastModified time.Time
	Changefreq   string
	Priority     float64
}

func (el Element) Body(ctx context.Context) view.View {
	return xml.NewNode(false, "url", nil)(
		xml.NewNode(true, "loc", nil)(xml.Url(el.Location)),
		xml.NewNode(true, "lastmod", nil)(xml.String(el.LastModified.Format("2006-01-02T15:04:05-07:00"))),
		xml.NewNode(true, "changefreq", nil)(xml.String(el.Changefreq)),
		xml.NewNode(true, "priority", nil)(xml.String("%v", el.Priority)),
	)
}
