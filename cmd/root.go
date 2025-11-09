package cmd

import (
	"fmt"
	"os"

	"catls/internal"
	"github.com/spf13/cobra"
)

var (
	maxDepth   int
	maxSize    int64
	outputFile string
)

var rootCmd = &cobra.Command{
	Use:   "catls [path]",
	Short: "Recursive file content dumper",
	Long:  `catls recursively dumps the contents of text files in a directory tree.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		root := "."
		if len(args) > 0 {
			root = args[0]
		}

		cfg := internal.DumpConfig{
			MaxDepth:    maxDepth,
			MaxFileSize: maxSize,
		}

		// Choose writer (stdout or file)
		out := os.Stdout
		if outputFile != "" {
			f, err := os.Create(outputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create output file: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()
			out = f
		}

		// Dump recursively
		if err := internal.DumpRecursive(root, cfg, 0, out); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if outputFile != "" {
				fmt.Printf("Output written to %s\n", outputFile)
			}
	},
}

func init() {
	rootCmd.Flags().IntVar(&maxDepth, "max-depth", 2, "Maximum recursion depth (-1 for unlimited)")
	rootCmd.Flags().Int64Var(&maxSize, "max-size", 512*1024, "Maximum file size to include (default 512KB)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write output to file instead of stdout")
}

// Run executes the CLI
func Run() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
