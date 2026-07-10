package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func check(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", e)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <markdown-filename> [-nocss]\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	var fname string
	var noCSS bool

	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "--help" || arg == "-help" {
			usage()
		} else if arg == "-nocss" {
			noCSS = true
		} else if !strings.HasPrefix(arg, "-") {
			fname = arg
		}
	}

	if fname == "" {
		fmt.Fprintf(os.Stderr, "error: missing markdown filename\n")
		usage()
	}

	ext := filepath.Ext(fname)
	if strings.ToLower(ext) != ".md" {
		fmt.Fprintf(os.Stderr, "error: input file must be a markdown file (.md)\n")
		os.Exit(1)
	}

	base := strings.TrimSuffix(fname, ext)

	wfile, err := os.Create(base + ".html")
	check(err)

	defer func() {
		if err := wfile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}()
	writer := bufio.NewWriter(wfile)
	if !noCSS {
		_, err = fmt.Fprintln(writer, css())
		check(err)
	}

	rfile, err := os.Open(fname)
	check(err)

	defer func() {
		if err := rfile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}()
	reader := bufio.NewReader(rfile)

	lines := make([]Line, 0)
	for {
		line, err := reader.ReadString('\n')
		if err != nil { // io.EOF
			break
		}
		// Remove newline at the end of line
		if len(line) > 1 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		lines = append(lines, convert(line))
	}

	_, _ = writer.WriteString("<body>")
	_, _ = writer.WriteString(generate(lines))
	_, _ = writer.WriteString("</body>")
	err = writer.Flush()
	check(err)
}
