package core

// Encryptor is the interface for encrypting and decrypting sensitive information in wallets.
type Encryptor interface {
	// Name() provides the name of the encryptor
	Name() string

	// Version() provides the version of the encryptor
	Version() uint

	// Encrypt encrypts a byte array with its encryption mechanism and key
	Encrypt(data []byte, key string) (map[string]interface{}, error)

	// Decrypt encrypts a byte array with its encryption mechanism and key
	Decrypt(data map[string]interface{}, key string) ([]byte, error)
}
