// cmd/root.go
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

// rootCmd is the base command
var rootCmd = &cobra.Command{
	Use:   "catls [path]",
	Short: "catls merges cat + ls to serialize structure and content",
	Long: `catls recursively walks directories, reading both structure and file contents
to produce AI-friendly Markdown or JSON output.

Examples:
  catls ./ --max-depth=2
  catls ./ --format=json --output=project.json
  catls ./ --ignore=.git,node_modules --show-content --lines=5

Flags:
  --max-depth <n>     Limit recursion depth
  --max-size <bytes>  Max bytes per file to read (default 64000)
  --format <mode>     Output format: markdown | json
  --ignore <list>     Comma-separated ignore patterns
  --summary           Only show structure (no file contents)
  --show-content      Include file contents (Markdown only)
  --lines <n>         Number of preview lines with --show-content
  --output <path>     Write output to file instead of stdout
  -h, --help          Show this help message`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// --- Assign runtime arguments
		cfg.Path = args[0]
		fullContent := cfg.OutputMode == "json"

		// --- Display configuration summary for verification
		fmt.Printf("\n▶ Running catls on: %s\n", cfg.Path)
		fmt.Printf("• Depth: %d | MaxSize: %d bytes\n", cfg.MaxDepth, cfg.MaxSize)
		fmt.Printf("• Format: %s | Summary: %v | ShowContent: %v | Lines: %d\n",
			cfg.OutputMode, cfg.Summary, cfg.ShowContent, cfg.Lines)
		fmt.Printf("• Ignore: %v\n", cfg.Ignore)
		if cfg.OutputFile != "" {
			fmt.Printf("• Output File: %s\n\n", cfg.OutputFile)
		} else {
			fmt.Println("• Output: stdout")
		}

		// --- Choose destination
		out := os.Stdout
		if cfg.OutputFile != "" {
			f, err := os.Create(cfg.OutputFile)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer f.Close()
			out = f
		}

		// --- Stream output
		return internal.StreamFormatEntries(
			out,
			cfg.Path,
			internal.FormatMode(cfg.OutputMode),
			cfg.MaxDepth,
			cfg.Ignore,
			cfg.ShowContent,
			cfg.Lines,
			fullContent,
		)
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
	// Register flags
	rootCmd.Flags().IntVar(&cfg.MaxDepth, "max-depth", -1, "Limit recursion depth")
	rootCmd.Flags().Int64Var(&cfg.MaxSize, "max-size", 64000, "Max bytes per file to read")
	rootCmd.Flags().StringVar(&cfg.OutputMode, "format", "markdown", "Output format: markdown or json")
	rootCmd.Flags().StringSliceVar(&cfg.Ignore, "ignore", []string{".git", "node_modules"}, "Ignore patterns (comma separated)")
	rootCmd.Flags().BoolVar(&cfg.Summary, "summary", false, "Structure only, no content")
	rootCmd.Flags().StringVar(&cfg.OutputFile, "output", "", "Write output to file instead of stdout")
	rootCmd.Flags().BoolVar(&cfg.ShowContent, "show-content", false, "Include file contents (Markdown mode only)")
	rootCmd.Flags().IntVar(&cfg.Lines, "lines", 10, "Number of preview lines to read when using --show-content")

	// Cobra automatically wires --help and -h
}
