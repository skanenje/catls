package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FormatMode defines output type.
type FormatMode string

const (
	FormatMarkdown FormatMode = "markdown"
	FormatJSON     FormatMode = "json"
)

// StreamFormatEntries writes directory structure incrementally.
func StreamFormatEntries(
	out io.Writer,
	root string,
	mode FormatMode,
	maxDepth int,
	ignore []string,
	showContent bool,
	lines int,
	fullContent bool,
) error {

	first := true // For JSON commas

	if mode == FormatJSON {
		_, _ = fmt.Fprint(out, "[")
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			entry := FileEntry{Path: path, Error: err.Error()}
			return writeEntry(out, entry, mode, &first)
		}

		depth := strings.Count(filepath.ToSlash(path), "/") - strings.Count(filepath.ToSlash(root), "/")
		if maxDepth >= 0 && depth > maxDepth {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		for _, ig := range ignore {
			if strings.Contains(path, ig) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		info, err := d.Info()
		if err != nil {
			entry := FileEntry{Path: path, Error: err.Error()}
			return writeEntry(out, entry, mode, &first)
		}

		entry := FileEntry{
			Path:    path,
			Kind:    kindFromDirEntry(d),
			Size:    info.Size(),
			Depth:   depth,
			Ignored: false,
		}

		// Read file content only if requested
		if entry.Kind == "file" && (showContent || fullContent) && info.Size() > 0 {
			content, err := os.ReadFile(path)
			if err == nil {
				entry.Content = string(content)
			} else {
				entry.Error = err.Error()
			}
		}

		return writeEntry(out, entry, mode, &first)
	})

	if mode == FormatJSON {
		_, _ = fmt.Fprint(out, "]")
	}

	return err
}

func writeEntry(out io.Writer, e FileEntry, mode FormatMode, first *bool) error {
	switch mode {
	case FormatMarkdown:
		fmt.Fprintf(out, "\n### %s\n@type: %s\n@size: %d bytes\n@depth: %d\n",
			e.Path, e.Kind, e.Size, e.Depth)
		if e.Error != "" {
			fmt.Fprintf(out, "⚠️ Error: %s\n", e.Error)
		}
		if e.Content != "" {
			fmt.Fprintf(out, "```text\n%s\n```\n", e.Content)
		}
		fmt.Fprintln(out, "---")
	case FormatJSON:
		b, err := json.Marshal(e)
		if err != nil {
			return err
		}
		if !*first {
			fmt.Fprint(out, ",")
		} else {
			*first = false
		}
		fmt.Fprint(out, string(b))
	}
	return nil
}
