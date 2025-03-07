package main

import (
	"fmt"

	"github.com/charmbracelet/log"

	fs "github.com/nanoDFS/client-sdk/filesystem"
)

func main() {
	fileId, userId, err := fs.NewFileSystem().Upload("./test.txt")
	if err != nil {
		log.Errorf("got error got this: %v", err)
	}
	fmt.Printf("File id %s \n", fileId)

	err = fs.NewFileSystem().Download(fileId, userId, "./temp")
	if err != nil {
		log.Errorf("got error got this: %v", err)
	}

	err = fs.NewFileSystem().Delete(fileId, userId)
	if err != nil {
		log.Errorf("got error got this: %v", err)
	}

	//select {}
}
