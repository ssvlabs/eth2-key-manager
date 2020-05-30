package encryptors

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestEncryption(t *testing.T) {
	encryptor := NewPlainTextEncryptor()

	tests := []struct {
		name string
		data []byte
		expectedEncrypted []byte
		password []byte
	}{
		{
			name: "simple password",
			data:[]byte("data"),
			expectedEncrypted:[]byte("ZGF0YQ=="),
			password:[]byte("password"),
		},
		{
			name: "empty password",
			data:[]byte("data"),
			expectedEncrypted:[]byte("ZGF0YQ=="),
			password:[]byte(""),
		},
		{
			name: "seed",
			data:[]byte("de66c3680d21a10fa5079b62e89b983165f4a4897cdabeaa34e195c56e14c86b"),
			expectedEncrypted:[]byte("ZGU2NmMzNjgwZDIxYTEwZmE1MDc5YjYyZTg5Yjk4MzE2NWY0YTQ4OTdjZGFiZWFhMzRlMTk1YzU2ZTE0Yzg2Yg=="),
			password:[]byte("1234"),
		},
	}

	for _,test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encrypted,err := encryptor.Encrypt(test.data,test.password)
			if err != nil {
				t.Error(err)
				return
			}

			if strings.Compare(fmt.Sprintf("%s",encrypted["original"]),string(test.expectedEncrypted)) != 0 {
				t.Error(fmt.Errorf("encryptd data not matching expected"))
				return
			}

			decrypted,err := encryptor.Decrypt(encrypted,test.password)
			if err != nil {
				t.Error(err)
				return
			}
			if bytes.Compare(decrypted,test.data) != 0 {
				t.Error(fmt.Errorf("decrypted data not matching original"))
			}
		})
	}
}
