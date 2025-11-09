package internal

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FileEntry represents one filesystem entity.
type FileEntry struct {
	Path    string
	Kind    string // "file" or "dir"
	Size    int64
	Depth   int
	Content string
	Ignored bool
	Error   string
}

// ScanDir walks the directory recursively and returns entries.
func ScanDir(root string, maxDepth int, ignore []string, showContent bool, lines int) ([]FileEntry, error) {
	var entries []FileEntry

	rootAbs, _ := filepath.Abs(root)

	err := filepath.WalkDir(rootAbs, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			entries = append(entries, FileEntry{Path: path, Error: err.Error()})
			return nil
		}

		// --- Depth calculation ---
		rel, err := filepath.Rel(rootAbs, path)
		if err != nil {
			rel = path
		}
		depth := 0
		if rel != "." {
			depth = strings.Count(filepath.ToSlash(rel), "/") + 1
		}

		if maxDepth >= 0 && depth > maxDepth {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Ignore patterns
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
			entries = append(entries, FileEntry{Path: path, Error: err.Error()})
			return nil
		}

		entry := FileEntry{
			Path:  path,
			Kind:  kindFromDirEntry(d),
			Size:  info.Size(),
			Depth: depth,
		}

		// --- Read file content if requested ---
		if showContent && !d.IsDir() && info.Size() > 0 {
			if content, ok := readFilePreview(path, lines); ok {
				entry.Content = content
			}
		}

		entries = append(entries, entry)
		return nil
	})

	return entries, err
}

func kindFromDirEntry(d fs.DirEntry) string {
	if d.IsDir() {
		return "dir"
	}
	if d.Type()&fs.ModeSymlink != 0 {
		return "symlink"
	}
	return "file"
}

// readFilePreview reads first n lines of a text file, skips binary files
func readFilePreview(path string, n int) (string, bool) {
	f, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer f.Close()

	// Simple binary check: skip files with null bytes
	buf := make([]byte, 8000)
	count, _ := f.Read(buf)
	if strings.ContainsRune(string(buf[:count]), '\x00') {
		return "", false
	}

	f.Seek(0, 0) // rewind

	scanner := bufio.NewScanner(f)
	var lines []string
	for i := 0; i < n && scanner.Scan(); i++ {
		lines = append(lines, scanner.Text())
	}

	return strings.Join(lines, "\n"), true
}
