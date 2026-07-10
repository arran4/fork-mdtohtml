package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Type int

// only block elements
const (
	Newline Type = iota
	P
	H1
	H2
	H3
	H4
	H5
	H6
	Li
	Hr
)

type Line struct {
	ty    Type
	val   string
	dep   int
	hasBr bool
}

func matchHeadingIn(line string) []int {
	idx := 0
	for idx < len(line) && line[idx] == ' ' {
		idx++
	}
	if idx >= len(line) || line[idx] != '-' {
		return nil
	}
	idx++
	if idx >= len(line) || line[idx] != ' ' {
		return nil
	}
	for idx < len(line) && line[idx] == ' ' {
		idx++
	}
	hashStart := idx
	hashCount := 0
	for idx < len(line) && line[idx] == '#' {
		hashCount++
		idx++
	}
	if hashCount < 1 || hashCount > 6 {
		return nil
	}
	hashEnd := idx
	if idx >= len(line) || line[idx] != ' ' {
		return nil
	}
	idx++
	if idx >= len(line) {
		return nil
	}
	return []int{0, len(line), hashStart, hashEnd, idx, len(line)}
}

func matchList(line string) []int {
	idx := 0
	for idx < len(line) && line[idx] == ' ' {
		idx++
	}
	spaces := idx
	if idx >= len(line) || line[idx] != '-' {
		return nil
	}
	idx++
	if idx >= len(line) || line[idx] != ' ' {
		return nil
	}
	idx++
	if idx >= len(line) {
		return nil
	}
	return []int{0, len(line), 0, spaces, idx, len(line)}
}

func matchHeading(line string) []int {
	idx := 0
	hashCount := 0
	for idx < len(line) && line[idx] == '#' {
		hashCount++
		idx++
	}
	if hashCount < 1 || hashCount > 6 {
		return nil
	}
	if idx >= len(line) || line[idx] != ' ' {
		return nil
	}
	idx++
	if idx >= len(line) {
		return nil
	}
	return []int{0, len(line), 0, hashCount, idx, len(line)}
}

func matchHorizontal(line string) bool {
	if strings.HasPrefix(line, "---") {
		return true
	}
	if strings.Contains(line, "___") {
		return true
	}
	if strings.Contains(line, "***") {
		return true
	}
	return false
}

func matchWhitespace(line string) []int {
	idx := 0
	for idx < len(line) && line[idx] == ' ' {
		idx++
	}
	if idx > 0 {
		return []int{0, len(line), 0, idx, idx, len(line)}
	}
	return nil
}

func matchStrong(line string) []int {
	idx1 := strings.LastIndex(line, "**")
	if idx1 != -1 {
		idx2 := strings.LastIndex(line[:idx1], "**")
		if idx2 != -1 && idx1-idx2 > 2 {
			return []int{0, len(line), idx2, idx1 + 2, -1, -1}
		}
	}
	idx1 = strings.LastIndex(line, "__")
	if idx1 != -1 {
		idx2 := strings.LastIndex(line[:idx1], "__")
		if idx2 != -1 && idx1-idx2 > 2 {
			return []int{0, len(line), -1, -1, idx2, idx1 + 2}
		}
	}
	return nil
}

func matchEmphasis(line string) []int {
	idx1 := strings.LastIndex(line, "*")
	if idx1 != -1 {
		idx2 := strings.LastIndex(line[:idx1], "*")
		if idx2 != -1 && idx1-idx2 > 1 {
			return []int{0, len(line), idx2, idx1 + 1, -1, -1}
		}
	}
	idx1 = strings.LastIndex(line, "_")
	if idx1 != -1 {
		idx2 := strings.LastIndex(line[:idx1], "_")
		if idx2 != -1 && idx1-idx2 > 1 {
			return []int{0, len(line), -1, -1, idx2, idx1 + 1}
		}
	}
	return nil
}

func matchLink(line string) []int {
	anchor := strings.LastIndex(line, "](")
	if anchor == -1 {
		return nil
	}
	start := strings.LastIndex(line[:anchor], "[")
	if start == -1 {
		return nil
	}
	end := strings.Index(line[anchor:], ")")
	if end == -1 {
		return nil
	}
	end += anchor
	if anchor-start > 1 && end-(anchor+1) > 1 {
		return []int{0, len(line), start, anchor + 1, anchor + 1, end + 1}
	}
	return nil
}

func ntoh(n int) Type {
	switch n {
	case 1:
		return H1
	case 2:
		return H2
	case 3:
		return H3
	case 4:
		return H4
	case 5:
		return H5
	case 6:
		return H6
	default:
		panic(fmt.Sprintf("a heading should be in the range of 1 to 6, but got %d", n))
	}

}

func hton(ty Type) int {
	switch ty {
	case H1:
		return 1
	case H2:
		return 2
	case H3:
		return 3
	case H4:
		return 4
	case H5:
		return 5
	case H6:
		return 6
	default:
		panic(fmt.Sprintf("a heading should be in the range of 1 to 6, but got %d", ty))
	}

}

func convert(line string) Line {
	// newline
	if line == "\n" || len(line) == 0 {
		return Line{Newline, " ", 0, false}
	}

	hasBr := false

	// ----- Inline Elements -----

	matchSomething := true
	for matchSomething {
		matchSomething = false

		// inline elements are replaced with HTML in this function.
		if loc := matchStrong(line); loc != nil {
			s := loc[2]
			e := loc[3]
			if s == -1 && e == -1 {
				s = loc[4]
				e = loc[5]
			}
			sttag := "<strong>" + line[s+2:e-2] + "</strong>"
			line = line[:s] + sttag + line[e:]
			matchSomething = true
			continue
		}

		if loc := matchEmphasis(line); loc != nil {
			s := loc[2]
			e := loc[3]
			if s == -1 && e == -1 {
				s = loc[4]
				e = loc[5]
			}
			emtag := "<em>" + line[s+1:e-1] + "</em>"
			line = line[:s] + emtag + line[e:]
			matchSomething = true
			continue
		}

		if loc := matchLink(line); loc != nil { // links <a href="url">link text</a>
			text := line[loc[2]+1 : loc[3]-1]
			url := line[loc[4]+1 : loc[5]-1]

			litag := "<a href=\"" + url + "\">" + text + "</a>"
			line = line[:loc[2]] + litag + line[loc[5]:]
			fmt.Println(loc)
			fmt.Println(text)
			fmt.Println(url)
			fmt.Println(line)
			matchSomething = true
			continue
		}

		// heading existing in another component
		if loc := matchHeadingIn(line); loc != nil {
			n := loc[3] - loc[2] // heading number
			htag := "<h" + strconv.Itoa(n) + ">" + line[loc[4]:loc[5]] + "</h" + strconv.Itoa(n) + ">"
			line = line[:loc[2]] + htag
			matchSomething = true
			continue
		}

		// break at the end of line
		if len(line) >= 2 && line[len(line)-2:] == "  " {
			line = line[:len(line)-2] + "<br>"
			hasBr = true
			// We don't continue or set matchSomething to true because
			// this should just process the line breaks at the end.
			// Actually the original logic set `matchSomething = false` after this,
			// so it wouldn't loop unless inline elements triggered `continue`.
			// Since we changed it to check everything without looping internally over `MatchString`,
			// this handles the break at the very end of processing.
		} else if len(line) >= 1 && line[len(line)-1] == '\\' {
			line = line[:len(line)-1] + "<br>"
			hasBr = true
		}
	}

	// ----- Block Elements -----

	// block elements will be replaced with HTML in the generate().
	if loc := matchList(line); loc != nil {
		dep := loc[3] / 2
		return Line{Li, line[loc[4]:loc[5]], dep, false}
	}

	if loc := matchHeading(line); loc != nil {
		n := loc[3]
		return Line{ntoh(n), line[loc[4]:loc[5]], 0, false}
	}

	if matchHorizontal(line) {
		return Line{Hr, "", 0, false}
	}

	// replace white spaces with a white space at the start of a line
	if loc := matchWhitespace(line); loc != nil {
		line = " " + line[loc[4]:loc[5]]
	}

	return Line{P, line, 0, hasBr}
}
