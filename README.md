# reflag

A tool that translates command-line flags between different CLI tools. Currently supports translating `ls` flags to their [eza](https://github.com/eza-community/eza) equivalents.

## Installation

### From releases

Download the latest binary from the [releases page](https://github.com/kluzzebass/reflag/releases).

### Using Go

```bash
go install github.com/kluzzebass/reflag@latest
```

### From source

```bash
git clone https://github.com/kluzzebass/reflag.git
cd reflag
make build
```

## Usage

### Explicit Mode

Specify the source and target tools explicitly:

```bash
$ reflag ls eza -la
eza -l -a

$ reflag ls eza -ltr
eza -l --sort=modified

$ reflag ls eza -lSh /tmp
eza -l --sort=size --reverse /tmp
```

### Symlink Mode

Create a symlink named `<source>2<target>` pointing to reflag:

```bash
ln -s $(which reflag) ~/bin/ls2eza
ls2eza -la  # outputs: eza -l -a
```

### Shell Integration

You can create a shell function to automatically translate and execute:

```bash
# bash/zsh
ls() { eval $(reflag ls eza "$@"); }

# Or using symlink mode
ln -s $(which reflag) ~/bin/ls2eza
ls() { eval $(ls2eza "$@"); }
```

For fish shell:

```fish
function ls
    eval (reflag ls eza $argv)
end
```

### List Available Translators

```bash
$ reflag --list
ls2eza: ls -> eza
```

## ls2eza Translator

The ls2eza translator converts `ls` flags to `eza` equivalents.

### BSD vs GNU ls Compatibility

reflag supports both BSD ls (macOS, FreeBSD) and GNU ls (Linux) flag conventions. By default, it auto-detects based on your operating system:

- **macOS, FreeBSD, OpenBSD, NetBSD, DragonFly** -> BSD mode
- **Linux, Windows, others** -> GNU mode

Override with the `LS2EZA_MODE` environment variable:

```bash
export LS2EZA_MODE=bsd   # Force BSD mode
export LS2EZA_MODE=gnu   # Force GNU mode
```

### Supported Flags

#### Display Format

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-l` | `-l` | Long format |
| `-1` | `-1` | One entry per line |
| `-C` | `--grid` | Multi-column output |
| `-x` | `--across` | Sort grid across |
| `-m` | `--oneline` | Stream output |

#### Show/Hide Entries

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-a` | `-a` | Show all (including . and ..) |
| `-A` | `-A` | Show hidden (except . and ..) |
| `-d` | `-d` | List directories themselves |
| `-R` | `--recurse` | Recurse into directories |

#### Sorting

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

#### Long Format Options

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

### Sort Order Handling

ls and eza have opposite default sort orders for time and size sorting. reflag automatically adds `--reverse` when needed so that the output matches ls behavior:

- `ls -lt` shows newest first -> `eza --sort=modified --reverse`
- `ls -ltr` shows oldest first -> `eza --sort=modified` (no reverse needed)
- `ls -lS` shows largest first -> `eza --sort=size --reverse`

### Conflicting Flags (BSD vs GNU)

| Flag | BSD ls | GNU ls |
|------|--------|--------|
| `-T` | Full timestamp display | Tab size (ignored) |
| `-X` | Don't cross filesystems (ignored) | Sort by extension |
| `-I` | Prevent auto -A (ignored) | Ignore pattern (`-I PATTERN`) |
| `-w` | Raw non-printable chars (ignored) | Output width (`-w COLS`) |
| `-D` | Date format (`-D FORMAT`) | Dired mode (ignored) |

## Adding New Translators

reflag is designed to be extensible. To add a new translator:

1. Create a new package under `translator/`
2. Implement the `translator.Translator` interface
3. Register it in `init()` using `translator.Register()`

See `translator/ls2eza/` for an example implementation.

## License

MIT
