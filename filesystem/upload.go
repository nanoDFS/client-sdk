package filesystem

import (
	"github.com/charmbracelet/log"
	"github.com/nanoDFS/client-sdk/crypto"
	"github.com/nanoDFS/client-sdk/filesystem/upload"
)

// Upload returns fileId, userId, err
func (t *FileSystem) Upload(key crypto.CryptoKey, filePath string) (string, string, error) {
	fileId, userId, err := upload.NewUploader(key, filePath).Upload()
	if err != nil {
		log.Errorf("failed to upload file: %v", err)
	}
	return fileId, userId, err
}
