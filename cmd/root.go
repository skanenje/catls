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
	var showHelp bool

	flag.IntVar(&maxDepth, "max-depth", 2, "Maximum recursion depth (-1 for unlimited)")
	flag.Int64Var(&maxSize, "max-size", 512*1024, "Maximum file size to include (default 512KB)")
	flag.BoolVar(&showHelp, "help", false, "Show help information")
	flag.Parse()

	if showHelp {
		fmt.Println(`catls — Recursive file content dumper
Usage:
  catls [path] [flags]

Flags:
  --max-depth N    Limit recursion depth (-1 = unlimited, default 2)
  --max-size BYTES Skip files larger than given size (default 512KB)
  --help           Show this message

Examples:
  catls .                  # Dump current directory recursively
  catls ./src --max-depth=1
  catls ./ --max-size=100000`)
		return
	}

	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	cfg := internal.DumpConfig{
		MaxDepth:    maxDepth,
		MaxFileSize: maxSize,
	}

	if err := internal.DumpRecursive(root, cfg, 0); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
