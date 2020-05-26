package encryptors

import (
	"encoding/json"
)

type PlainTextEncryptor struct {

}

type plainTextEncryption struct {
	original string
}

func NewPlainTextEncryptor() *PlainTextEncryptor {
	return &PlainTextEncryptor{}
}

// Name() provides the name of the encryptor
func (encryptor *PlainTextEncryptor) Name() string {
	return "Plain Text Encryptor"
}

// Version() provides the version of the encryptor
func (encryptor *PlainTextEncryptor) Version() uint {
	return 1
}

// Encrypt encrypts a byte array with its encryption mechanism and key
func (encryptor *PlainTextEncryptor) Encrypt(data []byte, key []byte) (map[string]interface{}, error) {
	output := &plainTextEncryption{original:string(data)}
	bytes, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Decrypt encrypts a byte array with its encryption mechanism and key
func (encryptor *PlainTextEncryptor) Decrypt(data map[string]interface{}, key []byte) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	ks := &plainTextEncryption{}
	err = json.Unmarshal(b, &ks)
	if err != nil {
		return nil, err
	}
	return []byte(ks.original),nil
}
