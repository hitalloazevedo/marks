package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/browser"
)

//go:embed static/github-markdown.css
var githubCSS string

//go:embed static/marks.css
var marksCSS string

//go:embed static/marks.js
var marksJS string

//go:embed static/page.html
var pageHTML string

var pageTemplate = template.Must(template.New("page").Parse(pageHTML))

func main() {
	if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Fprintln(os.Stderr, "usage: marks <file.md>")
		os.Exit(2)
	}

	path, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	md, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	html, err := renderHTML(md, path)
	if err != nil {
		log.Fatal(err)
	}

	tmp, err := os.CreateTemp("", "marks-*.html")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := tmp.Write(html); err != nil {
		log.Fatal(err)
	}
	if err := tmp.Close(); err != nil {
		log.Fatal(err)
	}

	if err := browser.OpenFile(tmp.Name()); err != nil {
		log.Fatal(err)
	}
}

func renderHTML(md []byte, path string) ([]byte, error) {
	p, err := buildPage(md, path)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = pageTemplate.Execute(&buf, struct {
		*page
		CSS template.CSS
		JS  template.JS
	}{p, template.CSS(githubCSS + "\n" + codeCSS() + "\n" + marksCSS), template.JS(marksJS)})
	return buf.Bytes(), err
}
