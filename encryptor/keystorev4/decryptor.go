package keystorev4

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

// Decrypt decrypts the data provided, returning the secret.
func (e *Encryptor) Decrypt(data map[string]interface{}, passphrase string) ([]byte, error) {
	if data == nil {
		return nil, errors.New("no data supplied")
	}
	// Marshal the map and unmarshal it back in to a keystore format so we can work with it.
	b, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("failed to parse keystore")
	}
	ks := &keystoreV4{}
	err = json.Unmarshal(b, &ks)
	if err != nil {
		return nil, errors.New("failed to parse keystore")
	}

	// Checksum and cipher are required
	if ks.Checksum == nil {
		return nil, errors.New("no checksum")
	}
	if ks.Cipher == nil {
		return nil, errors.New("no cipher")
	}

	normedPassphrase := []byte(normPassphrase(passphrase))
	// Decryption key
	var decryptionKey []byte
	if ks.KDF == nil {
		decryptionKey = normedPassphrase
	} else {
		kdfParams := ks.KDF.Params
		salt, err := hex.DecodeString(kdfParams.Salt)
		if err != nil {
			return nil, errors.New("invalid KDF salt")
		}
		switch ks.KDF.Function {
		case "scrypt":
			decryptionKey, err = scrypt.Key(normedPassphrase, salt, kdfParams.N, kdfParams.R, kdfParams.P, kdfParams.DKLen)
		case "pbkdf2":
			switch kdfParams.PRF {
			case "hmac-sha256":
				decryptionKey = pbkdf2.Key(normedPassphrase, salt, kdfParams.C, kdfParams.DKLen, sha256.New)
			default:
				return nil, errors.Errorf("unsupported PBKDF2 PRF %q", kdfParams.PRF)
			}
		default:
			return nil, errors.Errorf("unsupported KDF %q", ks.KDF.Function)
		}
		if err != nil {
			return nil, errors.New("invalid KDF parameters")
		}
	}

	// Checksum
	if len(decryptionKey) < 32 {
		return nil, errors.New("decryption key must be at least 32 bytes")
	}
	cipherMsg, err := hex.DecodeString(ks.Cipher.Message)
	if err != nil {
		return nil, errors.New("invalid cipher message")
	}
	h := sha256.New()
	if _, err := h.Write(decryptionKey[16:32]); err != nil {
		return nil, err
	}
	if _, err := h.Write(cipherMsg); err != nil {
		return nil, err
	}
	checksum := h.Sum(nil)
	checksumMsg, err := hex.DecodeString(ks.Checksum.Message)
	if err != nil {
		return nil, errors.New("invalid checksum message")
	}
	if !bytes.Equal(checksum, checksumMsg) {
		return nil, errors.New("invalid checksum")
	}

	// Decrypt
	res := make([]byte, len(cipherMsg))
	switch ks.Cipher.Function {
	case "xor":
		for i := range decryptionKey {
			res[i] = decryptionKey[i] ^ cipherMsg[i]
		}
	case "aes-128-ctr":
		aesCipher, err := aes.NewCipher(decryptionKey[:16])
		if err != nil {
			return nil, err
		}
		iv, err := hex.DecodeString(ks.Cipher.Params.IV)
		if err != nil {
			return nil, errors.New("invalid IV")
		}
		stream := cipher.NewCTR(aesCipher, iv)
		stream.XORKeyStream(res, cipherMsg)
	default:
		return nil, errors.Errorf("unsupported cipher %q", ks.Cipher.Function)
	}

	return res, nil
}
