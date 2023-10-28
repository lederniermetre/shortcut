# shortcut

Small helper for missing features on shortcut.com

## Usage

```bash
go run cmd/cli/main.go -iteration "Iteration name" -debug
```

Parameters

- `-iterration` is a search pattern.
- `-debug` set Debug level on logger

## Development

### Requirements

- Swagger cli

## Init

```bash
swagger generate client -f https://developer.shortcut.com/api/rest/v3/shortcut.swagger.json --target pkg/shortcut/gen/
```
