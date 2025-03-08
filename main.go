package main

import (
	"fmt"

	"github.com/charmbracelet/log"

	"github.com/nanoDFS/client-sdk/crypto"
	fs "github.com/nanoDFS/client-sdk/filesystem"
)

func main() {
	key := crypto.DefaultCryptoKey()

	fileId, userId, err := fs.NewFileSystem().Upload(key, "./test.mp4")
	if err != nil {
		log.Errorf("got error got this: %v", err)
	}
	fmt.Printf("File id %s \n", fileId)

	err = fs.NewFileSystem().Download(fileId, userId, key, "./temp")
	if err != nil {
		log.Errorf("got error got this: %v", err)
	}

	err = fs.NewFileSystem().Delete(fileId, userId)
	if err != nil {
		log.Errorf("got error got this: %v", err)
	}

	//select {}
}
