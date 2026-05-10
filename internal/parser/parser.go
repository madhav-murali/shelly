package parser

import "strings"

func ParseLine(line string) []string {
	var inQuotes bool
	var current strings.Builder
	var args []string
	for _, runeVal := range line {
		switch runeVal {
		case '\'':
			inQuotes = !inQuotes
			continue
		case ' ':
			if inQuotes {
				current.WriteRune(runeVal)
			} else {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(runeVal)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}
