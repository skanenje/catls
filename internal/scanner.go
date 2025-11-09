package internal

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// FileEntry represents one filesystem entity.
type FileEntry struct {
	Path      string
	Kind      string // "file" or "dir"
	Size      int64
	Depth     int
	Ignored   bool
	Error     string
}

// ScanDir walks the directory recursively and returns entries.
func ScanDir(root string, maxDepth int, ignore []string) ([]FileEntry, error) {
	var entries []FileEntry

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			entries = append(entries, FileEntry{Path: path, Error: err.Error()})
			return nil
		}

		depth := strings.Count(filepath.ToSlash(path), "/") - strings.Count(filepath.ToSlash(root), "/")
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
