package main

import (
	"github.com/bloxapp/KeyVault/cli/cmd"
	_ "github.com/bloxapp/KeyVault/cli/cmd/portfolio"
	_ "github.com/bloxapp/KeyVault/cli/cmd/portfolio/cmd/seed"
)

func main() {
	cmd.Execute()
}
