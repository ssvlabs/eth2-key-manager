package keystorev4

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecrypt(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		passphrase string
		output     []byte
		err        string
	}{
		{
			name:       "NoCipher",
			input:      `{"checksum":{"function":"SHA256","message":"cb27fe860c96f269f7838525ba8dce0886e0b7753caccc14162195bcdacbf49e","params":{}},"kdf":{"function":"scrypt","message":"","params":{"dklen":32,"n":262144,"p":8,"r":1,"salt":"ab0c7876052600dd703518d6fc3fe8984592145b591fc8fb5c6d43190334ba19"}}}`,
			passphrase: "testpassword",
			err:        "no cipher",
		},
		{
			name:       "ShortPassphrase",
			input:      `{"checksum":{"function":"SHA256","message":"cb27fe860c96f269f7838525ba8dce0886e0b7753caccc14162195bcdacbf49e","params":{}},"cipher":{"function":"xor","message":"e18afad793ec8dc3263169c07add77515d9f301464a05508d7ecb42ced24ed3a","params":{}}}`,
			passphrase: "testpassword",
			err:        "decryption key must be at least 32 bytes",
		},
		{
			name:       "BadSalt",
			input:      `{"checksum":{"function":"SHA256","message":"cb27fe860c96f269f7838525ba8dce0886e0b7753caccc14162195bcdacbf49e","params":{}},"cipher":{"function":"xor","message":"e18afad793ec8dc3263169c07add77515d9f301464a05508d7ecb42ced24ed3a","params":{}},"kdf":{"function":"scrypt","message":"","params":{"dklen":32,"n":262144,"p":8,"r":1,"salt":"hb0c7876052600dd703518d6fc3fe8984592145b591fc8fb5c6d43190334ba19"}}}`,
			passphrase: "testpassword",
			err:        "invalid KDF salt",
		},
		{
			name:       "BadPRF",
			input:      `{"kdf":{"function":"pbkdf2","params":{"dklen":32,"c":262144,"prf":"hmac-sha128","salt":"d4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3"},"message":""},"checksum":{"function":"sha256","params":{},"message":"18b148af8e52920318084560fd766f9d09587b4915258dec0676cba5b0da09d8"},"cipher":{"function":"aes-128-ctr","params":{"iv":"264daa3f303d7259501c93d997d84fe6"},"message": "a9249e0ca7315836356e4c7440361ff22b9fe71e2e2ed34fc1eb03976924ed48"}}`,
			passphrase: "testpassword",
			err:        `unsupported PBKDF2 PRF "hmac-sha128"`,
		},
		{
			name:       "BadKDF",
			input:      `{"kdf":{"function":"magic","params":{"dklen":32,"c":262144,"prf":"hmac-sha128","salt":"d4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3"},"message":""},"checksum":{"function":"sha256","params":{},"message":"18b148af8e52920318084560fd766f9d09587b4915258dec0676cba5b0da09d8"},"cipher":{"function":"aes-128-ctr","params":{"iv":"264daa3f303d7259501c93d997d84fe6"},"message": "a9249e0ca7315836356e4c7440361ff22b9fe71e2e2ed34fc1eb03976924ed48"}}`,
			passphrase: "testpassword",
			err:        `unsupported KDF "magic"`,
		},
		{
			name:       "InvalidScryptParams",
			input:      `{"checksum":{"function":"SHA256","message":"cb27fe860c96f269f7838525ba8dce0886e0b7753caccc14162195bcdacbf49e","params":{}},"cipher":{"function":"xor","message":"e18afad793ec8dc3263169c07add77515d9f301464a05508d7ecb42ced24ed3a","params":{}},"kdf":{"function":"scrypt","message":"","params":{"dklen":0,"n":3,"p":-4,"r":1,"salt":"ab0c7876052600dd703518d6fc3fe8984592145b591fc8fb5c6d43190334ba19"}}}`,
			passphrase: "testpassword",
			err:        "invalid KDF parameters",
		},
		{
			name:       "BadCipherMessage",
			input:      `{"checksum":{"function":"SHA256","message":"cb27fe860c96f269f7838525ba8dce0886e0b7753caccc14162195bcdacbf49e","params":{}},"cipher":{"function":"xor","message":"h18afad793ec8dc3263169c07add77515d9f301464a05508d7ecb42ced24ed3a","params":{}},"kdf":{"function":"scrypt","message":"","params":{"dklen":32,"n":262144,"p":8,"r":1,"salt":"ab0c7876052600dd703518d6fc3fe8984592145b591fc8fb5c6d43190334ba19"}}}`,
			passphrase: "testpassword",
			err:        "invalid cipher message",
		},
		{
			name:       "BadChecksumMessage",
			input:      `{"checksum":{"function":"SHA256","message":"hb27fe860c96f269f7838525ba8dce0886e0b7753caccc14162195bcdacbf49e","params":{}},"cipher":{"function":"xor","message":"e18afad793ec8dc3263169c07add77515d9f301464a05508d7ecb42ced24ed3a","params":{}},"kdf":{"function":"scrypt","message":"","params":{"dklen":32,"n":262144,"p":8,"r":1,"salt":"ab0c7876052600dd703518d6fc3fe8984592145b591fc8fb5c6d43190334ba19"}}}`,
			passphrase: "testpassword",
			err:        "invalid checksum message",
		},
		{
			name:       "InvalidChecksum",
			input:      `{"checksum":{"function":"SHA256","message":"db27fe860c96f269f7838525ba8dce0886e0b7753caccc14162195bcdacbf49e","params":{}},"cipher":{"function":"xor","message":"e18afad793ec8dc3263169c07add77515d9f301464a05508d7ecb42ced24ed3a","params":{}},"kdf":{"function":"scrypt","message":"","params":{"dklen":32,"n":262144,"p":8,"r":1,"salt":"ab0c7876052600dd703518d6fc3fe8984592145b591fc8fb5c6d43190334ba19"}}}`,
			passphrase: "testpassword",
			err:        "invalid checksum",
		},
		{
			name:       "BadIV",
			input:      `{"checksum":{"function":"SHA256","message":"cb27fe860c96f269f7838525ba8dce0886e0b7753caccc14162195bcdacbf49e","params":{}},"cipher":{"function":"xor","message":"e18afad793ec8dc3263169c07add77515d9f301464a05508d7ecb42ced24ed3a","params":{}},"kdf":{"function":"scrypt","message":"","params":{"dklen":32,"n":262144,"p":8,"r":1,"salt":"ab0c7876052600dd703518d6fc3fe8984592145b591fc8fb5c6d43190334ba19"}}}`,
			passphrase: "testpassword",
			output:     []byte{0x1b, 0x4b, 0x68, 0x19, 0x26, 0x11, 0xfa, 0xea, 0x20, 0x8f, 0xca, 0x21, 0x62, 0x7b, 0xe9, 0xda, 0xe6, 0xc3, 0xf2, 0x56, 0x4d, 0x42, 0x58, 0x8f, 0xb1, 0x11, 0x9d, 0xae, 0x7c, 0x9f, 0x4b, 0x87},
		},
		{
			name:       "BadIV",
			input:      `{"kdf":{"function":"pbkdf2","params":{"dklen":32,"c":262144,"prf":"hmac-sha256","salt":"d4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3"},"message":""},"checksum":{"function":"sha256","params":{},"message":"18b148af8e52920318084560fd766f9d09587b4915258dec0676cba5b0da09d8"},"cipher":{"function":"aes-128-ctr","params":{"iv":"h64daa3f303d7259501c93d997d84fe6"},"message": "a9249e0ca7315836356e4c7440361ff22b9fe71e2e2ed34fc1eb03976924ed48"}}`,
			passphrase: "testpassword",
			err:        "invalid IV",
		},
		{
			name:       "BadCipher",
			input:      `{"kdf":{"function":"pbkdf2","params":{"dklen":32,"c":262144,"prf":"hmac-sha256","salt":"d4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3"},"message":""},"checksum":{"function":"sha256","params":{},"message":"18b148af8e52920318084560fd766f9d09587b4915258dec0676cba5b0da09d8"},"cipher":{"function":"aes-64-ctr","params":{"iv":"264daa3f303d7259501c93d997d84fe6"},"message": "a9249e0ca7315836356e4c7440361ff22b9fe71e2e2ed34fc1eb03976924ed48"}}`,
			passphrase: "testpassword",
			err:        `unsupported cipher "aes-64-ctr"`,
		},
		{
			name:       "Good",
			input:      `{"checksum":{"function":"SHA256","message":"cb27fe860c96f269f7838525ba8dce0886e0b7753caccc14162195bcdacbf49e","params":{}},"cipher":{"function":"xor","message":"e18afad793ec8dc3263169c07add77515d9f301464a05508d7ecb42ced24ed3a","params":{}},"kdf":{"function":"scrypt","message":"","params":{"dklen":32,"n":262144,"p":8,"r":1,"salt":"ab0c7876052600dd703518d6fc3fe8984592145b591fc8fb5c6d43190334ba19"}}}`,
			passphrase: "testpassword",
			output:     []byte{0x1b, 0x4b, 0x68, 0x19, 0x26, 0x11, 0xfa, 0xea, 0x20, 0x8f, 0xca, 0x21, 0x62, 0x7b, 0xe9, 0xda, 0xe6, 0xc3, 0xf2, 0x56, 0x4d, 0x42, 0x58, 0x8f, 0xb1, 0x11, 0x9d, 0xae, 0x7c, 0x9f, 0x4b, 0x87},
		},
		{
			name: "Good2",
			input: `{"checksum":{"function":"sha256","message":"a230c7d50dc1e141433559a12cedbe2db2014012b7d5bcda08f399d06ec9bd87","params":{}},"cipher":{"function":"aes-128-ctr","message":"5263382e2ae83dd06020baac533e0173f195be6726f362a683de885c0bdc8e0cec93a411ebc10dfccf8408e23a0072fadc581ab1fcd7a54faae8d2db0680fa76","params":{"iv":"c6437d26eb11abafd373bfb470fd0ad4"}},"kdf":{"function":"scrypt","message":"","params":{"dklen":32,"n":16,"p":8,"r":1,"salt":"20c085c4048f5592cc36bb2a6aa16f0d887f4eb4110849830ceb1eb2dfc0d1be"}}}
`,
			passphrase: "wallet passphrase",
			output: []byte{
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
				0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
				0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encryptor := New()
			input := make(map[string]interface{})
			err := json.Unmarshal([]byte(test.input), &input)
			require.Nil(t, err)
			output, err := encryptor.Decrypt(input, test.passphrase)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.Nil(t, err)
				require.Equal(t, test.output, output)
			}
		})
	}
}

func TestDecryptBadInput(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
		err   string
	}{
		{
			name: "Nil",
			err:  "no data supplied",
		},
		{
			name:  "Empty",
			input: map[string]interface{}{},
			err:   "no checksum",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encryptor := New()
			_, err := encryptor.Decrypt(test.input, "irrelevant")
			require.EqualError(t, err, test.err)
		})
	}
}
