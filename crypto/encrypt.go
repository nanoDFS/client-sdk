package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

type Encryptor struct {
}

func NewEncryptor() *Encryptor {
	return &Encryptor{}
}

func (t *Encryptor) Encrypt(plaintext, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nil
}
