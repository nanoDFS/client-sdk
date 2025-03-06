package test

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/log"
	fs "github.com/nanoDFS/client-sdk/filesystem"
)

func TestUpload(t *testing.T) {
	fileId, userId, err := fs.NewFileSystem().Upload("./test_file.txt")
	fmt.Print(fileId, userId, userId)
	if err != nil {
		t.Errorf("got error got this: %v", err)
	}
	log.Printf("%s %s", fileId, userId)
}
