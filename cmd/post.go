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
	Title     string    `yaml:"title"`
	Date      time.Time `yaml:"date"`
	Teaser    string    `yaml:"teaser"`
	Thumbnail string    `yaml:"thumbnail"`
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

const pageTemplate = `
<head>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <!-- Pick the theme you like: github-markdown.css auto-switches light/dark -->
  <link
    rel="stylesheet"
    href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.8.1/github-markdown.min.css"
  >
</head>
<body>
  <!-- This wrapper applies GitHub’s spacing, fonts, headings, tables, etc. -->
  <article class="markdown-body">
    <!-- insert your rendered HTML here -->
    {{ .Content }}
  </article>
</body>

`

// tmpl is the parsed template ready for execution.
var tmpl = template.Must(template.New("page").Parse(pageTemplate))

func RenderPage(content string) (string, error) {
	// Use template.HTML to tell the engine this is trusted HTML.
	data := struct {
		Content template.HTML
	}{
		Content: template.HTML(content),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
