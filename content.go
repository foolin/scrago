package scrago

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)


func UGCContent(content string) string {
	p := bluemonday.UGCPolicy()
	content = p.Sanitize(content)
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	miniContent, err := m.String("text/html", content)
	if err != nil {
		return content
	}
	return miniContent
}
