package main

import (
	"regexp"
	"testing"
)

var (
	headingRegExp    = regexp.MustCompile(`(^#{1,6}) (.+)`)
	headingInRegExp  = regexp.MustCompile(`^ *- +(#{1,6}) (.+)`)
	listRegExp       = regexp.MustCompile(`^( *)- (.+)`)
	linkRegExp       = regexp.MustCompile(`.*(\[.+?\])(\(.+?\)).*`)
	emphasisRegExp   = regexp.MustCompile(`.*(\*.+\*).*|.*(\_.+\_).*`)
	strongRegExp     = regexp.MustCompile(`.*(\*\*.+\*\*).*|.*(\_\_.+\_\_).*`)
	horizontalRegExp = regexp.MustCompile(`^-{3}|_{3}|\*{3}`)
	whitespaceRegExp = regexp.MustCompile(`^( +)(.*)`)
)

func TestMatchStrong(t *testing.T) {
	tests := []string{
		"**hello**",
		"__hello__",
		"a **b** c",
		"**a** and **b**",
		"a __b__ c",
		"__a__ and __b__",
		"no strong",
		"**",
		"****",
		"**a**",
		"**a**** and **b**", // stack approach testing
	}

	for _, tc := range tests {
		expected := strongRegExp.FindStringSubmatchIndex(tc)
		actual := matchStrong(tc)
		if !slicesEqual(expected, actual) {
			t.Errorf("matchStrong(%q): expected %v, got %v", tc, expected, actual)
		}
	}
}

func TestMatchEmphasis(t *testing.T) {
	tests := []string{
		"*hello*",
		"_hello_",
		"a *b* c",
		"*a* and *b*",
		"a _b_ c",
		"_a_ and _b_",
		"no emphasis",
		"*",
		"**",
		"*a*",
		"*a*b*", // stack testing
	}

	for _, tc := range tests {
		expected := emphasisRegExp.FindStringSubmatchIndex(tc)
		actual := matchEmphasis(tc)
		if !slicesEqual(expected, actual) {
			t.Errorf("matchEmphasis(%q): expected %v, got %v", tc, expected, actual)
		}
	}
}

func TestMatchLink(t *testing.T) {
	tests := []string{
		"[a](b)",
		"x [a](b) y",
		"[a](b) and [c](d)",
		"no link",
		"[]()",
		"[a]()",
		"[](b)",
		"[a](b) and [c](d", // test backwards iteration for anchors
	}

	for _, tc := range tests {
		expected := linkRegExp.FindStringSubmatchIndex(tc)
		actual := matchLink(tc)
		if !slicesEqual(expected, actual) {
			t.Errorf("matchLink(%q): expected %v, got %v", tc, expected, actual)
		}
	}
}

func TestMatchHeadingIn(t *testing.T) {
	tests := []string{
		"- # h1",
		"  - ## h2",
		" - ###### h6",
		" - ####### h7", // invalid
		"- # ", // invalid
		"- # h",
		"- h1",
		" -  # h1",
	}

	for _, tc := range tests {
		expected := headingInRegExp.FindStringSubmatchIndex(tc)
		actual := matchHeadingIn(tc)
		if !slicesEqual(expected, actual) {
			t.Errorf("matchHeadingIn(%q): expected %v, got %v", tc, expected, actual)
		}
	}
}

func TestMatchList(t *testing.T) {
	tests := []string{
		"- a",
		"  - b",
		"- ",
		" -a",
		" - ",
	}

	for _, tc := range tests {
		expected := listRegExp.FindStringSubmatchIndex(tc)
		actual := matchList(tc)
		if !slicesEqual(expected, actual) {
			t.Errorf("matchList(%q): expected %v, got %v", tc, expected, actual)
		}
	}
}

func TestMatchHeading(t *testing.T) {
	tests := []string{
		"# h1",
		"## h2",
		"###### h6",
		"####### h7",
		"# ",
		"# h",
		" # h1",
	}

	for _, tc := range tests {
		expected := headingRegExp.FindStringSubmatchIndex(tc)
		actual := matchHeading(tc)
		if !slicesEqual(expected, actual) {
			t.Errorf("matchHeading(%q): expected %v, got %v", tc, expected, actual)
		}
	}
}

func TestMatchHorizontal(t *testing.T) {
	tests := []string{
		"---",
		"___",
		"***",
		"a---",
		"a___",
		"a***",
		"--",
		"__",
		"**",
	}

	for _, tc := range tests {
		expected := horizontalRegExp.MatchString(tc)
		actual := matchHorizontal(tc)
		if expected != actual {
			t.Errorf("matchHorizontal(%q): expected %v, got %v", tc, expected, actual)
		}
	}
}

func TestMatchWhitespace(t *testing.T) {
	tests := []string{
		" a",
		"  a",
		"a",
		" ",
	}

	for _, tc := range tests {
		expected := whitespaceRegExp.FindStringSubmatchIndex(tc)
		actual := matchWhitespace(tc)
		if !slicesEqual(expected, actual) {
			t.Errorf("matchWhitespace(%q): expected %v, got %v", tc, expected, actual)
		}
	}
}

func slicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
