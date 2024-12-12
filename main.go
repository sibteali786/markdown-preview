package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="content-type" content="text/html; charset=utf-8">
		<title>{{ .Title }}</title>
	</head>
	<body>
		<p>Previewing file: {{ .FileName }}</p>
		{{ .Body }}
	</body>
</html>
`
)

type content struct {
	Title    string
	Body     template.HTML
	FileName string
}

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto preview")
	tFname := flag.String("t", "", "Alternate template file path")

	// Enhance usage information
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `
Usage:
  mdp [options]

Options:
  -file string
        Path to the Markdown file to preview (required).
  -s    
        Skip opening the preview in a browser (optional).
  -t string
        Path to an alternate HTML template file (optional).

Environment Variables:
  MDP_TEMPLATE
        Path to a default template file. Used when -t is not specified.

Examples:
  Use the default template:
      mdp -file example.md

  Use a custom template file via flag:
      mdp -file example.md -t /path/to/template.html

  Use a custom template via environment variable:
      export MDP_TEMPLATE=/path/to/template.html
      mdp -file example.md

  Skip opening the preview in a browser:
      mdp -file example.md -s
`)

		// Show the default usage text
		flag.PrintDefaults()
	}

	flag.Parse()

	// If user did not provide input file, show usage
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Resolve the template to use
	templatePath := resolveTemplate(*tFname)

	if err := run(*filename, templatePath, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func resolveTemplate(tFname string) string {
	// If the -t flag is provided, use it
	if tFname != "" {
		return tFname
	}

	// If the environment variable is set, use it
	if envTemplate := os.Getenv("MDP_TEMPLATE"); envTemplate != "" {
		return envTemplate
	}

	// Fallback to using the default template
	return ""
}

func run(filename string, templatePath string, out io.Writer, skipPreview bool) error {
	// Read all the data from the input file and check for errors
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Extract only the file name from the full path
	fileNameOnly := filepath.Base(filename)

	htmlData, err := parseContent(input, templatePath, fileNameOnly)
	if err != nil {
		return err
	}

	// Create a temporary file
	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
	defer temp.Close()

	outName := temp.Name()
	fmt.Fprintln(out, outName)

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}
	defer os.Remove(outName)

	return preview(outName)
}

func parseContent(input []byte, templatePath string, fileName string) ([]byte, error) {
	// Parse the markdown file through blackfriday and bluemonday to generate valid and safe HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	var t *template.Template
	var err error

	// Use the custom template file if provided
	if templatePath != "" {
		t, err = template.ParseFiles(templatePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load template file: %w", err)
		}
	} else {
		// Use the default template string
		t, err = template.New("mdp").Parse(defaultTemplate)
		if err != nil {
			return nil, err
		}
	}

	// Instantiate the content with the title, body, and file name
	c := content{
		Title:    "Markdown Preview Tool",
		Body:     template.HTML(body),
		FileName: fileName,
	}

	// Create a buffer of bytes to write to file
	var buffer bytes.Buffer

	// Execute the template with the content
	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func saveHTML(outFname string, data []byte) error {
	// Write the bytes to the file
	return os.WriteFile(outFname, data, 0644)
}

func preview(fname string) error {
	cName := ""
	cParams := []string{}
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	// Append the filename to parameters slice
	cParams = append(cParams, fname)

	// Locate the executable in PATH
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	// Open the file using the default program
	err = exec.Command(cPath, cParams...).Run()

	// Give the browser some time to open the file before it's deleted
	time.Sleep(2 * time.Second)

	return err
}
