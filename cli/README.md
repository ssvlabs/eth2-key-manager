# Eth Key Manager CLI

This is the CLI for using KeyVault functionality through command line interface.

## Files structure

`main.go` includes all packages that contain CLI commands.

- `./cmd` contains commands and its handlers.
- `./util` contains utilities related to CLI commands.

## How to build

Execute the following commands from the root of the repo:

```bash
$ go build -o keyvault-cli ./cli
```

## Commands

- Create validator(s) included making deposits:
    ```sh
    $ keyvault-cli validator create \
      --wallet-private-key=<eth1-wallet-private-key> \
      --wallet-addr=<eth1-wallet-address> \
      --seeds-count=<seeds-count> \
      --validators-per-seed=<number-of-validators-per-seed>
    ```  
  [There](https://metamask.zendesk.com/hc/en-us/articles/360015289632-How-to-Export-an-Account-Private-Key) is a doc how to get a private key in MetaMask.