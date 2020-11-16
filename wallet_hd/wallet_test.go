package wallet_hd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/bloxapp/eth2-key-manager/encryptor"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth2-key-manager/core"
)

type dummyStorage struct{}

func (s *dummyStorage) Name() string                                    { return "" }
func (s *dummyStorage) Network() core.Network                           { return core.MainNetwork }
func (s *dummyStorage) SaveWallet(wallet core.Wallet) error             { return nil }
func (s *dummyStorage) OpenWallet() (core.Wallet, error)                { return nil, nil }
func (s *dummyStorage) ListAccounts() ([]core.ValidatorAccount, error)  { return nil, nil }
func (s *dummyStorage) SaveAccount(account core.ValidatorAccount) error { return nil }
func (s *dummyStorage) OpenAccount(accountId uuid.UUID) (core.ValidatorAccount, error) {
	return nil, nil
}
func (s *dummyStorage) DeleteAccount(accountId uuid.UUID) error                     { return nil }
func (s *dummyStorage) SetEncryptor(encryptor encryptor.Encryptor, password []byte) {}

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func _bigIntFromSkHex(input string) *big.Int {
	res, _ := new(big.Int).SetString(input, 16)
	return res
}

func storage() core.Storage {
	return &dummyStorage{}
}

// ethereum foundation launched a launchpad for making deposits.
// this test compares the launchpad results and KeyVault
// Updated for V1.0.0 https://github.com/ethereum/eth2.0-deposit-cli/releases/tag/v1.0.0
func TestAccountDerivationComparedToOfficialLaunchPad(t *testing.T) {
	e2types.InitBLS()

	tests := []struct {
		mnemonic        string
		password        string
		validatorPubKey string
	}{
		{
			mnemonic:        "vocal differ audit mom unique physical evolve cave retire design achieve pupil odor hockey drive animal habit fluid belt height vintage crack rigid sphere",
			password:        "",
			validatorPubKey: "858da1ba93ea436c93c8f1aac0d508130da2696f7394f9f88088b35050670f6dbf6a9d491cd386d28420fbef684c48e0",
		},
		{
			mnemonic:        "magnet burden popular race night clown moral sorry situate worth sorry solution live custom message finger soon month invest battle fade funny bright basket",
			password:        "",
			validatorPubKey: "b784d12bedcb1469200e8b2a0b00bcbe0be4cda19f0d05c307df8c16ec7b3a1f6244e23fc71f5e3d1e24f0c2231a3e03",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test vector: %d", i), func(t *testing.T) {
			seed, err := core.SeedFromMnemonic(test.mnemonic, test.password)
			require.NoError(t, err)

			//
			storage := storage()
			w := &HDWallet{
				id:          uuid.New(),
				indexMapper: make(map[string]uuid.UUID),
				context: &core.WalletContext{
					Storage: storage,
				},
			}
			account, err := w.CreateValidatorAccount(seed, nil)
			require.NoError(t, err)
			require.Equal(t, test.validatorPubKey, hex.EncodeToString(account.ValidatorPublicKey().Marshal()))
		})
	}
}

func TestAccountDerivation(t *testing.T) {
	e2types.InitBLS()

	// create wallet
	storage := storage()
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	w := &HDWallet{
		id:          uuid.New(),
		indexMapper: make(map[string]uuid.UUID),
		context: &core.WalletContext{
			Storage: storage,
		},
	}

	tests := []struct {
		testName              string
		index                 int
		expectedValidationKey *big.Int
		expectedWithdrawalKey *big.Int
	}{
		{
			testName:              "account 0",
			index:                 0,
			expectedValidationKey: _bigIntFromSkHex("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"),
			expectedWithdrawalKey: _bigIntFromSkHex("a0b9324da8a8a696c53950e984de25b299c123d17bab972eca1ac2c674964c9f817047bc6048ef0705d7ec6fae6d5da6"),
		},
		{
			testName:              "account 1",
			index:                 1,
			expectedValidationKey: _bigIntFromSkHex("b41df3c322a6fd305fc9425df52501f7f8067dbba551466d82d506c83c6ab287580202aa1a3449f54b9bc464a04b70e6"),
			expectedWithdrawalKey: _bigIntFromSkHex("858e30df33bfdd613234abc9359ccd924f4807f1ba21de328d361e72f8c9ca94c9b7c225536405141df8239db87bd510"),
		},
		{
			testName:              "account 2",
			index:                 2,
			expectedValidationKey: _bigIntFromSkHex("9415b51f7996d6872f32c9bf7c259fad10e211d6097ff52ae99520a0ab3b916b3570073abbb83fa87da66936d351010d"),
			expectedWithdrawalKey: _bigIntFromSkHex("85586894abb77e41ba5dc3cfa2a7506c7584d024f028501da1e766792bcf6cd79ae17ff68eee84315eba9a2a8e7f89fe"),
		},
		{
			testName:              "account 3",
			index:                 3,
			expectedValidationKey: _bigIntFromSkHex("80b42ed53fe82598d055c2723bce9b1dde249d0497291856ef77fc75b094c60aca9dcc648e414dc9db41f8b8dc2f13e4"),
			expectedWithdrawalKey: _bigIntFromSkHex("afb22992f52aaf46c461ad1013e88c2c3ca8656c58170a9d08aaaeb9eac404fba839b313150f8f4b2f9fe23e64119c1f"),
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			account, err := w.CreateValidatorAccount(seed, &test.index)
			require.NoError(t, err)

			val, err := e2types.BLSPublicKeyFromBytes(test.expectedValidationKey.Bytes())
			require.NoError(t, err)
			with, err := e2types.BLSPublicKeyFromBytes(test.expectedWithdrawalKey.Bytes())
			require.NoError(t, err)

			require.Equal(t, val.Marshal(), account.ValidatorPublicKey().Marshal(), fmt.Sprintf("expceted validation pk: %s\n", hex.EncodeToString(account.ValidatorPublicKey().Marshal())))
			require.Equal(t, with.Marshal(), account.WithdrawalPublicKey().Marshal(), fmt.Sprintf("expceted withdrawal pk: %s\n", hex.EncodeToString(account.WithdrawalPublicKey().Marshal())))
		})
	}
}

func TestCreateAccounts(t *testing.T) {
	tests := []struct {
		testName        string
		newAccounttName string
		expectedErr     string
	}{
		{
			testName:        "Add new account",
			newAccounttName: "account1",
			expectedErr:     "",
		},
	}

	// create key and wallet
	storage := storage()
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")

	w := &HDWallet{
		id:          uuid.New(),
		indexMapper: make(map[string]uuid.UUID),
		//key:key,
		context: &core.WalletContext{
			Storage: storage,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			_, err := w.CreateValidatorAccount(seed, nil)
			if test.expectedErr != "" {
				require.Errorf(t, err, test.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestWalletMarshaling(t *testing.T) {
	tests := []struct {
		id          uuid.UUID
		testName    string
		walletType  core.WalletType
		indexMapper map[string]uuid.UUID
		seed        []byte
		path        string
	}{
		{
			testName:    "simple wallet, no account",
			id:          uuid.New(),
			walletType:  core.HDWallet,
			indexMapper: map[string]uuid.UUID{},
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/0/0",
		},
		{
			testName:   "simple wallet with accounts",
			id:         uuid.New(),
			walletType: core.HDWallet,
			indexMapper: map[string]uuid.UUID{
				"ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3279": uuid.New(),
			},
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path: "/0/0",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// setup storage
			storage := storage()

			w := &HDWallet{
				walletType:  test.walletType,
				id:          test.id,
				indexMapper: test.indexMapper,
				//key:key,
			}

			// marshal
			byts, err := json.Marshal(w)
			require.NoError(t, err)

			//unmarshal
			w1 := &HDWallet{context: &core.WalletContext{Storage: storage}}
			err = json.Unmarshal(byts, w1)
			require.NoError(t, err)

			require.Equal(t, w.id, w1.id)
			require.Equal(t, w.walletType, w1.walletType)
			for k := range w.indexMapper {
				v := w.indexMapper[k]
				require.Equal(t, v, w1.indexMapper[k])
			}
		})
	}

}
