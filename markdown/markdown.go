package markdown

import (
	"fmt"
	"strings"
)

var (
	// is thread safe/goroutine safe
	markdownReplacer = strings.NewReplacer(
		"\\", "\\\\",
		"`", "\\`",
		"*", "\\*",
		"_", "\\_",
		"{", "\\{",
		"}", "\\}",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		".", "\\.",
		"!", "\\!",
	)
)

// Escape user input outside of inline code blocks
func Escape(userInput string) string {
	return markdownReplacer.Replace(userInput)
}

// WrapInInlineCodeBlock puts the user input into a inline codeblock that is properly escaped.
func WrapInInlineCodeBlock(text string) string {
	return WrapInCustom(text, "`")
}

func WrapInFat(text string) string {
	return WrapInCustom(text, "**")
}

func WrapInCustom(text, wrap string) (result string) {
	if text == "" {
		return ""
	}

	numWraps := strings.Count(text, wrap) + 1
	result = text
	for idx := 0; idx < numWraps; idx++ {
		result = fmt.Sprintf("%s%s%s", wrap, result, wrap)
	}
	return result
}
