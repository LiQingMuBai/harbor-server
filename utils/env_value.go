package utils

import "strings"

func StripInlineComment(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	for i := 0; i < len(value); i++ {
		if value[i] != '#' {
			continue
		}
		if i == 0 {
			return ""
		}
		prev := value[i-1]
		if prev == ' ' || prev == '\t' {
			return strings.TrimSpace(value[:i])
		}
	}
	return value
}

