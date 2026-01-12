# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ls2eza is a Go CLI tool that translates `ls` command-line flags to their `eza` equivalents. It outputs the translated eza command (does not execute it).

## Build and Test Commands

```bash
make build               # Build with version info
make test                # Run all tests
make clean               # Remove binary
make build-all           # Cross-compile for all platforms
go run main.go -la       # Run directly without building
go test -run TestName    # Run specific test
```

## Releasing

Push a semver tag to trigger a GitHub Actions release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Architecture

Single-file Go application (`main.go`) with these components:

1. **Mode detection** - `getLSMode()` determines BSD vs GNU ls compatibility:
   - Auto-detects based on OS (darwin/freebsd → BSD, linux/others → GNU)
   - Override with `LS2EZA_MODE=bsd` or `LS2EZA_MODE=gnu`
2. **Flag mappings** - Static maps defining ls→eza translations:
   - `reverseNeeded` - Flags that need sort order correction (`t`, `S`, `c`, `u`, `U`)
   - `flagMap` - Short flag translations (30+ flags)
   - `longFlagMap` - Long option translations (`--all`, `--recursive`, etc.)
   - `longFlagPrefixes` - Long options with =value that need prefix matching
3. **`translateFlags(mode)`** - Core logic that parses arguments, applies mode-specific mappings, handles reverse-sort semantics, and deduplicates flags
4. **`main()`** - Entry point that outputs the shell-quoted eza command

### BSD vs GNU conflicts

These flags have different meanings between BSD and GNU ls:
- `-T`: BSD=full time display, GNU=tab size (ignored)
- `-X`: BSD=don't cross filesystems (ignored), GNU=sort by extension
- `-I`: BSD=prevent auto -A (ignored), GNU=ignore pattern
- `-w`: BSD=raw non-printable (ignored), GNU=output width
- `-D`: BSD=date format, GNU=dired mode (ignored)

### Key behavior: Reverse sort handling

ls and eza have opposite default sort orders for time (`-t`, `-c`, `-u`, `-U`) and size (`-S`). The tool uses XOR logic to determine when to add `--reverse`:
- `ls -lt` needs `--reverse` (ls shows newest first, eza shows oldest first)
- `ls -ltr` does NOT need `--reverse` (user explicitly wants oldest first)
- Same logic applies to `-c` (changed), `-u` (accessed), `-U` (created), `-S` (size)
