# [[ .DisplayName ]]

[[ .Description ]]

## Features

- Bubble Tea v2 app shell
- Command palette
- Header/sidebar/main/footer layout
- Screen router
- Global keybindings
- Help overlay
- Theme system with tuitheme built-in themes
- Logs/debug screen
- GoReleaser release pipeline

## Development

```sh
go run ./cmd/[[ .BinaryName ]]
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
docker run --rm -it [[ .DockerImage ]]:latest
```

## Release

```sh
git tag v0.1.0
git push origin v0.1.0
```

## Version

```sh
go run ./cmd/[[ .BinaryName ]] --version
```
