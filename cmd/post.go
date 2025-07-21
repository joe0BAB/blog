package main

import (
	"bytes"
	"errors"
	"gopkg.in/yaml.v3"
	"html/template"
	"strings"
	"time"
)

type Meta struct {
	Title      string    `yaml:"title"`
	Date       time.Time `yaml:"date"`
	Teaser     string    `yaml:"teaser"`
	Thumbnail  string    `yaml:"thumbnail"`
	Background string    `yaml:"background"`
	Icon       string    `yaml:"icon"`
}

func ParseFrontMatter(input string) (Meta, string, error) {
	var meta Meta

	const delim = "---"
	if !strings.HasPrefix(input, delim) {
		return meta, input, nil
	}

	parts := strings.SplitN(input, delim, 3)
	if len(parts) < 3 {
		return meta, "", errors.New("unclosed front-matter block")
	}

	rawYAML := parts[1]
	body := parts[2]

	if len(body) > 0 && body[0] == '\n' {
		body = body[1:]
	}

	if err := yaml.Unmarshal([]byte(rawYAML), &meta); err != nil {
		return meta, "", err
	}

	return meta, body, nil
}

func RenderPage(content string, post Post) (string, error) {
	// 1) Create a new template and register your funcs.
	tmpl := template.New("post.html").Funcs(template.FuncMap{
		// formatDate takes a time.Time and a layout string.
		"formatDate": func(t time.Time, layout string) string {
			return t.Format(layout)
		},
	})

	// 2) Parse your files into that template instance.
	tmpl, err := tmpl.ParseFS(tmplFS, "templates/post.html")
	if err != nil {
		return "", err
	}

	// 3) Build your data context as before.
	data := struct {
		Content template.HTML
		Post
	}{
		Content: template.HTML(content),
		Post:    post,
	}

	// 4) Execute.
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
