package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kluzzebass/reflag/translator"
	_ "github.com/kluzzebass/reflag/translator/ls2eza" // Register ls2eza translator
)

// Version information - set via ldflags at build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func shellQuote(s string) string {
	if strings.ContainsAny(s, " \t\n\"'\\$`!") {
		return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
	}
	return s
}

func printVersion(name string) {
	fmt.Printf("%s %s\n", name, version)
	if commit != "none" {
		fmt.Printf("  commit: %s\n", commit)
	}
	if date != "unknown" {
		fmt.Printf("  built:  %s\n", date)
	}
}

func printUsage() {
	fmt.Println("reflag - translate command-line flags between tools")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  reflag <source> <target> [flags...]")
	fmt.Println("  reflag --list")
	fmt.Println("  reflag --version")
	fmt.Println()
	fmt.Println("Symlink mode:")
	fmt.Println("  Create a symlink named <source>2<target> pointing to reflag")
	fmt.Println("  Example: ln -s reflag ls2eza")
	fmt.Println("           ls2eza -la  # outputs: eza -l -a")
	fmt.Println()
	fmt.Println("Available translators:")
	names := translator.List()
	sort.Strings(names)
	for _, name := range names {
		t := translator.GetByName(name)
		fmt.Printf("  %s: %s -> %s\n", name, t.SourceTool(), t.TargetTool())
	}
}

// detectFromBinaryName parses a binary name like "ls2eza" into source and target
func detectFromBinaryName(name string) (source, target string, ok bool) {
	parts := strings.SplitN(name, "2", 2)
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], true
	}
	return "", "", false
}

func runTranslator(t translator.Translator, args []string) {
	// Handle version flag
	for _, arg := range args {
		if arg == "-V" || arg == "--version" {
			printVersion(t.Name())
			return
		}
	}

	ezaArgs := t.Translate(args)

	// Build and print the command
	parts := make([]string, len(ezaArgs)+1)
	parts[0] = t.TargetTool()
	for i, arg := range ezaArgs {
		parts[i+1] = shellQuote(arg)
	}
	fmt.Println(strings.Join(parts, " "))
}

func main() {
	binary := filepath.Base(os.Args[0])

	// Check if running as a symlink (e.g., ls2eza -> reflag)
	if source, target, ok := detectFromBinaryName(binary); ok {
		t := translator.Get(source, target)
		if t == nil {
			fmt.Fprintf(os.Stderr, "error: no translator registered for %s to %s\n", source, target)
			os.Exit(1)
		}
		runTranslator(t, os.Args[1:])
		return
	}

	// Running as reflag directly
	args := os.Args[1:]

	// Handle reflag's own flags
	if len(args) == 0 {
		printUsage()
		return
	}

	switch args[0] {
	case "--version", "-V":
		printVersion("reflag")
		return
	case "--list", "-l":
		names := translator.List()
		sort.Strings(names)
		for _, name := range names {
			t := translator.GetByName(name)
			fmt.Printf("%s: %s -> %s\n", name, t.SourceTool(), t.TargetTool())
		}
		return
	case "--help", "-h":
		printUsage()
		return
	}

	// Explicit mode: reflag <source> <target> [flags...]
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "error: expected <source> <target> arguments")
		fmt.Fprintln(os.Stderr, "usage: reflag <source> <target> [flags...]")
		os.Exit(1)
	}

	source, target := args[0], args[1]
	t := translator.Get(source, target)
	if t == nil {
		fmt.Fprintf(os.Stderr, "error: no translator registered for %s to %s\n", source, target)
		fmt.Fprintln(os.Stderr, "use 'reflag --list' to see available translators")
		os.Exit(1)
	}

	runTranslator(t, args[2:])
}
