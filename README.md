# Supercharged shortcut

[![Go Report Card](https://goreportcard.com/badge/github.com/lederniermetre/shortcut)](https://goreportcard.com/report/github.com/lederniermetre/shortcut) [![CI](https://github.com/lederniermetre/shortcut/actions/workflows/ci.yaml/badge.svg)](https://github.com/lederniermetre/shortcut/actions/workflows/ci.yaml)

Supercharged [shortcut](https://shortcut.com/) cli with stats. Powered by [cobra](github.com/spf13/cobra) and [pterm](https://github.com/pterm/pterm)

## Features

- Estimate charge by owners
- Report progression by Epic: stories and estimates
- Report global progress and stats of iteration
- Report postponed stories between iterations

## Installation

Get latest version from [release page](https://github.com/lederniermetre/shortcut/releases) or run `curl https://raw.githubusercontent.com/lederniermetre/shortcut/main/install.sh | bash`.

## Usage

```bash
shortcut help
```

Parameters

- `--iteration|-i` is a search pattern.
- `--debug|-d` set Debug level on logger
- `SHORTCUT_API_TOKEN` environment variable

## Example

```bash
shortcut iteration stats owners
Dec  1 21:42:34 INFO  Iteration retrieved name="#63 OPS"
WARNING  Story has no owners: ETQ SRE i love tests
WARNING  Story has no owners: ETQ SRE i want s3

Load by owners

Pierre Gasly ███████████████  1
Carlos Sainz ███████████████  1
```

## Development

### Requirements

- [Swagger cli](https://goswagger.io/install.html)
- Golang
- [task](https://taskfile.dev/installation/)

## Help

```bash
task --list-all
```

You can pass arguments to task:
```bash
task dev -- iteration stats owners
```
