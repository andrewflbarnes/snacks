package strs

import "unicode"

func IsPrintable(bytes []byte) bool {
	for _, c := range string(bytes) {
		if !unicode.IsPrint(c) &&
			c != '\n' &&
			c != '\r' &&
			c != '\t' {
			return false
		}
	}
	return true
}
