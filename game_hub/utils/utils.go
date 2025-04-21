package utils

import (
	"fmt"
	"strings"
	"unicode"
)

func WrapText(input string, width int) string {
	if width <= 0 {
		return input
	}

	lines := strings.SplitAfter(input, "\n")
	var result strings.Builder

	for _, line := range lines {
		trailing := ""
		if strings.HasSuffix(line, "\r\n") {
			trailing = "\r\n"
			line = strings.TrimSuffix(line, "\r\n")
		} else if strings.HasSuffix(line, "\n") {
			trailing = "\n"
			line = strings.TrimSuffix(line, "\n")
		} else if strings.HasSuffix(line, "\r") {
			trailing = "\r"
			line = strings.TrimSuffix(line, "\r")
		}

		if len(line) == 0 {
			result.WriteString(trailing)
			continue
		}

		words := strings.Fields(line)
		if len(words) == 0 {
			result.WriteString(trailing)
			continue
		}

		currentLineLen := 0
		for _, word := range words {
			wordLen := len([]rune(word))

			if currentLineLen+wordLen+1 > width && currentLineLen > 0 {
				result.WriteString("\n")
				currentLineLen = 0
			}

			if currentLineLen > 0 {
				result.WriteString(" ")
				currentLineLen++
			}

			result.WriteString(word)
			currentLineLen += wordLen
		}

		result.WriteString(trailing)
	}

	return result.String()
}
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
