package textutil

import (
	"errors"
	"strings"
)

func TruncateByRune(s string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}

	return string(runes[:maxRunes])
}

func ExtractJSONObject(content string) (string, error) {
	content = strings.TrimSpace(content)
	content = TrimMarkdownFence(content)

	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start < 0 || end < 0 || end <= start {
		return "", errors.New("json object not found")
	}

	return content[start : end+1], nil
}

func TrimMarkdownFence(content string) string {
	content = strings.TrimSpace(content)

	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimSuffix(content, "```")
		return strings.TrimSpace(content)
	}

	if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
		return strings.TrimSpace(content)
	}

	return content
}
