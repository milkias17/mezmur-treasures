package db

import "unicode"

func isEnglish(text string) bool {
	for _, r := range text {
		if r > unicode.MaxASCII {
			return false
		}
	}

	return true
}
