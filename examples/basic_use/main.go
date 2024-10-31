package main

import (
	"encoding/hex"
	"fmt"

	eth2keymanager "github.com/ssvlabs/eth2-key-manager"
	"github.com/ssvlabs/eth2-key-manager/core"
	"github.com/ssvlabs/eth2-key-manager/stores/inmemory"
)

func main() {
	entropy, _ := core.GenerateNewEntropy()

	// print out mnemonic
	mnemonic, _ := core.EntropyToMnemonic(entropy)
	fmt.Printf("Generated mnemonic: %s\n", mnemonic)

	// generate seed
	seed, _ := core.SeedFromEntropy(entropy, "")

	// create storage
	store := inmemory.NewInMemStore(core.PraterNetwork)

	// create options
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)

	// instantiate KeyVaul
	vault, _ := eth2keymanager.NewKeyVault(options)

	// create account
	wallet, _ := vault.Wallet()
	account, _ := wallet.CreateValidatorAccount(seed, nil)

	fmt.Printf("created validator account with pub key: %s\n", hex.EncodeToString(account.ValidatorPublicKey()))

}
