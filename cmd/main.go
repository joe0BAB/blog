package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func convertMarkdown(srcDir, dstDir string) ([]Post, error) {
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return nil, fmt.Errorf("creating output directory %s: %w", dstDir, err)
	}
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return nil, fmt.Errorf("reading source directory %s: %w", srcDir, err)
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			// Optional: syntax-highlight code blocks using Chroma with the “github” style
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
			),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	var posts []Post
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		srcPath := filepath.Join(srcDir, entry.Name())
		input, err := os.ReadFile(srcPath)
		if err != nil {
			return nil, fmt.Errorf("reading markdown file %s: %w", srcPath, err)
		}
		meta, body, err := ParseFrontMatter(string(input))
		if err != nil {
			return nil, fmt.Errorf("parsing metadata: %w", err)
		}
		fmt.Printf("Title: %s (%s)\n", meta.Title, meta.Date)

		var buf bytes.Buffer
		if err := md.Convert([]byte(body), &buf); err != nil {
			return nil, fmt.Errorf("converting markdown %s: %w", srcPath, err)
		}

		outName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())) + ".html"
		dstPath := filepath.Join(dstDir, outName)

		post := Post{
			Title:      meta.Title,
			Date:       meta.Date,
			Teaser:     meta.Teaser,
			Thumbnail:  meta.Thumbnail,
			Background: meta.Background,
			URL:        filepath.Join("posts", outName),
			Icon:       meta.Icon,
		}

		page, err := RenderPage(buf.String(), post)
		if err != nil {
			return nil, fmt.Errorf("rendering template: %w", err)
		}

		if err := os.WriteFile(dstPath, []byte(page), 0644); err != nil {
			return nil, fmt.Errorf("writing html file %s: %w", dstPath, err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	outPath := flags.String("out", "", "Destination path for index.html")
	flags.Parse(os.Args[1:])
	if *outPath == "" {
		fmt.Printf("Error: -out flag is required\n")
		os.Exit(1)
	}
	os.RemoveAll(*outPath)
	if err := os.MkdirAll(*outPath, 0755); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", *outPath, err)
	}
	if err := copyFolder("assets", filepath.Join(*outPath, "assets")); err != nil {
		fmt.Printf("Error copying assets: %v\n", err)
	}
	posts, err := convertMarkdown("posts", filepath.Join(*outPath, "posts"))
	if err != nil {
		fmt.Printf("Error converting markdown: %v\n", err)
	}
	if err := RenderPosts(posts, filepath.Join(*outPath, "index.html")); err != nil {
		fmt.Printf("Error rendering table of content: %v\n", err)
	}
	fmt.Println("Done.")
}
