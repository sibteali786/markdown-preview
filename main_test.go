package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	resultFile = "test1.md.html"
	goldenFile = "./testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}
	// Pass the filename to parseContent
	result, err := parseContent(input, "", "test1.md")
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}
	// Trim whitespace from result and expected
	result = bytes.TrimSpace(result)
	expected = bytes.TrimSpace(expected)
	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("Result content does not match golden file")
	}
}

func TestRunWithFile(t *testing.T) {
	var mockStdOut bytes.Buffer
	// Test the run function with file input
	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}
	fileNameOnly := "test1.md"
	if err := run(input, "", fileNameOnly, &mockStdOut, true); err != nil {
		t.Fatal(err)
	}
	resultFile := strings.TrimSpace(mockStdOut.String())
	result, err := os.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}
	// Trim whitespace from result and expected
	result = bytes.TrimSpace(result)
	expected = bytes.TrimSpace(expected)
	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("Result content does not match golden file")
	}
	os.Remove(resultFile)
}

func TestRunWithStdin(t *testing.T) {
	// Mock STDIN input
	stdinContent := "# Test Markdown File\n\nThis is a test."
	r, w, _ := os.Pipe()
	defer r.Close()
	os.Stdin = r
	go func() {
		w.Write([]byte(stdinContent))
		w.Close()
	}()

	// Mock STDOUT
	var mockStdOut bytes.Buffer

	// Run the tool with STDIN
	if err := run([]byte(stdinContent), "", "STDIN", &mockStdOut, true); err != nil {
		t.Fatal(err)
	}

	resultFile := strings.TrimSpace(mockStdOut.String())
	result, err := os.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}

	// Expected output from the golden file
	expected, err := parseContent([]byte(stdinContent), "", "STDIN")
	if err != nil {
		t.Fatal(err)
	}

	// Trim whitespace from result and expected
	result = bytes.TrimSpace(result)
	expected = bytes.TrimSpace(expected)
	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("Result content does not match golden file")
	}
	os.Remove(resultFile)
}
