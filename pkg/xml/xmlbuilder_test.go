package xml_test

import (
	"bytes"
	"testing"

	"github.com/ReanSn0w/gew/v2/pkg/view"
	"github.com/ReanSn0w/sitemap/pkg/xml"
)

func Test_Builder(t *testing.T) {
	val := "hello, world!"
	if val != buildnode(xml.String(val)) {
		t.Log("strings not equal")
		t.Fail()
	}
}

func Test_BuildNode(t *testing.T) {
	val := "<loc>\nhello world\n</loc>\n"
	if val != buildnode(xml.NewNode(false, "loc", nil)(xml.String("hello world"))) {
		t.Log("strings not equal")
		t.Fail()
	}
}

func Test_BuildInlineNode(t *testing.T) {
	val := "<loc>hello world</loc>\n"
	if val != buildnode(xml.NewNode(true, "loc", nil)(xml.String("hello world"))) {
		t.Log("strings not equal")
		t.Fail()
	}
}

func Test_BuildNodeWithParameters(t *testing.T) {
	val := "<loc first=\"first_val\" second>\nhello world\n</loc>\n"
	if val != buildnode(xml.NewNode(false, "loc", map[string]string{
		"first":  "first_val",
		"second": "",
	})(xml.String("hello world"))) {
		t.Log("strings not equal")
		t.Fail()
	}
}

func Test_BuildInlineNodeWithParameters(t *testing.T) {
	val := "<loc first=\"first_val\" second>hello world</loc>\n"
	if val != buildnode(xml.NewNode(true, "loc", map[string]string{
		"first":  "first_val",
		"second": "",
	})(xml.String("hello world"))) {
		t.Log("strings not equal")
		t.Fail()
	}
}

func buildnode(item view.View) string {
	buffer := new(bytes.Buffer)
	xml.Build(buffer, item)
	return buffer.String()
}
