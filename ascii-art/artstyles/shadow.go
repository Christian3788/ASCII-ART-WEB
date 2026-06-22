package artstyles

import "strings"

func Shadow() map[rune][6]string {
	base := Standard()
	shadow := make(map[rune][6]string, len(base))
	for ch, pattern := range base {
		var charLines [6]string
		for i, line := range pattern {
			shadowLine := strings.Map(func(r rune) rune {
				if r == ' ' {
					return ' '
				}
				return '.'
			}, line)
			charLines[i] = line + "  " + shadowLine
		}
		shadow[ch] = charLines
	}
	return shadow
}
