package test

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/nanoDFS/client-sdk/crypto"
	fs "github.com/nanoDFS/client-sdk/filesystem"
)

func TestUpload(t *testing.T) {
	key := crypto.DefaultCryptoKey()
	fileId, userId, err := fs.NewFileSystem().Upload(key, "./test_file.txt")
	fmt.Print(fileId, userId, userId)
	if err != nil {
		t.Errorf("got error got this: %v", err)
	}
	log.Printf("%s %s", fileId, userId)
}
