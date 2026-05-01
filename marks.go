package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pkg/browser"
	"github.com/yuin/goldmark"
)

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

	html := fmt.Sprintf(`
	<html>
	<body style="max-width:900px;margin:3rem auto;font-family:sans-serif">
	%s
	</body>
	</html>
	`, buf.String())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, html)
	})

	go http.ListenAndServe(":8080", nil)

	browser.OpenURL("http://localhost:8080")

	select {}
}
