// Package strs provides useful string functions. It is named strs to prevent it being
// confused with the standard strings lib
package strs

import "unicode"

// IsPrintable returns true if the bytes in a slice are all printable unicode characters
// or one of newline, carriage return or tab. The default encoding (UTF-8) is expected.
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
