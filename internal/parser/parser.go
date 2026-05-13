package parser

import "strings"

// Include the functionality of single quotes; currently implementing double quotes integ
func ParseLine(line string) []string {
	var inQuotes bool
	var inDQuotes bool
	var current strings.Builder
	var args []string

	for i := 0; i < len(line); i++ {
		runeVal := line[i]
		// if i > 0 && line[i-1] == '\\' && !inQuotes {
		// 	current.WriteByte(runeVal)
		// 	continue
		// }

		if inQuotes {
			if runeVal == '\'' {
				inQuotes = false
			} else {
				current.WriteByte(runeVal)
			}
			continue
		}

		if runeVal == '\\' {
			if i+1 < len(line) {
				nextChar := line[i+1]
				if inDQuotes {
					if nextChar == '$' || nextChar == '"' || nextChar == '`' || nextChar == '\n' || nextChar == '\\' {
						current.WriteByte(nextChar)
						i++
					} else {
						current.WriteByte('\\')
					}
				} else {
					current.WriteByte(nextChar)
					i++
				}
			}
			continue
		}

		switch runeVal {
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
				current.WriteByte(runeVal)
			} else {
				inQuotes = !inQuotes
			}
		case ' ':
			if inQuotes || inDQuotes {
				current.WriteByte(runeVal)
			} else {
				if current.Len() == 0 {
					continue
				}
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(runeVal)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}
