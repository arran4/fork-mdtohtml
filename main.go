package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"html"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed usage.tmpl
var usageTmpl string

func check(e error) {
	if e != nil {
		log.Printf("error: %v", e)
		os.Exit(1)
	}
}

func usage() {
	tmpl, err := template.New("usage").Parse(usageTmpl)
	if err != nil {
		log.Printf("error parsing usage template: %v", err)
		os.Exit(1)
	}
	data := struct {
		ProgName string
	}{
		ProgName: os.Args[0],
	}
	err = tmpl.Execute(log.Writer(), data)
	if err != nil {
		log.Printf("error executing usage template: %v", err)
	}
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	var fname string
	var noCSS bool
	var title string
	var headers []string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-h" || arg == "--help" || arg == "-help" {
			usage()
		} else if arg == "-nocss" {
			noCSS = true
		} else if arg == "-title" {
			if i+1 < len(args) {
				title = args[i+1]
				i++
			} else {
				log.Printf("error: missing argument for -title")
				usage()
			}
		} else if arg == "-header" {
			if i+1 < len(args) {
				headers = append(headers, args[i+1])
				i++
			} else {
				log.Printf("error: missing argument for -header")
				usage()
			}
		} else if strings.HasPrefix(arg, "-") {
			log.Printf("error: unknown flag %s", arg)
			usage()
		} else if fname != "" {
			log.Printf("error: multiple input files specified")
			usage()
		} else {
			fname = arg
		}
	}

	if fname == "" {
		log.Printf("error: missing markdown filename")
		usage()
	}

	ext := filepath.Ext(fname)
	if strings.ToLower(ext) != ".md" {
		log.Printf("error: input file must be a markdown file (.md)")
		os.Exit(1)
	}

	base := strings.TrimSuffix(fname, ext)

	wfile, err := os.Create(base + ".html")
	check(err)

	defer func() {
		if err := wfile.Close(); err != nil {
			log.Printf("error: %v", err)
			os.Exit(1)
		}
	}()
	writer := bufio.NewWriter(wfile)
	if title != "" {
		escapedTitle := html.EscapeString(title)
		_, err = fmt.Fprintf(writer, "<title>%s</title>\n", escapedTitle)
		check(err)
	}
	for _, h := range headers {
		_, err = fmt.Fprintln(writer, h)
		check(err)
	}
	if !noCSS {
		_, err = fmt.Fprintln(writer, css())
		check(err)
	}

	rfile, err := os.Open(fname)
	check(err)

	defer func() {
		if err := rfile.Close(); err != nil {
			log.Printf("error closing input file: %v", err)
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
