package main

import (
	"strings"
	"testing"
)

func runTest(t *testing.T, expect string, input string) {
	lines := make([]Line, 0)
	for _, in := range strings.Split(input, "\n") {
		lines = append(lines, convert(in))
	}
	html := generate(lines)

	if html != expect {
		t.Errorf("%q => expected %q but got %q", input, expect, html)
	}
}

func TestParagraph(t *testing.T) {
	runTest(t, "<p>a paragraph</p>", "a paragraph")
	runTest(t, "<p>a paragraph<br>hogehoge</p>", "a paragraph  \nhogehoge")
}

func TestHeading(t *testing.T) {
	runTest(t, "<h1>h1</h1>", "# h1")
	runTest(t, "<h2>h2</h2>", "## h2")
	runTest(t, "<h3>h3</h3>", "### h3")
	runTest(t, "<h4>h4</h4>", "#### h4")
	runTest(t, "<h5>h5</h5>", "##### h5")
	runTest(t, "<h6>h6</h6>", "###### h6")
	runTest(t, "<p>####### h7</p>", "####### h7")
	runTest(t, "<p>###dummyh3</p>", "###dummyh3")
	runTest(t, "<p>C## is not heading</p>", "C## is not heading")
}

func TestList(t *testing.T) {
	runTest(t, "<ul><li>list1</li></ul>", "- list1")
	runTest(t, "<ul><li>list1</li><li>list2</li></ul>", "- list1\n- list2")
	// TODO: Sublist is not a standard syntax.
	// It should be <li>list1<ul><li>sublist1</li></ul></li></ul>
	// but now got <li>list1</li><ul><li>sublist1</li></ul></ul>
	runTest(t, "<ul><li>list1</li><ul><li>sublist1</li></ul></ul>", "- list1\n  - sublist1")
	runTest(t, "<ul><li>list1</li><ul><li>sublist1</li><ul><li>subsublist1</li></ul></ul></ul>", "- list1\n  - sublist1\n    - subsublist1")
	runTest(t, "<ul><li>list1</li><ul><li>sublist1</li></ul><li>list2</li></ul>", "- list1\n  - sublist1\n- list2")
	runTest(t, "<ul><li>a</li><ul><li>aa</li><ul><li>aaa</li></ul></ul><li>b</li></ul>", "- a\n  - aa\n    - aaa\n- b")
	runTest(t, "<ul><li>a</li><ul><li>aa</li><ul><li>aaa</li></ul><li>bb</li></ul></ul>", "- a\n  - aa\n    - aaa\n  - bb")
	runTest(t, "<ul><li><h1>h1</h1></li></ul>", "- # h1")
	//runTest(t, "<ul><li>a -b</li><li>c</li></ul>", "- a\n  -b\n- c")
}

func TestLink(t *testing.T) {
	runTest(t, "<p><a href=\"http://example.com\">link</a></p>", "[link](http://example.com)")
	runTest(t, "<p><a href=\"http://example.com\">link(2)</a></p>", "[link(2)](http://example.com)")
	runTest(t, "<p>inline text<a href=\"http://example.com\">link</a>.</p>", "inline text[link](http://example.com).")
	runTest(t, "<p>[dummylink] (http://example.com)</p>", "[dummylink] (http://example.com)")
}

func TestHeadingWithInlineElements(t *testing.T) {
	runTest(t, "<h1><a href=\"http://example.com\">link</a></h1>", "# [link](http://example.com)")
	runTest(t, "<h1>- dummylist</h1>", "# - dummylist")
}

func TestListWithInlineElements(t *testing.T) {
	runTest(t, "<ul><li><a href=\"http://example.com\">link</a></li></ul>", "- [link](http://example.com)")
	runTest(t, "<ul><li>This is <a href=\"http://example.com\">link</a> list.</li></ul>", "- This is [link](http://example.com) list.")
	runTest(t, "<ul><li><h1>h1</h1></li></ul>", "- # h1")
}

func TestHeadingAfterList(t *testing.T) {
	runTest(t, "<ul><li>list1</li></ul><h1>h1</h1>", "- list1\n# h1")
	runTest(t, "<ul><li>list1</li></ul><h1>h1</h1>", "- list1\n\n# h1")
	runTest(t, "<ul><li>a</li><ul><li>b</li></ul></ul><h1>h1</h1>", "- a\n  - b\n# h1")
}

func TestMultipleLines(t *testing.T) {
	runTest(t, "<h1>h1</h1><p>text</p>", "# h1\ntext")
}

func TestEmphasis(t *testing.T) {
	runTest(t, "<p><em>emphasis</em></p>", "*emphasis*")
	runTest(t, "<p><em>emphasis</em></p>", "_emphasis_")
	runTest(t, "<p><strong>strong</strong></p>", "**strong**")
	runTest(t, "<p><strong>strong</strong></p>", "__strong__")
}
