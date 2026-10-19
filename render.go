package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// page is everything the HTML template needs to show one markdown file.
type page struct {
	Title       string
	FileName    string
	Dir         string
	BaseHref    template.URL
	FrontMatter template.HTML
	Body        template.HTML
	TOC         []tocEntry
	ReadMinutes int
	HasMermaid  bool
}

type tocEntry struct {
	Level int
	ID    string
	Text  string
}

var codeFormatter = chromahtml.New(chromahtml.WithClasses(true))

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM, // tables, task lists, strikethrough, autolinks
		extension.Footnote,
		extension.DefinitionList,
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
		parser.WithASTTransformers(util.Prioritized(alertTransformer{}, 100)),
	),
	goldmark.WithRendererOptions(
		// Local files written by you or your tools: allow <details>, <br>, <img> etc.
		html.WithUnsafe(),
		renderer.WithNodeRenderers(util.Prioritized(codeRenderer{}, 100)),
	),
)

func buildPage(src []byte, path string) (*page, error) {
	fm, body := splitFrontMatter(src)

	doc := md.Parser().Parse(text.NewReader(body))

	p := &page{
		FileName:    filepath.Base(path),
		Dir:         shortDir(filepath.Dir(path)),
		BaseHref:    template.URL(dirURL(filepath.Dir(path))),
		ReadMinutes: max(1, len(strings.Fields(string(body)))/220),
	}

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n := n.(type) {
		case *ast.Heading:
			title := plainText(n, body)
			if n.Level == 1 && p.Title == "" {
				p.Title = title
			}
			if id, ok := n.AttributeString("id"); ok && n.Level <= 3 {
				p.TOC = append(p.TOC, tocEntry{Level: n.Level, ID: string(id.([]byte)), Text: title})
			}
			return ast.WalkSkipChildren, nil
		case *ast.FencedCodeBlock:
			if isMermaid(n.Language(body)) {
				p.HasMermaid = true
			}
		}
		return ast.WalkContinue, nil
	})

	var buf bytes.Buffer
	if err := md.Renderer().Render(&buf, body, doc); err != nil {
		return nil, err
	}
	// Wrap tables so wide ones scroll on their own instead of stretching the page.
	out := strings.ReplaceAll(buf.String(), "<table>", `<div class="table-wrap"><table>`)
	out = strings.ReplaceAll(out, "</table>", "</table></div>")
	p.Body = template.HTML(out)

	if fm != nil {
		fields := parseFrontMatter(fm)
		p.FrontMatter = renderFrontMatter(fields)
		if p.Title == "" {
			p.Title = fields.get("title")
		}
	}
	if p.Title == "" {
		p.Title = p.FileName
	}
	if len(p.TOC) < 3 {
		p.TOC = nil
	}
	return p, nil
}

// shortDir abbreviates the home directory to "~" for display.
func shortDir(dir string) string {
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		if rel, err := filepath.Rel(home, dir); err == nil && !strings.HasPrefix(rel, "..") {
			return filepath.Join("~", rel)
		}
	}
	return dir
}

// dirURL returns a file:// URL for dir, used as <base> so relative images and links resolve.
func dirURL(dir string) string {
	s := filepath.ToSlash(dir)
	if !strings.HasPrefix(s, "/") {
		s = "/" + s // Windows drive letters
	}
	if !strings.HasSuffix(s, "/") {
		s += "/"
	}
	return (&url.URL{Scheme: "file", Path: s}).String()
}

func plainText(n ast.Node, src []byte) string {
	var b strings.Builder
	ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch c := c.(type) {
		case *ast.Text:
			b.Write(c.Segment.Value(src))
			if c.SoftLineBreak() {
				b.WriteByte(' ')
			}
		case *ast.String:
			if !c.IsCode() {
				b.Write(c.Value)
			}
		case *ast.RawHTML:
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimSpace(b.String())
}

// --- Code blocks: syntax highlighting and mermaid diagrams ---

type codeRenderer struct{}

func (r codeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.render)
}

func isMermaid(lang []byte) bool {
	return strings.EqualFold(string(lang), "mermaid")
}

func (codeRenderer) render(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.FencedCodeBlock)
	lang := string(n.Language(src))

	var code bytes.Buffer
	for i := 0; i < n.Lines().Len(); i++ {
		seg := n.Lines().At(i)
		code.Write(seg.Value(src))
	}

	if isMermaid(n.Language(src)) {
		// The source stays visible if mermaid can't load (offline) or fails to parse.
		w.WriteString(`<div class="diagram"><pre class="mermaid">`)
		template.HTMLEscape(w, code.Bytes())
		w.WriteString("</pre></div>\n")
		return ast.WalkSkipChildren, nil
	}

	w.WriteString(`<div class="code-block">`)
	if lang != "" {
		fmt.Fprintf(w, `<span class="code-lang">%s</span>`, template.HTMLEscapeString(lang))
	}
	w.WriteString(`<button class="code-copy" type="button">Copy</button>`)

	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	tokens, err := chroma.Coalesce(lexer).Tokenise(nil, code.String())
	if err == nil {
		err = codeFormatter.Format(w, styles.Get("github"), tokens)
	}
	if err != nil {
		w.WriteString(`<pre><code>`)
		template.HTMLEscape(w, code.Bytes())
		w.WriteString(`</code></pre>`)
	}
	w.WriteString("</div>\n")
	return ast.WalkSkipChildren, nil
}

// codeCSS returns the highlighting theme, following the system light/dark setting.
func codeCSS() string {
	var light, dark bytes.Buffer
	codeFormatter.WriteCSS(&light, styles.Get("github"))
	codeFormatter.WriteCSS(&dark, styles.Get("github-dark"))
	// Scope each theme to its own media query so token colors from one never leak into the other.
	return "@media not (prefers-color-scheme: dark) {\n" + light.String() + "}\n" +
		"@media (prefers-color-scheme: dark) {\n" + dark.String() + "}\n"
}

// --- GitHub alerts: > [!NOTE], > [!TIP], > [!IMPORTANT], > [!WARNING], > [!CAUTION] ---

var alertRe = regexp.MustCompile(`(?i)^\[!(note|tip|important|warning|caution)\]$`)

var alertIcons = map[string]string{
	"note":      `<circle cx="8" cy="8" r="6.25"/><path d="M8 7.25v3.75M8 5h.01"/>`,
	"tip":       `<path d="M5.75 11c-1-.9-2-2-2-3.75a4.25 4.25 0 0 1 8.5 0c0 1.75-1 2.85-2 3.75M6.25 13.25h3.5M7 15h2"/>`,
	"important": `<path d="M2.25 2.75h11.5v8.5h-6l-3 2.5v-2.5h-2.5z"/><path d="M8 5v3M8 9.75h.01"/>`,
	"warning":   `<path d="M8 1.75 14.75 13.75H1.25Z"/><path d="M8 6v3.5M8 11.5h.01"/>`,
	"caution":   `<path d="M5.25 1.75h5.5l3.5 3.5v5.5l-3.5 3.5h-5.5l-3.5-3.5v-5.5Z"/><path d="M8 4.75v3.75M8 10.75h.01"/>`,
}

type alertTransformer struct{}

func (alertTransformer) Transform(doc *ast.Document, reader text.Reader, _ parser.Context) {
	var quotes []*ast.Blockquote
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if bq, ok := n.(*ast.Blockquote); ok && entering {
			quotes = append(quotes, bq)
		}
		return ast.WalkContinue, nil
	})
	for _, bq := range quotes {
		makeAlert(bq, reader.Source())
	}
}

func makeAlert(bq *ast.Blockquote, src []byte) {
	p, ok := bq.FirstChild().(*ast.Paragraph)
	if !ok || p.Lines().Len() == 0 {
		return
	}
	first := p.Lines().At(0)
	m := alertRe.FindSubmatch(bytes.TrimSpace(first.Value(src)))
	if m == nil {
		return
	}
	kind := strings.ToLower(string(m[1]))

	// Drop the inline nodes that make up the "[!KIND]" marker line.
	for c := p.FirstChild(); c != nil; {
		next := c.NextSibling()
		t, ok := c.(*ast.Text)
		if !ok || t.Segment.Start >= first.Stop {
			break
		}
		p.RemoveChild(p, c)
		c = next
	}
	if p.ChildCount() == 0 {
		bq.RemoveChild(bq, p)
	}

	icon := ast.NewString([]byte(`<svg class="alert-icon" viewBox="0 0 16 16" aria-hidden="true">` + alertIcons[kind] + `</svg>`))
	icon.SetCode(true) // written verbatim
	title := ast.NewParagraph()
	title.SetAttributeString("class", []byte("markdown-alert-title"))
	title.AppendChild(title, icon)
	title.AppendChild(title, ast.NewString([]byte(strings.ToUpper(kind[:1])+kind[1:])))
	bq.InsertBefore(bq, bq.FirstChild(), title)
	bq.SetAttributeString("class", []byte("markdown-alert markdown-alert-"+kind))
}

// --- YAML front matter (common in AI-generated plans, skills and agent files) ---

var frontMatterKey = regexp.MustCompile(`^[A-Za-z0-9_-]+:`)

// splitFrontMatter separates a leading "---" YAML block from the markdown body.
func splitFrontMatter(src []byte) (fm, body []byte) {
	lines := bytes.SplitAfter(src, []byte("\n"))
	if len(lines) < 3 || string(bytes.TrimRight(lines[0], "\r\n")) != "---" {
		return nil, src
	}
	// Require a "key:" first line so a document that opens with a horizontal rule is left alone.
	if !frontMatterKey.Match(lines[1]) {
		return nil, src
	}
	offset := len(lines[0])
	for _, l := range lines[1:] {
		if t := string(bytes.TrimRight(l, "\r\n")); t == "---" || t == "..." {
			return src[len(lines[0]):offset], src[offset+len(l):]
		}
		offset += len(l)
	}
	return nil, src
}

type fmField struct{ Key, Value string }

type fmFields []fmField

func (f fmFields) get(key string) string {
	for _, field := range f {
		if field.Key == key {
			return strings.Trim(field.Value, `"'`)
		}
	}
	return ""
}

// parseFrontMatter reads top-level keys; nested or multi-line values are kept as text.
func parseFrontMatter(fm []byte) fmFields {
	var fields fmFields
	for _, line := range strings.Split(string(fm), "\n") {
		line = strings.TrimRight(line, "\r")
		if frontMatterKey.MatchString(line) {
			k, v, _ := strings.Cut(line, ":")
			fields = append(fields, fmField{k, strings.TrimSpace(v)})
		} else if len(fields) > 0 && strings.TrimSpace(line) != "" {
			last := &fields[len(fields)-1]
			last.Value = strings.TrimSpace(last.Value + "\n" + strings.TrimSpace(line))
		}
	}
	return fields
}

func renderFrontMatter(fields fmFields) template.HTML {
	if len(fields) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<details class="front-matter" open><summary>Front matter</summary><div class="table-wrap"><table><tbody>`)
	for _, f := range fields {
		fmt.Fprintf(&b, "<tr><th>%s</th><td>%s</td></tr>",
			template.HTMLEscapeString(f.Key), template.HTMLEscapeString(f.Value))
	}
	b.WriteString(`</tbody></table></div></details>`)
	return template.HTML(b.String())
}
