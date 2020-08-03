package main

import (
	"github.com/bloxapp/KeyVault/cli/cmd"
	_ "github.com/bloxapp/KeyVault/cli/cmd/mnemonic"
	_ "github.com/bloxapp/KeyVault/cli/cmd/seed"
	_ "github.com/bloxapp/KeyVault/cli/cmd/wallet"
	_ "github.com/bloxapp/KeyVault/cli/cmd/wallet/cmd/account"
)

func main() {
	cmd.Execute()
}
