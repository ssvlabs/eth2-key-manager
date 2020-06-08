package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	"math/big"
	"os"
	"testing"
)

const (
	basePath = "m/12381/3600"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func _bigInt(input string) *big.Int {
	res, _ := new(big.Int).SetString(input, 10)
	return res
}

func TestMarshaling(t *testing.T) {
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
			// create the key
			key,err := BaseKeyFromSeed(test.seed)
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
			require.Equal(t,key.Path,newKey.Path)
			require.Equal(t,key.Key.Marshal(),newKey.Key.Marshal())
			require.Equal(t,key.seed,newKey.seed)
		})
	}
}

func TestRelativePathDerivation(t *testing.T) {
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
			key,err := BaseKeyFromSeed(test.seed)
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

			assert.Equal(t,basePath + test.path,key.Path)
			assert.Equal(t,test.expectedKey.Bytes(),key.Key.Marshal())
		})
	}
}
