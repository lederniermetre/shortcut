# shortcut

[![Go Report Card](https://goreportcard.com/badge/github.com/lederniermetre/shortcut)](https://goreportcard.com/report/github.com/lederniermetre/shortcut) [![CI](https://github.com/lederniermetre/shortcut/actions/workflows/ci.yaml/badge.svg)](https://github.com/lederniermetre/shortcut/actions/workflows/ci.yaml)

Small helper for missing features on shortcut.com

## Features

- Allow to get by users the addition of estimate
- Display stories not finish in the previous iteration

## Usage

```bash
go run cmd/cli/main.go -iteration "Iteration name" -debug
```

Parameters

- `-iterration` is a search pattern.
- `-debug` set Debug level on logger

## Example

```bash
go run cmd/cli/main.go -iteration "Ops 61" -debug
Oct 29 11:48:19 INFO  cmd/cli/main.go:75 Retrieve iteration informations name="#61 OPS"
Oct 29 11:48:19 DEBUG cmd/cli/main.go:95 Compute story name="[IaC] Includes defaults in provider" owners="0" estimate="3"
Oct 29 11:48:19 WARN  cmd/cli/main.go:103 Story has no owners name="[IaC] Includes defaults in provider"
Oct 29 11:48:19 WARN  cmd/cli/main.go:91 OMG no estimate on story: Decomission Service A
Oct 29 11:48:19 WARN  cmd/cli/main.go:91 OMG no estimate on story: ETQ OPS I want to setup Service C
Oct 29 11:48:19 WARN  cmd/cli/main.go:91 OMG no estimate on story: Decomission Service B
Oct 29 11:48:19 WARN  cmd/cli/main.go:91 OMG no estimate on story: [IaC] PRA mode
Oct 29 11:48:19 DEBUG cmd/cli/main.go:95 Compute story name="[Cassandra] Update client" owners="2" estimate="5"
Oct 29 11:48:19 DEBUG cmd/cli/main.go:109 Story shared, split estimate name="[Cassandra] Update client"
Oct 29 11:48:19 INFO  cmd/cli/main.go:135 John Doe has 2 of load
Oct 29 11:48:19 INFO  cmd/cli/main.go:135 Michel Paul has 2 of load
```

## Development

### Requirements

- [Swagger cli](https://goswagger.io/install.html)
- Golang

## Init

```bash
swagger generate client -f https://developer.shortcut.com/api/rest/v3/shortcut.swagger.json --target pkg/shortcut/gen/
```
