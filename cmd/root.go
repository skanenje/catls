// cmd/root.go
package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

// Config holds runtime options
type Config struct {
	Path       string
	MaxDepth   int
	MaxSize    int64
	OutputMode string
	Ignore     []string
	Summary    bool
	OutputFile string
}

var cfg Config

// rootCmd is the base command
var rootCmd = &cobra.Command{
	Use:   "catls [path]",
	Short: "catls merges cat + ls to serialize structure and content",
	Long: `catls recursively walks directories, reading both structure and file contents
to produce AI-friendly Markdown or JSON output.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg.Path = args[0]

		// For now, just print parsed config (placeholder for pipeline)
		fmt.Printf("Running catls on: %s\n", cfg.Path)
		fmt.Printf("Depth: %d, MaxSize: %d bytes, Format: %s\n",
			cfg.MaxDepth, cfg.MaxSize, cfg.OutputMode)
		fmt.Printf("Ignore: %v\n", cfg.Ignore)
		fmt.Printf("Summary: %v\n", cfg.Summary)
		fmt.Printf("Output File: %s\n", cfg.OutputFile)

		// In Stage 4.2, we’ll call the scanner and formatter here.
		return nil
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
}
