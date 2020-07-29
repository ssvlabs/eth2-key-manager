package main

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/stores/in_memory"
)

func main() {
	seed, _ := KeyVault.GenerateNewSeed()

	// print out mnemonic
	mnemonic, _ := KeyVault.SeedToMnemonic(seed)
	fmt.Printf("Generated mnemonic: %s\n", mnemonic)

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
