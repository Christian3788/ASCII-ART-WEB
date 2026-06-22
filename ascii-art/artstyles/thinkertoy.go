package artstyles

import "strings"

func Thinkertoy() map[rune][6]string {
	base := Standard()
	thinker := make(map[rune][6]string, len(base))
	for ch, pattern := range base {
		var charLines [6]string
		for i, line := range pattern {
			charLines[i] = strings.Map(func(r rune) rune {
				if r == ' ' {
					return ' '
				}
				return '@'
			}, line)
		}
		thinker[ch] = charLines
	}
	return thinker
}
