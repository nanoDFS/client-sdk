package filesystem

import (
	"github.com/charmbracelet/log"
	"github.com/nanoDFS/client-sdk/filesystem/download"
)

func (t *FileSystem) Download(fileId string, userId string, filPath string) error {
	err := download.NewDownloader(fileId, userId, filPath).Download()
	if err != nil {
		log.Errorf("failed to download file: %v", err)
	}
	return err
}
