# reflag — command runner recipes

_default:
    @just --list

# Build the binary.
build:
    go build -o reflag .

# Run tests.
test:
    go test -v ./...

# Clean build artifacts.
clean:
    rm -f reflag
    rm -rf dist/

# Bump version, tag, and push (triggers GitHub Actions release). Usage: just release major|minor|patch
release bump:
    #!/usr/bin/env bash
    set -euo pipefail
    next=$(svu {{ bump }})
    echo "Releasing ${next}"
    git tag -a "${next}" -m "Release ${next}"
    git push origin "${next}"
