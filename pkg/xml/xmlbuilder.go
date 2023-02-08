package xml

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/ReanSn0w/gew/v2/pkg/view"
)

func Build(wr io.Writer, doc view.View) {
	view.Build(doc, context.Background(), func(i interface{}, ctx context.Context) {
		switch v := i.(type) {
		case *XMLNode:
			wr.Write([]byte(fmt.Sprintf("<%s%s>", v.Name, v.parameters())))

			// if !v.Inline {
			// 	wr.Write([]byte("\n"))
			// }

			Build(wr, v.Value)

			// if !v.Inline {
			// 	wr.Write([]byte("\n"))
			// }

			wr.Write([]byte(fmt.Sprintf("</%s>", v.Name)))
		case []byte:
			wr.Write(v)
		}
	})
}

func NewNode(inline bool, name string, params Parameters) func(content ...view.View) view.View {
	return func(content ...view.View) view.View {
		return view.External(&XMLNode{
			Name:   name,
			Value:  view.Group(content...),
			Params: params,
		})
	}
}

type Parameters map[string]string

type XMLNode struct {
	Inline bool
	Name   string
	Value  view.View
	Params Parameters
}

func (node *XMLNode) parameters() string {
	if len(node.Params) == 0 {
		return ""
	}

	res := ""
	for key, value := range node.Params {
		if value == "" {
			res += " " + key
		} else {
			res += " " + key + "=\"" + value + "\""
		}
	}

	return res
}

func String(val string, inserts ...interface{}) view.View {
	val = fmt.Sprintf(val, inserts...)
	return view.External([]byte(val))
}

func Url(val string) view.View {
	val = strings.Replace(val, "<", "&lt;", -1)
	val = strings.Replace(val, ">", "&gt;", -1)
	val = strings.Replace(val, "&", "&amp;", -1)
	val = strings.Replace(val, "\"", "&quot;", -1)
	val = strings.Replace(val, "'", "&apos;", -1)

	return view.External([]byte(val))
}
