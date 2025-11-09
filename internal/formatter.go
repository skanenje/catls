package internal

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatMode defines supported output types
type FormatMode string

const (
	Markdown FormatMode = "markdown"
	JSON     FormatMode = "json"
)

// FormatEntries formats a list of FileEntry into Markdown or JSON
func FormatEntries(entries []FileEntry, mode FormatMode) (string, error) {
	switch mode {
	case Markdown:
		return formatMarkdown(entries), nil
	case JSON:
		return formatJSON(entries)
	default:
		return "", fmt.Errorf("unsupported format: %s", mode)
	}
}

// ------------------- Markdown -------------------

func formatMarkdown(entries []FileEntry) string {
	var sb strings.Builder

	for _, e := range entries {
		if e.Kind == "dir" {
			sb.WriteString(fmt.Sprintf("\n## 📁 %s\n", e.Path))
		} else if e.Kind == "file" {
			sb.WriteString(fmt.Sprintf("\n### 📄 %s\n", e.Path))
			sb.WriteString(fmt.Sprintf("@type: file\n@size: %d bytes\n", e.Size))

			lang := detectLanguage(e.Path)
			if lang != "" {
				sb.WriteString(fmt.Sprintf("@language: %s\n", lang))
			}

			if e.Content != "" {
				sb.WriteString("```" + lang + "\n")
				sb.WriteString(e.Content + "\n")
				sb.WriteString("```\n")
			}

			sb.WriteString("---\n")
		}
	}

	return sb.String()
}

// ------------------- JSON -------------------

func formatJSON(entries []FileEntry) (string, error) {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ------------------- Helper -------------------

func detectLanguage(path string) string {
	if dot := strings.LastIndex(path, "."); dot != -1 && dot < len(path)-1 {
		ext := path[dot+1:]
		switch ext {
		case "go":
			return "go"
		case "rs":
			return "rust"
		case "js":
			return "javascript"
		case "ts":
			return "typescript"
		case "py":
			return "python"
		case "java":
			return "java"
		case "c":
			return "c"
		case "cpp":
			return "cpp"
		case "html":
			return "html"
		case "css":
			return "css"
		case "md":
			return "markdown"
		default:
			return ""
		}
	}
	return ""
}
