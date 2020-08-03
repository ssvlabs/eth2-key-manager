package core

import "github.com/tyler-smith/go-bip39"

func GenerateNewEntropy() ([]byte, error) {
	seed, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, err
	}

	return seed, nil
}

// Given an entropy, create the mnemonic passphrase.
func EntropyToMnemonic(entropy []byte) (string, error) {
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

// the seed is the product of applying a key derivation algo (PBKDF2) on the mnemonic (as the entropy)
// and the password as salt.
// please see https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki
func SeedFromMnemonic(mnemonic string, password string) ([]byte, error) {
	return bip39.NewSeed(mnemonic, password), nil
}


// the seed is the product of applying a key derivation algo (PBKDF2) on the mnemonic (as the entropy)
// and the password as salt.
// please see https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki
func SeedFromEntropy(entropy []byte, password string) ([]byte, error) {
	mnemonic,err := EntropyToMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	return bip39.NewSeed(mnemonic, password), nil
}

