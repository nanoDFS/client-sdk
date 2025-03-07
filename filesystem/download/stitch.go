package download

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func (t *Downloader) stitchFiles(outputFileDir string, inputFilesDir string) error {
	files := t.listChunkFiles(inputFilesDir)

	err := os.MkdirAll(outputFileDir, os.ModePerm)
	if err != nil {
		return err
	}
	outputFile := path.Join(outputFileDir, "output.bin")
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("could not create output file: %v", err)
	}
	defer output.Close()

	for _, inputFile := range files {
		fullName := path.Join(inputFilesDir, inputFile.name)
		input, err := os.Open(fullName)
		if err != nil {
			return fmt.Errorf("could not open input file %s: %v", inputFile.name, err)
		}
		defer input.Close()
		_, err = io.Copy(output, input)
		if err != nil {
			return fmt.Errorf("could not copy data from %s to output file: %v", inputFile.name, err)
		}
	}
	t.cleanUp(inputFilesDir, "^temp-.*")
	return nil
}

func (t *Downloader) listChunkFiles(dir string) []struct {
	name string
	num  int
} {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var sortedFiles []struct {
		name string
		num  int
	}

	for _, file := range files {
		if !file.IsDir() {
			if strings.HasPrefix(file.Name(), "temp-") {
				numStr := file.Name()[5:]
				num, err := strconv.Atoi(numStr)
				if err != nil {
					log.Println("Error parsing number:", err)
					continue
				}
				sortedFiles = append(sortedFiles, struct {
					name string
					num  int
				}{name: file.Name(), num: num})
			}
		}
	}

	sort.Slice(sortedFiles, func(i, j int) bool {
		return sortedFiles[i].num < sortedFiles[j].num
	})

	return sortedFiles
}

func (t *Downloader) cleanUp(dir, pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex: %v", err)
	}

	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && re.MatchString(d.Name()) {
			return os.Remove(path)
		}
		return nil
	})
}
