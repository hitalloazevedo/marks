package main

import (
	_ "embed"
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/pkg/browser"
	"github.com/yuin/goldmark"
)

//go:embed static/github-markdown.css
var githubCSS string

func main() {
	file := os.Args[1]

	md, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	if err := goldmark.Convert(md, &buf); err != nil {
		log.Fatal(err)
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
	<style>
		html {
			margin: 4rem;
		}
	</style>
  <style>%s</style>
<body class="markdown-body">
%s
</body>
</html>
`, githubCSS, buf.String())

	tmp, err := os.CreateTemp("", "marks-*.html")
	if err != nil {
		log.Fatal(err)
	}

	tmp.WriteString(html)
	tmp.Close()

	browser.OpenFile(tmp.Name())
}
