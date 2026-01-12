# ls2eza

A simple tool that translates `ls` command flags to their [eza](https://github.com/eza-community/eza) equivalents.

## Installation

### From releases

Download the latest binary from the [releases page](https://github.com/kluzzebass/ls2eza/releases).

### Using Go

```bash
go install github.com/kluzzebass/ls2eza@latest
```

### From source

```bash
git clone https://github.com/kluzzebass/ls2eza.git
cd ls2eza
make build
```

## Usage

ls2eza takes ls-style arguments and outputs the equivalent eza command:

```bash
$ ls2eza -la
eza -l -a

$ ls2eza -ltr
eza -l --sort=modified

$ ls2eza -lSh /tmp
eza -l --sort=size --reverse /tmp
```

### Using with an alias

You can create a shell alias to automatically translate and execute:

```bash
alias ls='eval $(ls2eza "$@")'
```

Or for fish shell:

```fish
function ls
    eval (ls2eza $argv)
end
```

## BSD vs GNU ls Compatibility

ls2eza supports both BSD ls (macOS, FreeBSD) and GNU ls (Linux) flag conventions. By default, it auto-detects based on your operating system:

- **macOS, FreeBSD, OpenBSD, NetBSD, DragonFly** → BSD mode
- **Linux, Windows, others** → GNU mode

You can override the detection with the `LS2EZA_MODE` environment variable:

```bash
# Force BSD mode
export LS2EZA_MODE=bsd

# Force GNU mode
export LS2EZA_MODE=gnu
```

### Conflicting Flags

Some flags have different meanings between BSD and GNU ls:

| Flag | BSD ls | GNU ls |
|------|--------|--------|
| `-T` | Full timestamp display | Tab size (ignored) |
| `-X` | Don't cross filesystems (ignored) | Sort by extension |
| `-I` | Prevent auto -A for superuser (ignored) | Ignore pattern (`-I PATTERN`) |
| `-w` | Raw non-printable chars (ignored) | Output width (`-w COLS`) |
| `-D` | Date format (`-D FORMAT`) | Dired mode (ignored) |

## Supported Flags

### Display Format

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-l` | `-l` | Long format |
| `-1` | `-1` | One entry per line |
| `-C` | `--grid` | Multi-column output |
| `-x` | `--across` | Sort grid across |
| `-m` | `--oneline` | Stream output |

### Show/Hide Entries

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-a` | `-a` | Show all (including . and ..) |
| `-A` | `-A` | Show hidden (except . and ..) |
| `-d` | `-d` | List directories themselves |
| `-R` | `--recurse` | Recurse into directories |

### Sorting

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-t` | `--sort=modified --reverse` | Sort by modification time |
| `-S` | `--sort=size --reverse` | Sort by size |
| `-c` | `--sort=changed --reverse` | Sort by change time |
| `-u` | `--sort=accessed --reverse` | Sort by access time |
| `-U` | `--sort=created --reverse` | Sort by creation time (BSD) |
| `-f` | `--sort=none -a` | Unsorted, show all |
| `-v` | `--sort=name` | Version/name sort |
| `-r` | `--reverse` | Reverse sort order |

### Long Format Options

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-i` | `--inode` | Show inode numbers |
| `-s` | `--blocksize` | Show allocated blocks |
| `-n` | `--numeric` | Numeric user/group IDs |
| `-o` | `-l --no-group` | Long format without group |
| `-g` | `-l --no-user` | Long format without owner |
| `-O` | `--flags` | Show file flags (BSD/macOS) |
| `-@` | `--extended` | Show extended attributes |
| `-h` | (default) | Human-readable sizes |

### Indicators

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-F` | `-F` | Append type indicators |
| `-p` | `--classify` | Append / to directories |

### Symlinks

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-L` | `-X` | Dereference symlinks |
| `-H` | `-X` | Follow symlinks on command line |

### BSD-Specific Flags

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-T` | `--time-style=full-iso` | Full timestamp display |
| `-D FORMAT` | `--time-style=+FORMAT` | Custom date format (strftime) |
| `-G` | (default) | Color output |
| `-I` | (ignored) | Prevent auto -A for superuser |
| `-W` | (ignored) | Display whiteouts |
| `-X` | (ignored) | Don't cross filesystems |
| `-w` | (ignored) | Raw non-printable characters |

### GNU-Specific Flags

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-X` | `--sort=extension` | Sort by file extension |
| `-I PATTERN` | `--ignore-glob=PATTERN` | Ignore files matching pattern |
| `-w COLS` | `--width=COLS` | Set output width |
| `-Z` | `-Z` | SELinux security context |
| `-N` | `--no-quotes` | Print names without quoting |
| `--group-directories-first` | `--group-directories-first` | List directories first |
| `--full-time` | `-l --time-style=full-iso` | Long format with full time |
| `--ignore=PATTERN` | `--ignore-glob=PATTERN` | Ignore files matching pattern |
| `--hyperlink` | `--hyperlink` | Hyperlink file names |

### Other

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-k` | (ignored) | Block size handling |
| `-V` | (ls2eza only) | Show version |

### Long Options

| ls option | eza equivalent |
|-----------|----------------|
| `--all` | `-a` |
| `--almost-all` | `-A` |
| `--directory` | `-d` |
| `--recursive` | `--recurse` |
| `--human-readable` | (default) |
| `--inode` | `--inode` |
| `--numeric-uid-gid` | `--numeric` |
| `--classify` | `-F` |
| `--dereference` | `-X` |
| `--color=WHEN` | `--color=WHEN` |

### Sort Order Handling

ls and eza have opposite default sort orders for time and size sorting. ls2eza automatically adds `--reverse` when needed so that the output matches ls behavior:

- `ls -lt` shows newest first → `eza --sort=modified --reverse`
- `ls -ltr` shows oldest first → `eza --sort=modified` (no reverse needed)
- `ls -lS` shows largest first → `eza --sort=size --reverse`
- `ls -lc` shows newest changed first → `eza --sort=changed --reverse`

## License

MIT
