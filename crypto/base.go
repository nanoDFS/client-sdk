package crypto

import (
	"crypto/rand"

	"github.com/charmbracelet/log"
)

type CryptoKey struct {
	Key   []byte
	Nonce []byte
}

func NewCryptoKey(key string, nonce string) CryptoKey {
	return CryptoKey{
		Key:   []byte(key),
		Nonce: []byte(nonce)[:12],
	}
}

func DefaultCryptoKey() CryptoKey {
	return CryptoKey{
		Key:   GenerateRandomBytes(32),
		Nonce: GenerateRandomBytes(12),
	}
}

func GenerateRandomBytes(length int) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Warn("random bytes generator error")
	}
	return bytes
}
