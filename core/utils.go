package core

import (
	"fmt"
	"strings"
	"unicode"
)

func SubstituteParams(message string, details map[string]any) string {
	if details == nil {
		return message
	}
	for key, value := range details {
		message = strings.ReplaceAll(message, "$"+key, fmt.Sprintf("%v", value))
	}
	return message
}

func Capitalize(str string) string {
	if len(str) == 0 {
		return str
	}
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
