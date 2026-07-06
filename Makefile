files=regexp.go gen.go css.go

mdtohtml: main.go $(files)
	go build -o mdtohtml main.go $(files)

test: main_test.go mdtohtml
	go test -v .

clean:
	rm mdtohtml
