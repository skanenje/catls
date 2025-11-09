// cmd/root.go
package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"catls/internal"

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

		entries, err := internal.ScanDir(cfg.Path, cfg.MaxDepth, cfg.Ignore)
		if err != nil {
			return err
		}

		for _, e := range entries {
			fmt.Printf("[%s] %s (%d bytes, depth=%d)\n", e.Kind, e.Path, e.Size, e.Depth)
		}


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
