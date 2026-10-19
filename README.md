# marks

A simple CLI tool written in Go to preview Markdown files in your browser with GitHub-style rendering.

`marks` reads a markdown file, converts it to HTML, applies GitHub markdown styling, creates a temporary HTML file, and opens it automatically in your default browser.

---

## Features

- Preview Markdown files in browser
- GitHub-style markdown rendering (light and dark, follows your system)
- GitHub Flavored Markdown: tables, task lists, strikethrough, autolinks, footnotes, definition lists
- GitHub alerts (`> [!NOTE]`, `> [!TIP]`, `> [!IMPORTANT]`, `> [!WARNING]`, `> [!CAUTION]`)
- Syntax highlighting for fenced code blocks, with a copy button
- Mermaid diagrams (```` ```mermaid ````), click a diagram to toggle full size
- Sticky table of contents with current-section highlight
- YAML front matter shown as a small table
- Relative images and links resolve against the markdown file's folder
- No local server required
- Automatic browser opening
- Lightweight and fast
- Written in Go

---

## Installation

## Option 1: Install from binary release

Download the latest release from:

https://github.com/hitalloazevedo/marks/releases/latest

After downloading:

```bash
chmod +x marks
mv marks ~/.local/bin/marks
```

Make sure `~/.local/bin` is in your PATH.

### zsh

```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### bash

```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

---

## Option 2: Build locally

Clone repository:

```bash
git clone https://github.com/hitalloazevedo/marks.git
cd marks
```

Install dependencies:

```bash
go mod tidy
```

Run install script:

```bash
chmod +x install.sh
./install.sh
```

This script will:

- build the binary
- install it to:

```bash
~/.local/bin/marks
```

---

## Usage

Preview a markdown file:

```bash
marks README.md
```

Preview any markdown file:

```bash
marks docs/setup.md
```

What happens:

1. Reads markdown file
2. Converts markdown to HTML
3. Applies GitHub markdown CSS styling
4. Creates temporary HTML file
5. Opens file in your default browser

No server is started.

---

## Example

```bash
marks notes.md
```

The browser opens a temporary rendered HTML preview automatically.

---

## Development

Run directly without installation:

```bash
go run . README.md
```

Build manually:

```bash
go build -o marks
```

---

## Dependencies

- github.com/yuin/goldmark
- github.com/alecthomas/chroma/v2 (syntax highlighting)
- github.com/pkg/browser

Install with:

```bash
go mod download
```

Mermaid is loaded from the jsDelivr CDN only when a document contains a diagram. Offline, the diagram source is shown instead.

---

## Project structure

```bash
.
├── marks.go          # CLI: read file, write temp HTML, open browser
├── render.go         # markdown → HTML (extensions, highlighting, alerts, TOC)
├── render_test.go
├── install.sh
├── static/
│   ├── github-markdown.css
│   ├── marks.css     # layout, tables, code blocks, diagrams
│   ├── marks.js      # anchors, copy buttons, TOC highlight
│   └── page.html     # page template
```

---

## Future improvements

- live reload on file changes
- math (KaTeX)
- custom themes
- watch mode
- export to HTML file

---

## License

MIT

## Screenshot
<img width="1920" height="965" alt="image" src="https://github.com/user-attachments/assets/b07e91e9-bc74-4818-9ee4-e0468c3334f7" />

