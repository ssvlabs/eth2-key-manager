package core

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

type mockedStorage struct {
	seed []byte
	err  error
}

func (s *mockedStorage) Name() string                                              { return "" }
func (s *mockedStorage) SaveWallet(wallet Wallet) error                            { return nil }
func (s *mockedStorage) OpenWallet() (Wallet, error)                               { return nil, nil }
func (s *mockedStorage) ListAccounts() ([]ValidatorAccount, error)                 { return nil, nil }
func (s *mockedStorage) SaveAccount(account ValidatorAccount) error                { return nil }
func (s *mockedStorage) OpenAccount(accountId uuid.UUID) (ValidatorAccount, error) { return nil, nil }
func (s *mockedStorage) SetEncryptor(encryptor Encryptor, password []byte)         {}
func (s *mockedStorage) SecurelyFetchPortfolioSeed() ([]byte, error)               { return s.seed, nil }
func (s *mockedStorage) SecurelySavePortfolioSeed(secret []byte) error             { return s.err }

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func _bigInt(input string) *big.Int {
	res, _ := new(big.Int).SetString(input, 10)
	return res
}

func TestMarshalingHDKey(t *testing.T) {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	tests := []struct {
		name string
		seed []byte
		path string
		err  error
	}{
		{
			name: "validation account derivation",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path: "/0/0/0", // after basePath
			err:  nil,
		},
		{
			name: "withdrawal account derivation",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path: "/0/0", // after basePath
			err:  nil,
		},
		{
			name: "Base account derivation (base path only)",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path: "", // after basePath
			err:  errors.New("invalid relative path. Example: /1/2/3"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//storage := storage(test.seed, nil)
			//err := storage.SecurelySavePortfolioSeed(test.seed)
			//require.NoError(t, err)

			// create the privKey
			key, err := MasterKeyFromSeed(test.seed, TestNetwork)
			require.NoError(t, err)

			hdKey, err := key.Derive(test.path)
			if test.err != nil {
				require.EqualError(t, test.err, err.Error())
				return
			} else {
				require.NoError(t, err)
			}

			// marshal and unmarshal
			byts, err := json.Marshal(hdKey)
			require.NoError(t, err)

			newKey := &HDKey{}
			err = json.Unmarshal(byts, newKey)
			require.NoError(t, err)

			// match
			require.Equal(t, hdKey.Path(), newKey.Path())
			require.Equal(t, hdKey.id.String(), newKey.id.String())
			require.Equal(t, hdKey.privKey.Marshal(), newKey.privKey.Marshal())
			require.Equal(t, hdKey.PublicKey().Marshal(), newKey.PublicKey().Marshal())
		})
	}
}

func TestDerivableKeyRelativePathDerivation(t *testing.T) {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	tests := []struct {
		name        string
		seed        []byte
		path        string
		err         error
		expectedKey *big.Int
	}{
		{
			name:        "validation account 0 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/0/0/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("16278447180917815188301017385774271592438483452880235255024605821259671216398"),
		},
		{
			name:        "validation account 1 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/1/0/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("22772506560955906640840029020628554414154538440282401807772339666252999598733"),
		},
		{
			name:        "validation account 2 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/2/0/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("39196384482644522441983190042722076264169843386078553516164086198183513560637"),
		},
		{
			name:        "validation account 3 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/3/0/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("28093661633617073106051830080274606181076423213304176144286257209925213345002"),
		},
		{
			name:        "withdrawal account 0 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/0/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("26551663876804375121305275007227133452639447817512639855729535822239507627836"),
		},
		{
			name:        "withdrawal account 1 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/1/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("35957947454275682122989949668683794518020231276710636838205992785623169821803"),
		},
		{
			name:        "withdrawal account 2 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/2/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("8862394884593725153617163219481465667794938944832130820949251394547028786321"),
		},
		{
			name:        "withdrawal account 3 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/3/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("24013488102538647731381570745201628464138315555327292772724806156501038782887"),
		},
		{
			name:        "Base account derivation (big index)",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/100/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("14004582289918763639923763455218870137436565566857894891588947000864308096613"),
		},
		{
			name:        "bad path",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "0/0", // after basePath
			err:         errors.New("invalid relative path. Example: /1/2/3"),
			expectedKey: nil,
		},
		{
			name:        "too large of an index, bad path",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/1000/0", // after basePath
			err:         errors.New("invalid relative path. Example: /1/2/3"),
			expectedKey: nil,
		},
		{
			name:        "not a relative path",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "m/0/0", // after basePath
			err:         errors.New("invalid relative path. Example: /1/2/3"),
			expectedKey: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			key, err := MasterKeyFromSeed(test.seed, MainNetwork)
			require.NoError(t, err)

			hdKey, err := key.Derive(test.path)
			if err != nil {
				if test.err != nil {
					require.Equal(t, test.err.Error(), err.Error())
				} else {
					t.Error(err)
				}
				return
			} else {
				require.NoError(t, test.err)
			}

			require.Equal(t, MainNetwork.FullPath(test.path), hdKey.Path())
			privkey, err := e2types.BLSPrivateKeyFromBytes(test.expectedKey.Bytes())
			require.NoError(t, err)
			require.Equal(t, privkey.PublicKey().Marshal(), hdKey.PublicKey().Marshal())
		})
	}
}
