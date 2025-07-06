package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
)

func copyFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("opening source file %s: %w", srcPath, err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("creating destination file %s: %w", dstPath, err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copying file: %w", err)
	}
	return nil
}

func convertMarkdown(srcDir, dstDir string) error {
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("creating output directory %s: %w", dstDir, err)
	}
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("reading source directory %s: %w", srcDir, err)
	}

	md := goldmark.New()
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
			return fmt.Errorf("reading markdown file %s: %w", srcPath, err)
		}

		var buf bytes.Buffer
		if err := md.Convert(input, &buf); err != nil {
			return fmt.Errorf("converting markdown %s: %w", srcPath, err)
		}

		outName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())) + ".html"
		dstPath := filepath.Join(dstDir, outName)
		if err := os.WriteFile(dstPath, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("writing html file %s: %w", dstPath, err)
		}
	}
	return nil
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
	srcPath := "index.html"
	dstPath := filepath.Join(*outPath, "index.html")

	if err := copyFile(srcPath, dstPath); err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		os.Exit(1)
	}
	if err := convertMarkdown("posts", filepath.Join(*outPath, "posts")); err != nil {
		fmt.Printf("Error converting markdown: %v\n", err)
	}
	fmt.Println("Done.")
}
