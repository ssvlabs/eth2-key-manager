package hashicorp

import (
	"context"
	"encoding/hex"
	"github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	"testing"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func baseKeyVault(seed []byte, acc1Name string, acc2Name string, t *testing.T) (core.Storage, core.Wallet, []core.ValidatorAccount) {
	// store
	inMemStore := in_memory.NewInMemStore()
	// seed
	// create keyvault in a normal in mem store
	options := &KeyVault.KeyVaultOptions{}
	options.SetStorage(inMemStore)
	options.SetSeed(seed)
	options.SetEncryptor(keystorev4.New())
	options.SetPassword("password")
	kv, err := KeyVault.NewKeyVault(options)
	require.NoError(t, err)
	require.NotNil(t, kv)
	// get wallet and accounts to compare
	inMemWallet, err := kv.Wallet()
	require.NoError(t, err)
	require.NotNil(t, inMemWallet)
	inMemAcc1, err := inMemWallet.CreateValidatorAccount(seed, acc1Name)
	require.NoError(t, err)
	require.NotNil(t, inMemAcc1)
	inMemAcc2, err := inMemWallet.CreateValidatorAccount(seed, acc2Name)
	require.NoError(t, err)
	require.NotNil(t, inMemAcc2)

	return inMemStore, inMemWallet, []core.ValidatorAccount{inMemAcc1, inMemAcc2}
}

func TestImportAndDeleteFromInMem(t *testing.T) {
	oldInMemStore, _, oldInMemAccounts := baseKeyVault(
		_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
		"acc1",
		"acc2",
		t,
	)

	hashiStorage := &logical.InmemStorage{}

	// import to hashicorp
	oldHashi, err := FromInMemoryStore(oldInMemStore.(*in_memory.InMemStore), hashiStorage, context.Background())
	require.NoError(t, err)

	// create another in mem base keyvault to override (different seed and account names)
	inMemStore, inMemWallet, inMemAccounts := baseKeyVault(
		_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fdf"),
		"acc3",
		"acc4",
		t,
	)

	// import to hashicorp, should override
	hashi, err := FromInMemoryStore(inMemStore.(*in_memory.InMemStore), hashiStorage, context.Background())
	require.NoError(t, err)

	// verify deletion
	// accounts fetched should no longer match old accounts
	res, err := oldHashi.OpenAccount(oldInMemAccounts[0].ID())
	require.Nil(t, res)
	res, err = oldHashi.OpenAccount(oldInMemAccounts[1].ID())
	require.Nil(t, res)

	// get hasicorp's wallet and accounts
	hashiWallet, err := hashi.OpenWallet()
	require.NoError(t, err)
	require.NotNil(t, hashiWallet)
	hashiAcc3, err := hashiWallet.AccountByPublicKey("90c7ad0b421a767b6d03c5b06fe57a6b9826dd5658dfccc069d85122451433b45014f0f4f7a6afc8bbe3b24c2344019b")
	require.NoError(t, err)
	require.NotNil(t, hashiAcc3)
	hashiAcc4, err := hashiWallet.AccountByPublicKey("813e00d6d7fae79275ca1552383dfbfc6a7db09baa9d981d797b2338e205e00bf2d8a7f9038b0d8e3d8e316ec555ecf6")
	require.NoError(t, err)
	require.NotNil(t, hashiAcc4)

	// compare
	require.Equal(t, inMemWallet.ID().String(), hashiWallet.ID().String())
	require.Equal(t, inMemAccounts[0].ID().String(), hashiAcc3.ID().String())
	require.Equal(t, inMemAccounts[0].ValidatorPublicKey().Marshal(), hashiAcc3.ValidatorPublicKey().Marshal())
	require.Equal(t, inMemAccounts[1].ID().String(), hashiAcc4.ID().String())
	require.Equal(t, inMemAccounts[1].ValidatorPublicKey().Marshal(), hashiAcc4.ValidatorPublicKey().Marshal())
}

func TestImportFromInMem(t *testing.T) {
	inMemStore, inMemWallet, inMemAccounts := baseKeyVault(
		_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
		"acc1",
		"acc2",
		t,
	)

	// import to hashicorp
	hashi, err := FromInMemoryStore(inMemStore.(*in_memory.InMemStore), &logical.InmemStorage{}, context.Background())
	require.NoError(t, err)

	// get hasicorp's wallet and accounts
	hashiWallet, err := hashi.OpenWallet()
	require.NoError(t, err)
	require.NotNil(t, hashiWallet)
	hashiAcc1, err := hashiWallet.AccountByPublicKey("ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3279")
	require.NoError(t, err)
	require.NotNil(t, hashiAcc1)
	hashiAcc2, err := hashiWallet.AccountByPublicKey("aedfb8f364cba8951565cf275d939a008a1bc8b8a26cff18681091da165cf7478a6bd71bbac87a6eae70d68fab384efd")
	require.NoError(t, err)
	require.NotNil(t, hashiAcc2)

	// compare
	require.Equal(t, inMemWallet.ID().String(), hashiWallet.ID().String())
	require.Equal(t, inMemAccounts[0].ID().String(), hashiAcc1.ID().String())
	require.Equal(t, inMemAccounts[0].ValidatorPublicKey().Marshal(), hashiAcc1.ValidatorPublicKey().Marshal())
	require.Equal(t, inMemAccounts[1].ID().String(), hashiAcc2.ID().String())
	require.Equal(t, inMemAccounts[1].ValidatorPublicKey().Marshal(), hashiAcc2.ValidatorPublicKey().Marshal())
}
