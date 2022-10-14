package main

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
				cmd := exec.Command("tesseract", path, fmt.Sprintf("%d", cnt), "-l", "chi_sim+eng", "-c", "preserve_interword_spaces=1")
				err := cmd.Run()
				if err != nil {
					log.Printf("failed to run tesseract, err: %+v\n", err)
				}

				handled[path] = true

				// clear \n
				outFile := fmt.Sprintf("%d.txt", cnt)
				f, err := os.OpenFile(outFile, os.O_RDWR, 0)
				if err != nil {
					log.Printf("failed to open file %s, err: %v", outFile, err)
					return nil
				}
				defer f.Close()
				bs, err := ioutil.ReadAll(f)
				if err != nil {
					log.Printf("failed to read file %s, err: %v", outFile, err)
					return nil
				}
				newContent := strings.Replace(string(bs), "\n", "", -1)
				err = os.Truncate(outFile, 0)
				if err != nil {
					log.Printf("failed to truncate file %s, err: %v", outFile, err)
					return nil
				}
				_, err = f.Seek(0, io.SeekStart)
				if err != nil {
					log.Printf("failed to seek file %s, err: %v", outFile, err)
				}
				n, err := f.WriteString(newContent)
				if err != nil {
					log.Printf("failed to write back to file %s, err: %v", outFile, err)
					return nil
				}
				if n != len(newContent) {
					log.Printf("not fully write, want: %d, write: %d", len(newContent), n)
					return nil
				}
			}
			return nil
		})
	}
}
