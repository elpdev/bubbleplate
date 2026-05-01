# Bubbleplate

Bubbleplate is an opinionated generator for building Go TUI projects with the [Charm](https://charm.sh) stack: Bubble Tea, Lip Gloss, and Bubbles.

## Install

Homebrew:

```sh
brew install elpdev/tap/bubbleplate
```

Arch Linux via AUR with yay:

```sh
yay -S bubbleplate-bin
```

Or install from the AUR manually:

```sh
git clone https://aur.archlinux.org/bubbleplate-bin.git
cd bubbleplate-bin
makepkg -si
```

## Generate a Project

Launch the interactive generator:

```sh
bubbleplate
```

Or generate non-interactively:

```sh
bubbleplate new myapp --module github.com/acme/myapp
cd myapp
go mod tidy
go test ./...
go run ./cmd/myapp
```

Useful options:

```sh
bubbleplate new myapp \
  --module github.com/acme/myapp \
  --output ../myapp \
  --display-name "My App" \
  --description "My terminal app"
```

## Features

- Bubble Tea v2 app shell
- Charm stack conventions
- Command palette
- Header/sidebar/main/footer layout
- Screen router
- Global keybindings
- Help overlay
- Theme system with tuitheme built-in themes
- Logs/debug screen
- GoReleaser release pipeline

## Development

Run the Bubbleplate demo shell:

```sh
bubbleplate demo
```

## Test

```sh
go test ./...
```

## Snapshot Release Build

```sh
goreleaser release --snapshot --clean
```

## Docker

Manual publishes build and push multi-arch images to GitHub Container Registry:

```sh
docker run --rm -it ghcr.io/elpdev/bubbleplate:latest
```

## Release

```sh
git tag v0.1.0
git push origin v0.1.0
```

## Version

```sh
bubbleplate --version
```
