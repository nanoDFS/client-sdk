package filesystem

import (
	"github.com/charmbracelet/log"
	"github.com/nanoDFS/client-sdk/filesystem/delete"
)

func (t *FileSystem) Delete(fileId string, userId string) error {
	err := delete.NewDeletor(fileId, userId).Delete()
	if err != nil {
		log.Errorf("failed to delete file: %v", err)
	}
	return err
}
