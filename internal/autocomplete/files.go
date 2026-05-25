package auto

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"golang.org/x/term"
)

func getLSP(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, entry := range strs {
		for !strings.HasPrefix(entry, prefix) {
			prefix = prefix[:len(prefix)-1]
		}
	}
	return prefix
}

func multiReturn(t *term.Terminal, matches []string, cmdPrefix string, word string, line string, pos int, tabCount *int16, lastLine *string) (newLine string, newPos int, ok bool) {
	if line != *lastLine {
		*tabCount = 0
	}
	*tabCount++
	*lastLine = line
	switch *tabCount {
	case 1:
		lcp := getLSP(matches)
		if len(lcp) > len(word) {
			*tabCount = 0
			completion := cmdPrefix + lcp
			return completion, len(completion), true
		}
		fmt.Fprintf(t, "\x07")
		return line, pos, false
	case 2:
		*tabCount = 0
		slices.Sort(matches)
		fmt.Fprintf(t, "$ %s", line)
		fmt.Fprintf(t, "\r\n%s\r\n", strings.Join(matches, "  "))
		return line, pos, false
	}
	return line, pos, false
}

func FileCompletetion(t *term.Terminal, line string, pos int, key rune, tabCount *int16, lastLine *string) (newLine string, newPos int, ok bool) {
	lastSpace := strings.LastIndex(line, " ")
	if lastSpace == -1 {
		return line, pos, false
	}
	// if line is: ls p -> cmdPrefix = ls and word = p
	cmdPrefix := line[:lastSpace+1]
	word := line[lastSpace+1:]

	searchDir := "."
	filePrefix := word

	if strings.Contains(word, "/") {
		idx := strings.LastIndex(word, "/")
		searchDir = word[:idx+1]
		filePrefix = word[idx+1:]
	}

	entries, err := os.ReadDir(searchDir)
	if err != nil {
		fmt.Fprintf(t, "\x07")
		return line, pos, false
	}

	var matches []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), filePrefix) {
			matchPath := searchDir + entry.Name()

			if searchDir == "." {
				matchPath = entry.Name()
			}
			if entry.IsDir() {
				matchPath += "/"
			}

			matches = append(matches, matchPath)
		}
	}

	switch len(matches) {
	case 0:
		break
	case 1:
		completion := matches[0]
		newLine = cmdPrefix + completion
		if !strings.HasSuffix(newLine, "/") {
			newLine += " "
		}
		return newLine, len(newLine), true
	default:
		//fmt.Fprintf(t, "matches has len: %d", len(matches))
		return multiReturn(t, matches, cmdPrefix, word, line, pos, tabCount, lastLine)
	}

	fmt.Fprintf(t, "\x07")
	return line, pos, false
}
