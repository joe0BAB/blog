// main.go
package main

import (
	"embed"
	"html/template"
	"os"
	"time"
)

// Post represents a blog post in the TOC.
type Post struct {
	Title      string
	Date       time.Time
	Teaser     string
	Thumbnail  string
	Background string
	URL        string
	Icon       string
}

//go:embed templates
var tmplFS embed.FS

func RenderPosts(posts []Post, dstFile string) error {
	tmpl, err := template.ParseFS(tmplFS, "templates/index.html")
	if err != nil {
		return err
	}

	f, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer f.Close()

	data := struct{ Posts []Post }{Posts: posts}
	if err := tmpl.Execute(f, data); err != nil {
		return err
	}

	return nil
}
