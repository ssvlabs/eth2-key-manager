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

func _bigInt(input string) *big.Int {
	res, _ := new(big.Int).SetString(input, 10)
	return res
}

func storage() core.Storage {
	return &dummyStorage{}
}

// ethereum foundation launched a launchpad for making deposits.
// this test compares the launchpad results and KeyVault
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
			validatorPubKey: "aa6e59b378b905a7454cf3a7a57e07ce89d5410fb3f96610aba0a8036984d7a6a2e1398fceb85611bc576cd349d7dcd2",
		},
		{
			mnemonic:        "magnet burden popular race night clown moral sorry situate worth sorry solution live custom message finger soon month invest battle fade funny bright basket",
			password:        "",
			validatorPubKey: "98ee5f5107f72bef05f59dee8de08223cd0db04be9f142806743cad44366f44184f0957935dc8e4994e038ad1fd5d821",
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
			expectedValidationKey: _bigInt("16278447180917815188301017385774271592438483452880235255024605821259671216398"),
			expectedWithdrawalKey: _bigInt("26551663876804375121305275007227133452639447817512639855729535822239507627836"),
		},
		{
			testName:              "account 1",
			index:                 1,
			expectedValidationKey: _bigInt("22772506560955906640840029020628554414154538440282401807772339666252999598733"),
			expectedWithdrawalKey: _bigInt("35957947454275682122989949668683794518020231276710636838205992785623169821803"),
		},
		{
			testName:              "account 2",
			index:                 2,
			expectedValidationKey: _bigInt("39196384482644522441983190042722076264169843386078553516164086198183513560637"),
			expectedWithdrawalKey: _bigInt("8862394884593725153617163219481465667794938944832130820949251394547028786321"),
		},
		{
			testName:              "account 3",
			index:                 3,
			expectedValidationKey: _bigInt("28093661633617073106051830080274606181076423213304176144286257209925213345002"),
			expectedWithdrawalKey: _bigInt("24013488102538647731381570745201628464138315555327292772724806156501038782887"),
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			account, err := w.CreateValidatorAccount(seed, &test.index)
			require.NoError(t, err)

			val, err := e2types.BLSPrivateKeyFromBytes(test.expectedValidationKey.Bytes())
			require.NoError(t, err)
			with, err := e2types.BLSPrivateKeyFromBytes(test.expectedWithdrawalKey.Bytes())
			require.NoError(t, err)

			require.Equal(t, val.PublicKey().Marshal(), account.ValidatorPublicKey().Marshal())
			require.Equal(t, with.PublicKey().Marshal(), account.WithdrawalPublicKey().Marshal())
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
