package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DumpConfig defines limits for recursive file dumping
type DumpConfig struct {
	MaxDepth    int
	MaxFileSize int64 // bytes
}

// Heuristic to detect text-like files
func IsTextFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	textExts := []string{
		".go", ".js", ".ts", ".json", ".py", ".rs", ".c", ".cpp",
		".h", ".html", ".css", ".md", ".txt", ".sh", ".yaml", ".yml",
	}
	for _, e := range textExts {
		if ext == e {
			return true
		}
	}
	return false
}

// DumpRecursive walks a directory tree and prints readable file contents.
func DumpRecursive(root string, cfg DumpConfig, depth int, out io.Writer) error {
	if cfg.MaxDepth >= 0 && depth > cfg.MaxDepth {
		return nil
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return fmt.Errorf("failed to read dir %s: %w", root, err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(root, entry.Name())

		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		if entry.IsDir() {
			_ = DumpRecursive(fullPath, cfg, depth+1, out)
			continue
		}

		if !IsTextFile(fullPath) || info.Size() > cfg.MaxFileSize {
			continue
		}

		data, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		fmt.Fprintf(out, "\n### %s\n", fullPath)
		fmt.Fprintf(out, "Size: %d bytes\n", info.Size())
		fmt.Fprintln(out, "---")
		fmt.Fprintln(out, string(data))
		fmt.Fprintln(out, "\n--- end of file ---")
	}

	return nil
}
