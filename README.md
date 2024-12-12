
# Markdown Preview Tool (`mdp`)

`mdp` is a lightweight and cross-platform command-line tool for previewing Markdown files as HTML. It supports flexible input methods, including file-based input and `STDIN`, and allows customization of HTML templates via flags or environment variables.

---

## Features

- **Preview Markdown Files:** Convert Markdown to HTML and open it in your default browser.
- **Support for STDIN:** Pipe Markdown content directly to the tool for quick previews.
- **Custom Templates:** Use custom HTML templates to style your output.
- **Cross-Platform Compatibility:** Works on Linux, macOS, and Windows.
- **Safe HTML Output:** Sanitizes output using `bluemonday` for security.
- **Quick Installation:** Install easily via Go.

---

## Installation

To install `mdp`, use the following command:

```bash
go get github.com/sibteali786/mdp
```

---

## Usage

### Basic Usage

Preview a Markdown file:

```bash
mdp -file example.md
```

### Using STDIN

Pipe Markdown content directly:

```bash
echo "# Hello World" | mdp
```

### Skipping Browser Preview

Generate the HTML without opening it in the browser:

```bash
mdp -file example.md -s
```

### Custom Templates

Use a custom HTML template file:

```bash
mdp -file example.md -t /path/to/template.html
```

Set a default template via an environment variable:

```bash
export MDP_TEMPLATE=/path/to/template.html
mdp -file example.md
```

### Priority of Input Methods

1. **`-file` flag:** If specified, the file is used as input.
2. **`STDIN`:** If no file is specified, `STDIN` is used.
3. If neither is provided, the tool exits with a usage message.

---

## Example

### Input Markdown (`example.md`)

```markdown
# Markdown Preview Tool

This is a simple tool to preview Markdown files as HTML.

## Features:

- Convert Markdown to HTML
- Use custom templates
- Secure output

## Example Code Block:
```

go run main.go

```

```

### Output HTML

```html
<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8" />
    <title>Markdown Preview Tool</title>
  </head>
  <body>
    <p>Previewing file: example.md</p>
    <h1>Markdown Preview Tool</h1>
    <p>This is a simple tool to preview Markdown files as HTML.</p>
    <h2>Features:</h2>
    <ul>
      <li>Convert Markdown to HTML</li>
      <li>Use custom templates</li>
      <li>Secure output</li>
    </ul>
    <h2>Example Code Block:</h2>
    <pre><code>go run main.go
        </code></pre>
  </body>
</html>
```

---

## Testing

Run unit and integration tests:

```bash
go test
```

---

## Contributing

Contributions are welcome! Feel free to open issues, suggest features, or submit pull requests.

---

## License

This project is licensed under the [MIT License](LICENSE).

---

## Acknowledgments

- [Blackfriday Markdown Processor](https://github.com/russross/blackfriday)
- [Bluemonday HTML Sanitizer](https://github.com/microcosm-cc/bluemonday)
