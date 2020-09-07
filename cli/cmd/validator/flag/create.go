package flag

import (
	"github.com/spf13/cobra"

	"github.com/bloxapp/eth2-key-manager/cli/util/cliflag"
)

// Flag names.
const (
	seedsCountFlag        = "seeds-count"
	validatorsPerSeedFlag = "validators-per-seed"
	walletAddrFlag        = "wallet-addr"
	walletPrivateKeyFlag  = "wallet-private-key"
	web3AddrFlag          = "web3-addr"
)

// Default values
const (
	web3AddrDefault = "https://goerli.prylabs.net"
)

// AddSeedsCountFlag adds the seeds count flag to the command
func AddSeedsCountFlag(c *cobra.Command) {
	cliflag.AddPersistentIntFlag(c, seedsCountFlag, 0, "count of seed to generate", true)
}

// GetSeedsCountFlagValue gets the seeds count flag from the command
func GetSeedsCountFlagValue(c *cobra.Command) (int, error) {
	return c.Flags().GetInt(seedsCountFlag)
}

// AddValidatorsPerSeedFlag adds the validators per seed flag to the command
func AddValidatorsPerSeedFlag(c *cobra.Command) {
	cliflag.AddPersistentIntFlag(c, validatorsPerSeedFlag, 0, "number of validators per one seed", true)
}

// GetValidatorsPerSeedFlagValue gets the validators per seed flag from the command
func GetValidatorsPerSeedFlagValue(c *cobra.Command) (int, error) {
	return c.Flags().GetInt(validatorsPerSeedFlag)
}

// AddWalletPrivateKeyFlag adds the wallet private key flag to the command
func AddWalletPrivateKeyFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, walletPrivateKeyFlag, "", "private key of ETH wallet", true)
}

// GetWalletPrivateKeyFlagValue gets the wallet private key flag from the command
func GetWalletPrivateKeyFlagValue(c *cobra.Command) (string, error) {
	return c.Flags().GetString(walletPrivateKeyFlag)
}

// AddWalletAddressFlag adds the wallet address flag to the command
func AddWalletAddressFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, walletAddrFlag, "", "ETH wallet address", true)
}

// GetWalletAddressFlagValue gets the wallet address flag from the command
func GetWalletAddressFlagValue(c *cobra.Command) (string, error) {
	return c.Flags().GetString(walletAddrFlag)
}

// AddWeb3AddrFlag adds the web3 address flag to the command
func AddWeb3AddrFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, web3AddrFlag, web3AddrDefault, "An eth1 web3 provider string http endpoint", false)
}

// GetWeb3AddrFlagValue gets the web3 addr flag from the command
func GetWeb3AddrFlagValue(c *cobra.Command) (string, error) {
	return c.Flags().GetString(web3AddrFlag)
}
