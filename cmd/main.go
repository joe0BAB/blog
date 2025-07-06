package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	fmt.Println("Done.")
}
