package keystorev4

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

const (
	name    = "keystore"
	version = 4

	// Scrypt parameters
	scryptN      = 262144
	scryptr      = 8
	scryptp      = 1
	scryptKeyLen = 32

	// PBKDF2 parameters
	pbkdf2KeyLen = 32
	pbkdf2c      = 262144
	pbkdf2PRF    = "hmac-sha256"
)

// Encryptor is an encryptor that follows the Ethereum keystore V4 specification.
type Encryptor struct {
	cipher string
}

type ksKDFParams struct {
	// Shared parameters
	Salt  string `json:"salt"`
	DKLen int    `json:"dklen"`
	// Scrypt-specific parameters
	N int `json:"n,omitempty"`
	P int `json:"p,omitempty"`
	R int `json:"r,omitempty"`
	// PBKDF2-specific parameters
	C   int    `json:"c,omitempty"`
	PRF string `json:"prf,omitempty"`
}

type ksKDF struct {
	Function string       `json:"function"`
	Params   *ksKDFParams `json:"params"`
	Message  string       `json:"message"`
}

type ksChecksum struct {
	Function string                 `json:"function"`
	Params   map[string]interface{} `json:"params"`
	Message  string                 `json:"message"`
}

type ksCipherParams struct {
	// AES-128-CTR-specific parameters
	IV string `json:"iv,omitempty"`
}

type ksCipher struct {
	Function string          `json:"function"`
	Params   *ksCipherParams `json:"params"`
	Message  string          `json:"message"`
}

type keystoreV4 struct {
	KDF      *ksKDF      `json:"kdf"`
	Checksum *ksChecksum `json:"checksum"`
	Cipher   *ksCipher   `json:"cipher"`
}

// options are the options for the keystore encryptor.
type options struct {
	cipher string
}

// Option gives options to New
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

// WithCipher sets the cipher for the encryptor.
func WithCipher(cipher string) Option {
	return optionFunc(func(o *options) {
		o.cipher = cipher
	})
}

// New creates a new keystore V4 encryptor.
// This takes the following options:
// - cipher: the cipher to use when encrypting the secret, can be either "pbkdf2" (default) or "scrypt"
func New(opts ...Option) *Encryptor {
	options := options{
		cipher: "pbkdf2",
	}
	for _, o := range opts {
		o.apply(&options)
	}

	return &Encryptor{
		cipher: options.cipher,
	}
}

// Name returns the name of this encryptor
func (e *Encryptor) Name() string {
	return name
}

// Version returns the version of this encryptor
func (e *Encryptor) Version() uint {
	return version
}

// Encrypt encrypts data.
func (e *Encryptor) Encrypt(secret []byte, passphrase string) (map[string]interface{}, error) {
	if secret == nil {
		return nil, errors.New("no secret")
	}

	// Random salt
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	normedPassphrase := []byte(normPassphrase(passphrase))
	// Create the decryption key
	var decryptionKey []byte
	var err error
	switch e.cipher {
	case "scrypt":
		decryptionKey, err = scrypt.Key(normedPassphrase, salt, scryptN, scryptr, scryptp, scryptKeyLen)
	case "pbkdf2":
		decryptionKey = pbkdf2.Key(normedPassphrase, salt, pbkdf2c, pbkdf2KeyLen, sha256.New)
	default:
		return nil, errors.Errorf("unknown cipher %q", e.cipher)
	}
	if err != nil {
		return nil, err
	}

	// Generate the cipher message
	cipherMsg := make([]byte, len(secret))
	aesCipher, err := aes.NewCipher(decryptionKey[:16])
	if err != nil {
		return nil, err
	}
	// Random IV
	iv := make([]byte, 16)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesCipher, iv)
	stream.XORKeyStream(cipherMsg, secret)

	// Generate the checksum
	h := sha256.New()
	if _, err := h.Write(decryptionKey[16:32]); err != nil {
		return nil, err
	}
	if _, err := h.Write(cipherMsg); err != nil {
		return nil, err
	}
	checksumMsg := h.Sum(nil)

	var kdf *ksKDF
	switch e.cipher {
	case "scrypt":
		kdf = &ksKDF{
			Function: "scrypt",
			Params: &ksKDFParams{
				DKLen: scryptKeyLen,
				N:     scryptN,
				P:     scryptp,
				R:     scryptr,
				Salt:  hex.EncodeToString(salt),
			},
			Message: "",
		}
	case "pbkdf2":
		kdf = &ksKDF{
			Function: "pbkdf2",
			Params: &ksKDFParams{
				DKLen: pbkdf2KeyLen,
				C:     pbkdf2c,
				PRF:   pbkdf2PRF,
				Salt:  hex.EncodeToString(salt),
			},
			Message: "",
		}
	}

	// Build the output
	output := &keystoreV4{
		KDF: kdf,
		Checksum: &ksChecksum{
			Function: "sha256",
			Params:   make(map[string]interface{}),
			Message:  hex.EncodeToString(checksumMsg),
		},
		Cipher: &ksCipher{
			Function: "aes-128-ctr",
			Params: &ksCipherParams{
				IV: hex.EncodeToString(iv),
			},
			Message: hex.EncodeToString(cipherMsg),
		},
	}

	// We need to return a generic map; go to JSON and back to obtain it
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
