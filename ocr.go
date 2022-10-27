package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var cnt int

func extSupported(filePath string) bool {
	ext := filepath.Ext(filePath)
	if ext == ".png" || ext == ".jpeg" || ext == ".jpg" {
		return true
	}
	return false
}

// allInOneMarkdownOCR put all text content in one large markdown formatted file.
func allInOneMarkdownOCR(filePath string) {
	if !extSupported(filePath) {
		log.Printf("extension not supported")
		return
	}

	// convert the image
	cnt++
	out := fmt.Sprintf("%d", cnt)
	cmd := exec.Command("tesseract", filePath, out, "-l", "chi_sim+eng", "-c", "preserve_interword_spaces=1")
	err := cmd.Run()
	if err != nil {
		log.Printf("failed to run tesseract, err: %+v\n", err)
	}
	outFile := fmt.Sprintf("%s.txt", out)
	defer func() {
		err := os.Remove(outFile)
		if err != nil {
			log.Printf("failed to remove file '%s', err: %v", outFile, err)
		}
	}()

	// clear \n
	clearNewLine(outFile)

	// fetch the text
	of, err := os.OpenFile(outFile, os.O_RDONLY, 0)
	if err != nil {
		log.Printf("failed to open '%s' to read, err: %v", outFile, err)
	}
	bs, err := ioutil.ReadAll(of)
	if err != nil {
		log.Printf("failed to read '%s', err: %v", outFile, err)
	}
	if of != nil {
		of.Close()
	}
	text := string(bs)

	// open the large markdown file
	f, err := os.OpenFile("ocr.md", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Errorf("failed to open/create md file, err: %v", err))
	}
	defer f.Close()

	// append the image and text
	imgLine := fmt.Sprintf("\n![](%s)\n", url.PathEscape(filePath))
	_, err = f.WriteString(imgLine)
	if err != nil {
		log.Printf("failed to write image line, err: %v", err)
	}
	txtLine := fmt.Sprintf("\n%s\n", text)
	_, err = f.WriteString(txtLine)
	if err != nil {
		log.Printf("failed to write text line, err: %v", err)
	}
}

// separateOCR put text content in a separate file for each image.
func separateOCR(filePath string) {
	if !extSupported(filePath) {
		log.Printf("extension not supported")
		return
	}

	cnt++
	cmd := exec.Command("tesseract", filePath, fmt.Sprintf("%d", cnt), "-l", "chi_sim+eng", "-c", "preserve_interword_spaces=1")
	err := cmd.Run()
	if err != nil {
		log.Printf("failed to run tesseract, err: %+v\n", err)
	}

	// clear \n
	outFile := fmt.Sprintf("%d.txt", cnt)
	clearNewLine(outFile)
}

func clearNewLine(filePath string) {
	f, err := os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil {
		log.Printf("failed to open file %s, err: %v", filePath, err)
	}
	defer f.Close()
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("failed to read file %s, err: %v", filePath, err)
	}
	newContent := strings.Replace(string(bs), "\n", "", -1)
	err = os.Truncate(filePath, 0)
	if err != nil {
		log.Printf("failed to truncate file %s, err: %v", filePath, err)
	}
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		log.Printf("failed to seek file %s, err: %v", filePath, err)
	}
	n, err := f.WriteString(newContent)
	if err != nil {
		log.Printf("failed to write back to file %s, err: %v", filePath, err)
	}
	if n != len(newContent) {
		log.Printf("not fully write, want: %d, write: %d", len(newContent), n)
	}
}
