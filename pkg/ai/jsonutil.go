package ai

import (
	"regexp"
	"strings"
)

// CleanJSON trims common wrappers and artifacts around LLM JSON outputs
// and fixes trailing commas so that standard JSON parsers can accept it.
func CleanJSON(s string) string {
	if s == "" {
		return s
	}

	str := strings.TrimSpace(s)

	// Strip markdown code fences
	str = strings.TrimPrefix(str, "```json")
	str = strings.TrimPrefix(str, "```")
	str = strings.TrimSuffix(str, "```")

	// Normalize whitespace and control chars
	str = strings.TrimSpace(str)
	str = strings.TrimPrefix(str, "\ufeff") // BOM
	str = strings.ReplaceAll(str, "\r", "")

	// Fix common JSON format issues
	str = fixJSONTrailingCommas(str)

	return str
}

// ExtractJSONFromText tries to extract a top-level JSON object from mixed text
// (e.g., includes thoughts, explanations, or markdown), then cleans it.
func ExtractJSONFromText(text string) string {
	if text == "" {
		return text
	}

	t := strings.TrimSpace(text)

	// Prefer fenced blocks if present
	if strings.Contains(t, "```json") && strings.Contains(t, "```") {
		// Take the first fenced block
		start := strings.Index(t, "```json")
		if start >= 0 {
			start += len("```json")
			// Skip following newlines
			for start < len(t) && (t[start] == '\n' || t[start] == '\r') {
				start++
			}
			if endRel := strings.Index(t[start:], "```"); endRel >= 0 {
				candidate := strings.TrimSpace(t[start : start+endRel])
				return CleanJSON(candidate)
			}
		}
	}

	// If contains special prefixes like "> *Thought", fall back to first {...} pair
	if strings.Contains(t, "> *Thought") {
		// fall through to brace scan
	}

	// Fallback: first '{' to last '}'
	l := strings.Index(t, "{")
	r := strings.LastIndex(t, "}")
	if l >= 0 && r >= 0 && r > l {
		candidate := t[l : r+1]
		return CleanJSON(candidate)
	}

	// Nothing better, return cleaned original
	return CleanJSON(t)
}

// fixJSONTrailingCommas removes dangling commas before '}' or ']' so that
// JSON becomes parseable by standard decoders.
func fixJSONTrailingCommas(jsonStr string) string {
	// Remove trailing comma before object end: "value", }
	re1 := regexp.MustCompile(`,(\s*})`)
	jsonStr = re1.ReplaceAllString(jsonStr, "$1")

	// Remove trailing comma before array end: "value", ]
	re2 := regexp.MustCompile(`,(\s*])`)
	jsonStr = re2.ReplaceAllString(jsonStr, "$1")

	return jsonStr
}
