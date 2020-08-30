package main

import (
	"github.com/bloxapp/eth-key-manager/cli/cmd"
	_ "github.com/bloxapp/eth-key-manager/cli/cmd/mnemonic"
	_ "github.com/bloxapp/eth-key-manager/cli/cmd/seed"
	_ "github.com/bloxapp/eth-key-manager/cli/cmd/wallet"
	_ "github.com/bloxapp/eth-key-manager/cli/cmd/wallet/cmd/account"
)

func main() {
	cmd.Execute()
}
