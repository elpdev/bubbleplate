# Bubbleplate

Bubbleplate is an opinionated generator for building Go TUI projects with the [Charm](https://charm.sh) stack: Bubble Tea, Lip Gloss, and Bubbles.

## Generate a Project

```sh
go run ./cmd/bubbleplate new myapp --module github.com/acme/myapp
cd myapp
go mod tidy
go test ./...
go run ./cmd/myapp
```

Useful options:

```sh
go run ./cmd/bubbleplate new myapp \
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
- Theme system with Phosphor, Muted Dark, and Miami themes
- Logs/debug screen
- GoReleaser release pipeline

## Development

Run the Bubbleplate demo shell:

```sh
go run ./cmd/bubbleplate demo
```

## Install

Homebrew:

```sh
brew install elpdev/tap/bubbleplate
```

Arch Linux via AUR:

```sh
yay -S bubbleplate-bin
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
go run ./cmd/bubbleplate --version
```
