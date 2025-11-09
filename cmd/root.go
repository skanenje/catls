package cmd

import (
	"flag"
	"fmt"
	"os"

	"catls/internal"
)

// Run executes the simplified CLI
func Run() {
	var maxDepth int
	var maxSize int64
	var outputFile string
	var showHelp bool

	flag.IntVar(&maxDepth, "max-depth", 2, "Maximum recursion depth (-1 for unlimited)")
	flag.Int64Var(&maxSize, "max-size", 512*1024, "Maximum file size to include (default 512KB)")
	flag.StringVar(&outputFile, "output", "", "Write output to file instead of stdout")
	flag.BoolVar(&showHelp, "help", false, "Show help information")
	flag.Parse()

	if showHelp {
		fmt.Println(`catls — Recursive file content dumper
Usage:
  catls [path] [flags]

Flags:
  --max-depth N    Limit recursion depth (-1 = unlimited, default 2)
  --max-size BYTES Skip files larger than given size (default 512KB)
  --output FILE    Write output to file instead of stdout
  --help           Show this message

Examples:
  catls .                  # Dump current directory recursively
  catls -max-depth=1 ./src
  catls -max-size=100000 -output=dump.txt ./`)
		return
	}

	// Default root is current directory
	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
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
}
