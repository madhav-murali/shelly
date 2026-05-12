package parser

import "strings"

// Include the functionality of single quotes; currently implementing double quotes integ
func ParseLine(line string) []string {
	var inQuotes bool
	var inDQuotes bool
	var current strings.Builder
	var args []string
	for i, runeVal := range line {
		if i > 0 && line[i-1] == '\\' {
			current.WriteRune(runeVal)
			continue
		}
		switch runeVal {
		case '\\':
			if inQuotes {
				current.WriteRune(runeVal)
			}
			continue
		case '"':
			if i != len(line)-1 && i != 0 {
				if line[i+1] == '"' {
					continue
				} else if line[i-1] == '"' { // this is stupid
					continue
				}
			}
			inDQuotes = !inDQuotes
		case '\'':
			if inDQuotes {
				current.WriteRune(runeVal)
			} else {
				inQuotes = !inQuotes
			}
		case ' ':
			if inQuotes || inDQuotes {
				current.WriteRune(runeVal)
			} else {
				if current.Len() == 0 {
					continue
				}
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
