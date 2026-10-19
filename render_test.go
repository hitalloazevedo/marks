package main

import (
	"strings"
	"testing"
)

func render(t *testing.T, md string) string {
	t.Helper()
	out, err := renderHTML([]byte(md), "/tmp/docs/plan.md")
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func assertContains(t *testing.T, html string, want ...string) {
	t.Helper()
	for _, w := range want {
		if !strings.Contains(html, w) {
			t.Errorf("output missing %q", w)
		}
	}
}

func assertNotContains(t *testing.T, html string, unwanted ...string) {
	t.Helper()
	for _, u := range unwanted {
		if strings.Contains(html, u) {
			t.Errorf("output should not contain %q", u)
		}
	}
}

func TestTables(t *testing.T) {
	html := render(t, "| Step | Owner |\n|:-----|------:|\n| Parse | `go` |\n")
	assertContains(t, html, `<div class="table-wrap"><table>`, "<th", "Parse", "<code>go</code>", "</table></div>")
}

func TestGFMExtras(t *testing.T) {
	html := render(t, "- [x] done\n- [ ] todo\n\n~~old~~ https://example.com\n\nText[^1].\n\n[^1]: Note.\n")
	assertContains(t, html, `type="checkbox"`, "<del>old</del>", `<a href="https://example.com">`, `class="footnotes"`)
}

func TestAlerts(t *testing.T) {
	html := render(t, "> [!WARNING]\n> Back up the database first.\n\n> plain quote\n")
	assertContains(t, html,
		`<blockquote class="markdown-alert markdown-alert-warning">`,
		`<p class="markdown-alert-title"><svg`,
		"Warning</p>",
		"Back up the database first.",
		"<blockquote>\n<p>plain quote</p>",
	)
	assertNotContains(t, html, "[!WARNING]")
}

func TestAlertWithContentInNextParagraph(t *testing.T) {
	html := render(t, "> [!tip]\n>\n> Use the cache.\n")
	assertContains(t, html, "markdown-alert-tip", "Tip</p>", "<p>Use the cache.</p>")
	assertNotContains(t, html, "<p></p>")
}

func TestCodeHighlighting(t *testing.T) {
	html := render(t, "```go\nfunc main() {}\n```\n")
	assertContains(t, html, `<div class="code-block"><span class="code-lang">go</span>`, `class="chroma"`, `<span class="kd">func</span>`)
	assertNotContains(t, html, "mermaid.esm")
}

func TestMermaid(t *testing.T) {
	html := render(t, "```mermaid\ngraph TD\n  A-->B\n```\n")
	assertContains(t, html, `<pre class="mermaid">graph TD`, "A--&gt;B", "mermaid.esm.min.mjs")
}

func TestFrontMatter(t *testing.T) {
	html := render(t, "---\ntitle: Migration plan\ntags:\n  - db\n---\n\nBody text\n")
	assertContains(t, html, `class="front-matter"`, "<th>title</th><td>Migration plan</td>", "<title>Migration plan</title>", "<p>Body text</p>")
	assertNotContains(t, html, "<hr>")
}

func TestLeadingRuleIsNotFrontMatter(t *testing.T) {
	html := render(t, "---\n\nIntro\n\n---\n")
	assertNotContains(t, html, `class="front-matter"`)
	assertContains(t, html, "<p>Intro</p>")
}

func TestTitleAndTOC(t *testing.T) {
	html := render(t, "# Big Plan\n\n## Goals\n\n## Steps\n\n### Step `one`\n")
	assertContains(t, html, "<title>Big Plan</title>", `class="toc"`, `<a href="#goals">Goals</a>`, `<a href="#step-one">Step one</a>`)
}

func TestBaseHref(t *testing.T) {
	html := render(t, "![d](img/diagram.png)\n")
	assertContains(t, html, `<base href="file:///tmp/docs/">`, `src="img/diagram.png"`)
}
