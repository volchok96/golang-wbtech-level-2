package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputCh := make(chan string)
	go func() {
		var buf strings.Builder
		io.Copy(&buf, r)
		outputCh <- buf.String()
	}()

	f()

	w.Close()
	output := <-outputCh
	os.Stdout = old
	return output
}

func TestCut(t *testing.T) {
	content := "one\ttwo\tthree\nfour\tfive\tsix\nseven\teight\tnine"
	tmpfile := createTempFile(t, content)
	defer removeTempFile(t, tmpfile)

	tests := []struct {
		options  cutOptions
		expected string
	}{
		{cutOptions{fields: "1", delimiter: "\t", separated: false}, "one\nfour\nseven\n"},
		{cutOptions{fields: "1,3", delimiter: "\t", separated: false}, "one\tthree\nfour\tsix\nseven\tnine\n"},
		{cutOptions{fields: "2", delimiter: "\t", separated: false}, "two\nfive\neight\n"},
		{cutOptions{fields: "1-2", delimiter: "\t", separated: false}, "one\ttwo\nfour\tfive\nseven\teight\n"},
		{cutOptions{fields: "2-3", delimiter: "\t", separated: false}, "two\tthree\nfive\tsix\neight\tnine\n"},
		{cutOptions{fields: "1-3", delimiter: "\t", separated: false}, "one\ttwo\tthree\nfour\tfive\tsix\nseven\teight\tnine\n"},
		{cutOptions{fields: "1,3", delimiter: "\t", separated: true}, "one\tthree\nfour\tsix\nseven\tnine\n"},
		{cutOptions{fields: "1", delimiter: "\t", separated: true}, "one\nfour\nseven\n"}, // Тест на наличие разделителя
	}

	for _, test := range tests {
		input, _ := os.Open(tmpfile)
		output := captureOutput(func() {
			cutLines(input, test.options)
		})

		if output != test.expected {
			t.Errorf("For options %+v, expected %q but got %q", test.options, test.expected, output)
		}
	}
}

func createTempFile(t *testing.T, content string) string {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	return tmpfile.Name()
}

func removeTempFile(t *testing.T, name string) {
	if err := os.Remove(name); err != nil {
		t.Fatal(err)
	}
}
