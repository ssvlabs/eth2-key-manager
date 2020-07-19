# KeyVault CLI

This is the CLI for using KeyVault functionality through command line interface.

## Files structure

`main.go` includes all packages that contain CLI commands.

- `./cmd` contains commands and its handlers.
- `./util` contains utilities related to CLI commands.
- `./portfolio/flag` (doesn't exists) contains flags/parameters of the command. 
    There should be defined all flags related to portfolio and its sub-commands.
- `./portfolio/cmd/seed/handler` contains handlers of the portfolio seed related commands.

## How to build

Execute the following commands from the root of the repo:

```bash
$ go build -o keyvault-cli ./cli
```