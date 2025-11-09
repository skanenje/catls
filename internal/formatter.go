package internal

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

// FormatMode defines output format
type FormatMode string

const (
	Markdown FormatMode = "markdown"
	JSONMode FormatMode = "json"
)

// StreamFormatEntries consumes entries from channel and writes formatted output
func StreamFormatEntries(entries <-chan FileEntry, mode FormatMode, writer func(string)) error {
	if mode == JSONMode {
		writer("[\n")
		first := true
		for e := range entries {
			if !first {
				writer(",\n")
			}
			first = false
			data, _ := json.MarshalIndent(e, "  ", "  ")
			writer(string(data))
		}
		writer("\n]\n")
	} else { // Markdown
		for e := range entries {
			if e.Kind == "dir" {
				writer(fmt.Sprintf("## 📁 %s\n\n", e.Path))
			} else {
				writer(fmt.Sprintf("### 📄 %s\n@type: %s\n@size: %d bytes\n\n```%s\n%s\n```\n\n",
					e.Path, e.Kind, e.Size, detectLanguage(e.Path), e.Content))
			}
		}
	}
	return nil
}

// Simple extension-based language detection
func detectLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go":
		return "go"
	case ".rs":
		return "rust"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".py":
		return "python"
	case ".java":
		return "java"
	default:
		return ""
	}
}
