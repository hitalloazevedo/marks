# marks

A simple CLI tool written in Go to preview Markdown files in your browser.

`marks` reads a markdown file, converts it to HTML, serves it locally, and automatically opens your default browser.

Repository: https://github.com/hitalloazevedo/marks

---

## Features

- Preview Markdown in browser
- Local HTTP server
- Automatic browser opening
- Lightweight and fast
- Written in Go

---

## Installation

### Clone repository

```bash
git clone https://github.com/hitalloazevedo/marks.git
cd marks
```

### Install dependencies

```bash
go mod tidy
```

### Run install script

```bash
chmod +x install.sh
./install.sh
```

This will:

- build the binary
- install it to:

```bash
~/.local/bin/marks
```

---

## PATH setup

If `~/.local/bin` is not in your PATH, add this to your shell config.

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

## Usage

Preview a markdown file:

```bash
marks README.md
```

Preview any markdown file:

```bash
marks docs/setup.md
```

This will:

1. convert markdown to HTML
2. start a local server
3. open browser automatically

Default URL:

```text
http://localhost:8080
```

---

## Example

```bash
marks notes.md
```

Browser opens:

```text
http://localhost:8080
```

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
- github.com/pkg/browser

Install manually:

```bash
go get github.com/yuin/goldmark
go get github.com/pkg/browser
```

---

## Future improvements

- live reload on file changes
- syntax highlighting
- GitHub markdown styling
- custom port support
- dark mode
- configurable theme

---

## License

MIT
