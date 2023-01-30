# Eth Key Manager CLI

This is the CLI for using KeyVault functionality through command line interface.

## Files structure

`main.go` includes all packages that contain CLI commands.

- `./cmd` contains commands and its handlers.
- `./util` contains utilities related to CLI commands.

## How to build

Execute the following commands from the root of the repo:

Mac:
```bash
$ go build -o keyvault-cli ./cli
```

Windows:
build win cli on mac machine

1. brew install mingw-w64
```bash
$ brew install mingw-w64
```
2. build bin
```bash
$ env CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -o keyvault-cli.exe ./cli
```
