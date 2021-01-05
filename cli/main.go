package main

import (
	"github.com/bloxapp/eth2-key-manager/cli/cmd"
	_ "github.com/bloxapp/eth2-key-manager/cli/cmd/config"
	_ "github.com/bloxapp/eth2-key-manager/cli/cmd/mnemonic"
	_ "github.com/bloxapp/eth2-key-manager/cli/cmd/seed"
	_ "github.com/bloxapp/eth2-key-manager/cli/cmd/validator"
	_ "github.com/bloxapp/eth2-key-manager/cli/cmd/wallet"
	_ "github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account"
	_ "github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/publickey"
)

func main() {
	cmd.Execute()
}
