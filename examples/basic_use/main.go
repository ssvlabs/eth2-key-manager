package main

import (
	"encoding/hex"
	"fmt"

	"github.com/bloxapp/eth-key-manager"
	"github.com/bloxapp/eth-key-manager/core"
	"github.com/bloxapp/eth-key-manager/stores/in_memory"
)

func main() {
	entropy, _ := core.GenerateNewEntropy()

	// print out mnemonic
	mnemonic, _ := core.EntropyToMnemonic(entropy)
	fmt.Printf("Generated mnemonic: %s\n", mnemonic)

	// generate seed
	seed, _ := core.SeedFromEntropy(entropy, "")

	// create storage
	store := in_memory.NewInMemStore()

	// create options
	options := &KeyVault.KeyVaultOptions{}
	options.SetStorage(store)

	// instantiate KeyVaul
	vault, _ := KeyVault.NewKeyVault(options)

	// create account
	wallet, _ := vault.Wallet()
	account, _ := wallet.CreateValidatorAccount(seed, "account test")

	fmt.Printf("created validator account with pub key: %s\n", hex.EncodeToString(account.ValidatorPublicKey().Marshal()))

}
