package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// Version information - set via ldflags at build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// LSMode determines which ls flavor to emulate
type LSMode int

const (
	ModeBSD LSMode = iota
	ModeGNU
)

// getLSMode returns the ls compatibility mode based on OS or environment override
// Override with LS2EZA_MODE=bsd or LS2EZA_MODE=gnu
func getLSMode() LSMode {
	// Check environment override first
	if mode := os.Getenv("LS2EZA_MODE"); mode != "" {
		switch strings.ToLower(mode) {
		case "bsd":
			return ModeBSD
		case "gnu":
			return ModeGNU
		}
	}

	// Default based on OS
	switch runtime.GOOS {
	case "darwin", "freebsd", "openbsd", "netbsd", "dragonfly":
		return ModeBSD
	default:
		// linux, windows, and others default to GNU
		return ModeGNU
	}
}

// Flags that need --reverse in eza to match ls default behavior
// (ls shows newest/largest first, eza shows oldest/smallest first)
var reverseNeeded = map[rune]bool{
	't': true, // time sort: ls=newest first, eza=oldest first
	'S': true, // size sort: ls=largest first, eza=smallest first
	'c': true, // change time sort: ls=newest first, eza=oldest first
	'u': true, // access time sort: ls=newest first, eza=oldest first
	'U': true, // creation time sort (BSD): ls=newest first, eza=oldest first
}

// Simple 1:1 flag mappings
var flagMap = map[rune][]string{
	// Display format
	'l': {"-l"},                // long format
	'1': {"-1"},                // one entry per line
	'C': {"--grid"},            // multi-column output (default in terminal)
	'x': {"--across"},          // sort grid across rather than down
	'm': {"--oneline"},         // stream output (eza doesn't have comma-separated, use oneline)

	// Show/hide entries
	'a': {"-a"},                // show all including . and ..
	'A': {"-A"},                // show hidden but not . and ..
	'd': {"-d"},                // list directories themselves, not contents
	'R': {"--recurse"},         // recurse into directories

	// Sorting
	't': {"--sort=modified"},   // sort by modification time
	'S': {"--sort=size"},       // sort by size
	'c': {"--sort=changed"},    // sort by change time
	'u': {"--sort=accessed"},   // sort by access time
	'U': {"--sort=created"},    // sort by creation time (BSD)
	'f': {"--sort=none", "-a"}, // unsorted, show all
	'v': {"--sort=name"},       // natural version sort (approximate)

	// File size display
	'h': {},                    // human-readable (default in eza)
	'k': {},                    // 1024-byte blocks (eza handles differently)
	's': {"--blocksize"},       // show allocated blocks

	// Indicators and classification
	'F': {"-F"},                // append file type indicators (*/=>@|)
	'p': {"--classify"},        // append / to directories

	// Long format options
	'i': {"--inode"},           // show inode numbers
	'n': {"--numeric"},         // numeric user/group IDs
	'o': {"-l", "--no-group"},  // long format without group (BSD)
	'g': {"-l", "--no-user"},   // long format without owner (GNU style)
	'O': {"--flags"},           // show file flags (BSD/macOS)
	'e': {},                    // show ACL (no eza equivalent)
	// 'T' handled specially - BSD=full time, GNU=tabsize
	'@': {"--extended"},        // show extended attributes

	// Symlink handling
	'L': {"-X"},                // dereference symlinks
	'H': {"-X"},                // follow symlinks on command line
	'P': {},                    // don't follow symlinks (default)

	// Color
	'G': {},                    // color output (default in eza)

	// Misc BSD
	'q': {},                    // replace non-printable with ? (no eza equivalent)
	// 'w' handled specially - GNU uses -w COLS for width
	'b': {},                    // C-style escapes (no eza equivalent)
	'B': {},                    // octal escapes (BSD) / ignore-backups (GNU) - ignore both
	'W': {},                    // display whiteouts (BSD, no eza equivalent)

	// GNU ls specific (non-conflicting)
	// 'X' handled specially - BSD=ignore, GNU=sort by extension
	'Z': {"-Z"},                // SELinux security context
	'N': {"--no-quotes"},       // print names without quoting
	'Q': {},                    // quote names (eza doesn't quote by default anyway)
}

// Long option mappings (ls long options to eza equivalents)
var longFlagMap = map[string][]string{
	// Common options
	"--all":             {"-a"},
	"--almost-all":      {"-A"},
	"--directory":       {"-d"},
	"--recursive":       {"--recurse"},
	"--human-readable":  {}, // default in eza
	"--inode":           {"--inode"},
	"--numeric-uid-gid": {"--numeric"},
	"--classify":        {"-F"},
	"--file-type":       {"--classify"},
	"--dereference":     {"-X"},
	"--no-group":        {"--no-group"},

	// GNU ls specific
	"--group-directories-first": {"--group-directories-first"},
	"--reverse":                 {"--reverse"},
	"--size":                    {"--blocksize"},
	"--context":                 {"-Z"},
	"--literal":                 {"--no-quotes"},
	"--quote-name":              {}, // eza doesn't quote by default
	"--hide-control-chars":      {}, // no eza equivalent
	"--show-control-chars":      {}, // no eza equivalent
	"--hyperlink":               {"--hyperlink"},
	"--full-time":               {"-l", "--time-style=full-iso"},
	"--author":                  {}, // no eza equivalent
	"--escape":                  {}, // no eza equivalent
	"--ignore-backups":          {}, // no eza equivalent
	"--kibibytes":               {}, // eza uses binary by default
	"--si":                      {}, // no eza equivalent (powers of 1000)
	"--dired":                   {}, // no eza equivalent (emacs mode)
	"--zero":                    {}, // no eza equivalent (NUL terminated)
}

// Long options with =value that need prefix matching
var longFlagPrefixes = []struct {
	prefix string
	pass   bool // true = pass through to eza, false = ignore
}{
	{"--color", true},
	{"--colour", true},
	{"--sort=", true},
	{"--time=", true},
	{"--time-style=", true},
	{"--hyperlink=", true},
	{"--width=", true},
	{"--ignore=", true},      // GNU -I/--ignore=PATTERN → eza -I/--ignore-glob
	{"--hide=", false},       // no direct equivalent
	{"--block-size=", false}, // no eza equivalent
	{"--indicator-style=", false},
	{"--quoting-style=", false},
	{"--tabsize=", false},
}

func translateFlags(args []string, mode LSMode) []string {
	var ezaArgs []string
	var paths []string
	userReverse := false
	needsReverse := false
	skipNext := false

	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}

		if strings.HasPrefix(arg, "--") {
			// Long option handling
			if arg == "--reverse" {
				userReverse = true
				continue
			}

			// Check for prefix matches (options with =value)
			handled := false
			for _, pf := range longFlagPrefixes {
				if strings.HasPrefix(arg, pf.prefix) {
					if pf.pass {
						// Special case: --ignore= → --ignore-glob
						if strings.HasPrefix(arg, "--ignore=") {
							pattern := strings.TrimPrefix(arg, "--ignore=")
							ezaArgs = append(ezaArgs, "--ignore-glob="+pattern)
						} else {
							ezaArgs = append(ezaArgs, arg)
						}
					}
					// If !pf.pass, we just ignore it
					handled = true
					break
				}
			}
			if handled {
				continue
			}

			// Check exact long option matches
			if mapped, ok := longFlagMap[arg]; ok {
				ezaArgs = append(ezaArgs, mapped...)
			} else {
				// Unknown long option - pass through
				ezaArgs = append(ezaArgs, arg)
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			// Short options - translate each character
			flags := arg[1:]
			for j, c := range flags {
				if c == 'r' {
					userReverse = true
					continue
				}
				if c == 'D' {
					// BSD: -D FORMAT (date format) → --time-style=+FORMAT
					// GNU: -D (dired mode) → ignore
					if mode == ModeBSD {
						remaining := flags[j+1:]
						var format string
						if len(remaining) > 0 {
							format = string(remaining)
						} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
							format = args[i+1]
							skipNext = true
						}
						if format != "" {
							ezaArgs = append(ezaArgs, "--time-style=+"+format)
						}
						break // -D consumes rest of this arg
					}
					// GNU dired mode - ignore
					continue
				}
				if c == 'I' {
					// BSD: -I (prevent auto -A for superuser) → ignore (no argument)
					// GNU: -I PATTERN (ignore glob) → --ignore-glob=PATTERN
					remaining := flags[j+1:]
					if mode == ModeGNU {
						var pattern string
						if len(remaining) > 0 {
							pattern = string(remaining)
						} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
							pattern = args[i+1]
							skipNext = true
						}
						if pattern != "" {
							ezaArgs = append(ezaArgs, "--ignore-glob="+pattern)
						}
						break // -I consumes rest of this arg in GNU mode
					}
					// BSD -I has no argument, continue processing other flags
					continue
				}
				if c == 'w' {
					// BSD: -w (raw non-printable) → ignore
					// GNU: -w COLS (width) → --width=COLS
					if mode == ModeGNU {
						remaining := flags[j+1:]
						var width string
						if len(remaining) > 0 {
							width = string(remaining)
						} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
							width = args[i+1]
							skipNext = true
						}
						if width != "" {
							ezaArgs = append(ezaArgs, "--width="+width)
						}
						break // -w consumes rest of this arg
					}
					// BSD -w - ignore
					continue
				}
				if c == 'T' {
					// BSD: -T (full time) → --time-style=full-iso
					// GNU: -T COLS (tab size) → ignore
					if mode == ModeBSD {
						ezaArgs = append(ezaArgs, "--time-style=full-iso")
						continue
					}
					// GNU -T COLS - consume argument and ignore
					remaining := flags[j+1:]
					if len(remaining) > 0 {
						break // consumed as part of -T
					} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						skipNext = true
					}
					continue
				}
				if c == 'X' {
					// BSD: -X (don't cross filesystems) → ignore
					// GNU: -X (sort by extension) → --sort=extension
					if mode == ModeGNU {
						ezaArgs = append(ezaArgs, "--sort=extension")
					}
					// BSD -X - ignore
					continue
				}
				if reverseNeeded[c] {
					needsReverse = true
				}
				if mapped, ok := flagMap[c]; ok {
					ezaArgs = append(ezaArgs, mapped...)
				} else {
					// Unknown flag - try passing it through
					ezaArgs = append(ezaArgs, "-"+string(c))
				}
			}
		} else {
			// Not a flag - it's a path
			paths = append(paths, arg)
		}
	}

	// XOR logic: reverse if exactly one of (needsReverse, userReverse) is true
	// - ls -lt → need reverse to get newest first (needsReverse=true, userReverse=false) → add --reverse
	// - ls -ltr → user wants oldest first (needsReverse=true, userReverse=true) → don't add --reverse
	// - ls -lr → user wants reverse alpha (needsReverse=false, userReverse=true) → add --reverse
	if needsReverse != userReverse {
		ezaArgs = append(ezaArgs, "--reverse")
	}

	// Deduplicate flags
	seen := make(map[string]bool)
	var deduped []string
	for _, f := range ezaArgs {
		if !seen[f] {
			seen[f] = true
			deduped = append(deduped, f)
		}
	}

	return append(deduped, paths...)
}

func shellQuote(s string) string {
	if strings.ContainsAny(s, " \t\n\"'\\$`!") {
		return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
	}
	return s
}

func printVersion() {
	fmt.Printf("ls2eza %s\n", version)
	if commit != "none" {
		fmt.Printf("  commit: %s\n", commit)
	}
	if date != "unknown" {
		fmt.Printf("  built:  %s\n", date)
	}
}

func main() {
	args := os.Args[1:]

	// Handle version flag
	for _, arg := range args {
		if arg == "-V" || arg == "--version" {
			printVersion()
			return
		}
	}

	mode := getLSMode()
	ezaArgs := translateFlags(args, mode)

	// Build and print the command
	parts := make([]string, len(ezaArgs)+1)
	parts[0] = "eza"
	for i, arg := range ezaArgs {
		parts[i+1] = shellQuote(arg)
	}
	fmt.Println(strings.Join(parts, " "))
}
