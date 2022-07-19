package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var cnt int

func main() {
	err := os.Chdir("./data")
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
				cnt++
				cmd := exec.Command("tesseract", path, fmt.Sprintf("%d", cnt), "-l", "chi_sim", "-c", "preserve_interword_spaces=1")
				err := cmd.Run()
				if err != nil {
					log.Printf("failed to run tesseract, err: %+v\n", err)
				}
				handled[path] = true
			}
			return nil
		})
	}
}
