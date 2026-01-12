# ls2eza

A simple tool that translates `ls` command flags to their [eza](https://github.com/eza-community/eza) equivalents.

## Installation

```bash
go install github.com/kluzzebass/ls2eza@latest
```

Or build from source:

```bash
git clone https://github.com/kluzzebass/ls2eza.git
cd ls2eza
go build -o ls2eza
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

## Supported Flags

| ls flag | eza equivalent |
|---------|----------------|
| `-l` | `-l` |
| `-a` | `-a` |
| `-A` | `-A` |
| `-h` | (default in eza) |
| `-t` | `--sort=modified --reverse` |
| `-S` | `--sort=size --reverse` |
| `-r` | `--reverse` |
| `-R` | `--recurse` |
| `-1` | `-1` |
| `-d` | `-d` |
| `-F` | `-F` |
| `-G` | (default in eza) |
| `-i` | `--inode` |
| `-s` | `--blocksize` |
| `-n` | `--numeric` |
| `-o` | `-l --no-group` |
| `-g` | `-l --no-user` |
| `-p` | `--classify` |
| `-c` | `--sort=changed` |
| `-u` | `--sort=accessed` |
| `-x` | `--across` |
| `-C` | `--grid` |
| `-T` | `--tree` |

### Sort order handling

ls and eza have opposite default sort orders for time (`-t`) and size (`-S`) sorting. ls2eza automatically adds `--reverse` when needed so that the output matches ls behavior:

- `ls -lt` shows newest first → `eza --sort=modified --reverse`
- `ls -ltr` shows oldest first → `eza --sort=modified` (no reverse needed)

## License

MIT
