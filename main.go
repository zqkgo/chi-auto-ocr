package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

var imgDir string

func init() {
	flag.StringVar(&imgDir, "image_dir", "./data", "the directory containing images")
}

func main() {
	flag.Parse()

	err := os.Chdir(imgDir)
	if err != nil {
		log.Panic("failed to change dir, err: ", err)
	}
	handled := make(map[string]bool)
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		filepath.Walk("./", func(path string, info fs.FileInfo, err error) error {
			if handled[path] {
				return nil
			}
			if info.IsDir() && info.Name() != "." {
				log.Println("meet dir: ", info.Name())
				return nil
			}
			ext := filepath.Ext(path)
			if ext == ".png" || ext == ".jpeg" || ext == ".jpg" {
				allInOneMarkdownOCR(path)
				handled[path] = true
			}
			return nil
		})
	}
}
