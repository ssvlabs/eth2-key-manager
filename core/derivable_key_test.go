package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
	"math/big"
	"os"
	"testing"
)

const (
	basePath = "m/12381/3600"
)

type mockedStorage struct {
	seed []byte
	err error
}
func (s *mockedStorage) Name() string {return ""}
func (s *mockedStorage) SavePortfolio(portfolio Portfolio) error {return nil}
func (s *mockedStorage) OpenPortfolio() (Portfolio,error) {return nil,nil}
func (s *mockedStorage) ListWallets() ([]Wallet,error) {return nil,nil}
func (s *mockedStorage) SaveWallet(wallet Wallet) error {return nil}
func (s *mockedStorage) OpenWallet(uuid uuid.UUID) (Wallet,error) {return nil,nil}
func (s *mockedStorage) ListAccounts(walletID uuid.UUID) ([]Account,error) {return nil,nil}
func (s *mockedStorage) SaveAccount(account Account) error {return nil}
func (s *mockedStorage) OpenAccount(uuid uuid.UUID) (Account,error) {return nil,nil}
func (s *mockedStorage) SetEncryptor(encryptor types.Encryptor, password []byte) {}
func (s *mockedStorage) SecurelyFetchPortfolioSeed() ([]byte,error) {return s.seed,nil}
func (s *mockedStorage) SecurelySavePortfolioSeed(secret []byte) error {return s.err}

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func _bigInt(input string) *big.Int {
	res, _ := new(big.Int).SetString(input, 10)
	return res
}

func storage(seed []byte, err error) Storage {
	return &mockedStorage{seed:seed,err:err}
}

func TestMarshalingDerivableKey(t *testing.T) {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	tests := []struct{
		name string
		seed []byte
		path string
		err  error
	} {
		{
			name: "validation account derivation",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"),
			path:  "/0/0/0", // after basePath
			err: nil,
		},
		{
			name: "withdrawal account derivation",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"),
			path:  "/0/0", // after basePath
			err: nil,
		},
		{
			name: "Base account derivation (base path only)",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"),
			path:  "", // after basePath
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := storage(test.seed,test.err)
			err := storage.SecurelySavePortfolioSeed(test.seed)
			if err != nil {
				t.Error(err)
				return
			}

			// create the privKey
			key,err := BaseKeyFromSeed(test.seed, storage)
			if err != nil {
				t.Error(err)
				return
			}
			if len(test.path) > 0 {
				key,err = key.Derive(test.path)
				if err != nil {
					t.Error(err)
					return
				}
			}

			// marshal and unmarshal
			byts,err := json.Marshal(key)
			if err != nil {
				t.Error(err)
				return
			}
			newKey := &DerivableKey{}
			err = json.Unmarshal(byts,newKey)
			if err != nil {
				t.Error(err)
				return
			}

			// match
			require.Equal(t,key.Path(),newKey.Path())
			require.Equal(t,key.id.String(),newKey.id.String())
			require.Equal(t,key.PublicKey().Marshal(),newKey.PublicKey().Marshal())
		})
	}
}

func TestDerivableKeyRelativePathDerivation(t *testing.T) {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	tests := []struct {
		name string
		seed []byte
		path string
		err  error
		expectedKey *big.Int
	}{
		{
			name: "validation account derivation",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"),
			path:  "/0/0/0", // after basePath
			err: nil,
			expectedKey: _bigInt("34956507762652225510006274457607191231356413500508203267956091970659502095935"),
		},
		{
			name: "withdrawal account derivation",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"),
			path:  "/0/0", // after basePath
			err: nil,
			expectedKey: _bigInt("31676788419929922777864946442677915531199062343799598297489487887255736884383"),
		},
		{
			name: "Base account derivation (big index)",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"),
			path:  "/100/0", // after basePath
			err: nil,
			expectedKey: _bigInt("40407741422272659004469348930958444733127038589615463764403690368629477657256"),
		},
		{
			name: "bad path",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"),
			path:  "0/0", // after basePath
			err: fmt.Errorf("invalid relative path. Example: /1/2/3"),
			expectedKey: nil,
		},
		{
			name: "too large of an index, bad path",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"),
			path:  "/1000/0", // after basePath
			err: fmt.Errorf("invalid relative path. Example: /1/2/3"),
			expectedKey: nil,
		},
		{
			name: "not a relative path",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"),
			path:  "m/0/0", // after basePath
			err: fmt.Errorf("invalid relative path. Example: /1/2/3"),
			expectedKey: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := storage(test.seed,test.err)
			key,err := BaseKeyFromSeed(test.seed, storage)
			if err != nil {
				t.Error(err)
				return
			}

			key,err = key.Derive(test.path)
			if err != nil {
				if test.err != nil {
					assert.Equal(t,test.err.Error(),err.Error())
				} else {
					t.Error(err)
				}
				return
			} else {
				if test.err != nil {
					t.Errorf("should have returned error but didn't")
					return
				}
			}

			assert.Equal(t,basePath + test.path,key.Path())
			privkey,err := e2types.BLSPrivateKeyFromBytes(test.expectedKey.Bytes())
			assert.NoError(t,err)
			assert.Equal(t,privkey.PublicKey().Marshal(),key.PublicKey().Marshal())
		})
	}
}
