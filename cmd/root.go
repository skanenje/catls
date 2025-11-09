package cmd

import (
	"fmt"
	"os"

	"catls/internal"

	"github.com/spf13/cobra"
)

// Config holds runtime options
type Config struct {
	Path        string
	MaxDepth    int
	MaxSize     int64
	OutputMode  string
	Ignore      []string
	Summary     bool
	OutputFile  string
	ShowContent bool
	Lines       int
}

var cfg Config

var rootCmd = &cobra.Command{
	Use:   "catls [path]",
	Short: "catls merges cat + ls to serialize structure and content",
	Long: `catls recursively walks directories, reading both structure and file contents
to produce AI-friendly Markdown or JSON output.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg.Path = args[0]

		// Determine if full content is needed (JSON always includes full content)
		fullContent := cfg.OutputMode == "json"

		// Buffered channel for streaming entries
		entriesChan := make(chan internal.FileEntry, 100)

		// Launch scanner in a goroutine
		go func() {
			err := internal.ScanDirStream(cfg.Path, cfg.MaxDepth, cfg.Ignore, cfg.ShowContent, cfg.Lines, fullContent, entriesChan)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error scanning directory: %v\n", err)
			}
		}()

		// Writer function for formatter
		writerFunc := func(s string) {
			if cfg.OutputFile != "" {
				f, err := os.OpenFile(cfg.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
					return
				}
				defer f.Close()
				f.WriteString(s)
			} else {
				fmt.Print(s)
			}
		}

		// Summary-only mode: print basic info without content
		if cfg.Summary {
			for entry := range entriesChan {
				fmt.Printf("[%s] %s (%d bytes, depth=%d)\n", entry.Kind, entry.Path, entry.Size, entry.Depth)
			}
			return nil
		}

		// Stream entries to formatter
		return internal.StreamFormatEntries(entriesChan, internal.FormatMode(cfg.OutputMode), writerFunc)
	},
}

// Execute bootstraps CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVar(&cfg.MaxDepth, "max-depth", -1, "Limit recursion depth")
	rootCmd.Flags().Int64Var(&cfg.MaxSize, "max-size", 64000, "Max bytes per file")
	rootCmd.Flags().StringVar(&cfg.OutputMode, "format", "markdown", "Output format: markdown or json")
	rootCmd.Flags().StringSliceVar(&cfg.Ignore, "ignore", []string{".git", "node_modules"}, "Ignore patterns")
	rootCmd.Flags().BoolVar(&cfg.Summary, "summary", false, "Structure only, no content")
	rootCmd.Flags().StringVar(&cfg.OutputFile, "output", "", "Write output to file instead of stdout")
	rootCmd.Flags().BoolVar(&cfg.ShowContent, "show-content", false, "Show file content preview in Markdown")
	rootCmd.Flags().IntVar(&cfg.Lines, "lines", 5, "Number of lines to preview if --show-content is true")
}
