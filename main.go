package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type DuplicateFiles struct {
	size  int64
	num   int
	files []string
}

func main() {
	files, filesNum, err := FindFilesInSameSize(os.Args[1])
	if err != nil {
		fmt.Printf("FindFilesInSameSize failed, err: %s\n", err.Error())
		return
	}

	allDuplicateFiles, err := FindDuplicateFiles(files)
	if err != nil {
		fmt.Printf("FindDuplicateFiles failed, err: %s\n", err.Error())
		return
	}

	fmt.Printf("Get duplicate files success.\nAll duplicate files: %#v\nTotal number of detected files: %d\n", allDuplicateFiles, filesNum)
}

func FindFilesInSameSize(filedir string) (files map[int64][]string, filesNum int, err error) {
	files = make(map[int64][]string)

	err = filepath.Walk(filedir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filesNum++
			s := info.Size()
			if _, ok := files[s]; !ok {
				files[s] = []string{path}
			} else {
				files[s] = append(files[s], path)
			}
		}
		return nil
	})

	for size, sameSizeFiles := range files {
		if len(sameSizeFiles) < 2 {
			delete(files, size)
		}
	}

	return
}

func FindDuplicateFiles(files map[int64][]string) (allDuplicateFiles []DuplicateFiles, err error) {
	for size, sameSizeFiles := range files {
		filesmap := make(map[string][]string)
		for _, file := range sameSizeFiles {
			m, err := GetFileMd5(file)
			if err != nil {
				continue
			}
			if _, ok := filesmap[m]; !ok {
				filesmap[m] = []string{file}
			} else {
				filesmap[m] = append(filesmap[m], file)
			}
		}
		for _, files := range filesmap {
			if len(files) > 1 {
				allDuplicateFiles = append(allDuplicateFiles, DuplicateFiles{
					size:  size,
					num:   len(files),
					files: files,
				})
			}
		}

	}
	return
}

func GetFileMd5(path string) (string, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
