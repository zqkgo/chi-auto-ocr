package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// allInOneMarkdownOCR put all text content in one large markdown formatted file.
func allInOneMkdOCR(imgPath string) {
	if !extSupported(imgPath) {
		log.Printf("extension not supported")
		return
	}

	// convert the image
	outFile, err := ocr(imgPath)
	defer func() {
		err := os.Remove(outFile)
		if err != nil {
			log.Printf("failed to remove file '%s', err: %v", outFile, err)
		}
	}()
	if err != nil {
		log.Printf("failed to exec ocr, err: %v", err)
		return
	}

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
	text = splitLinesMdFmt(text)

	// open the large markdown file
	f, err := os.OpenFile("ocr.md", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Errorf("failed to open/create md file, err: %v", err))
	}
	defer f.Close()

	// append the image and text
	// imgLine := fmt.Sprintf("\n![](%s)\n", url.PathEscape(imgPath))
	// _, err = f.WriteString(imgLine)
	// if err != nil {
	// 	log.Printf("failed to write image line, err: %v", err)
	// }
	txtLine := fmt.Sprintf("\n%s\n", text)
	_, err = f.WriteString(txtLine)
	if err != nil {
		log.Printf("failed to write text line, err: %v", err)
	}
}

func extSupported(filePath string) bool {
	ext := filepath.Ext(filePath)
	if ext == ".png" || ext == ".jpeg" || ext == ".jpg" {
		return true
	}
	return false
}

func ocr(imgPath string) (outFile string, err error) {
	out := fmt.Sprintf("%d", rand.Intn(math.MaxInt64))
	outFile = fmt.Sprintf("%s.txt", out)
	cmd := exec.Command("tesseract", imgPath, out, "-l", "chi_sim+eng", "-c", "preserve_interword_spaces=1")
	err = cmd.Run()
	if err != nil {
		log.Printf("failed to run tesseract, err: %+v\n", err)
		return
	}
	return
}

// md 格式的行。
func splitLinesMdFmt(txt string) string {
	sep := "。"
	if !strings.Contains(txt, sep) {
		sep = "."
	}
	lines := strings.Split(txt, sep)
	txt = ""
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		txt += "- " + l + sep
		txt += "\n"
	}
	return strings.TrimRight(txt, "\n")
}

// 文本格式的行。
func splitLines(txt string) string {
	sep := "。"
	if !strings.Contains(txt, sep) {
		sep = "."
	}
	lines := strings.Split(txt, sep)
	txt = ""
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		txt += l + sep
		txt += "\n\n"
	}
	return strings.TrimRight(txt, "\n")
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
	newContent = strings.TrimSpace(newContent)
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
