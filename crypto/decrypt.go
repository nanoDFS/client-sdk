package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

type Decryptor struct {
}

func NewDecryptor() *Decryptor {
	return &Decryptor{}
}

func (t *Decryptor) Decrypt(ciphertext, nonce, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
