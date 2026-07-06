package main

import (
	"embed"
	"io/fs"
	"path"
	"strings"
	"testing"
	"testing/fstest"

	"golang.org/x/tools/txtar"
)

//go:embed testdata/txtar/*.txtar
var testdataFS embed.FS

func SplitInputExpected(ar *txtar.Archive) (input, expected fstest.MapFS) {
	input = fstest.MapFS{}
	expected = fstest.MapFS{}

	for _, f := range ar.Files {
		switch {
		case f.Name == "input.txt":
			input[f.Name] = &fstest.MapFile{Data: f.Data}
		case f.Name == "expected.html":
			expected[f.Name] = &fstest.MapFile{Data: f.Data}
		}
	}
	return input, expected
}

func TestTxtar(t *testing.T) {
	entries, err := fs.Glob(testdataFS, "testdata/txtar/*.txtar")
	if err != nil {
		t.Fatalf("glob fixtures: %v", err)
	}

	for _, fixture := range entries {
		fixture := fixture
		t.Run(strings.TrimSuffix(path.Base(fixture), ".txtar"), func(t *testing.T) {
			raw, err := testdataFS.ReadFile(fixture)
			if err != nil {
				t.Fatalf("read fixture %s: %v", fixture, err)
			}
			ar := txtar.Parse(raw)

			inputFS, expectedFS := SplitInputExpected(ar)

			inputRaw, err := fs.ReadFile(inputFS, "input.txt")
			if err != nil {
				t.Fatalf("read input.txt: %v", err)
			}
			input := strings.TrimSpace(string(inputRaw))

			lines := make([]Line, 0)
			for _, in := range strings.Split(input, "\n") {
				lines = append(lines, convert(in))
			}
			html := generate(lines)

			expectedRaw, err := fs.ReadFile(expectedFS, "expected.html")
			if err != nil {
				t.Fatalf("read expected.html: %v", err)
			}
			expected := strings.TrimSpace(string(expectedRaw))

			if html != expected {
				t.Errorf("%q => expected %q but got %q", input, expected, html)
			}
		})
	}
}
