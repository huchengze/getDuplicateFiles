package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
)

type DuplicateFiles struct {
	size  int64
	num   int
	files []string
}

func main() {
	fileDir := ""
	if len(os.Args) > 1 {
		fileDir = os.Args[1]
	}

	files, filesNum, err := FindFilesInSameSize(fileDir)
	if err != nil {
		fmt.Printf("FindFilesInSameSize failed, err: %s\n", err.Error())
		return
	}

	allDuplicateFiles, err := FindDuplicateFiles(files)
	if err != nil {
		fmt.Printf("FindDuplicateFiles failed, err: %s\n", err.Error())
		return
	}

	err = WriteResultFile(fileDir, allDuplicateFiles)
	if err != nil {
		fmt.Printf("WriteResultFile failed, err: %s\n", err.Error())
		return
	}

	fmt.Printf("Get duplicate files success.\nTotal number of detected files: %d\nTotal number of detected files groups: %d\n", filesNum, len(allDuplicateFiles))
}

func FindFilesInSameSize(fileDir string) (files map[int64][]string, filesNum int, err error) {
	files = make(map[int64][]string)
	if fileDir == "" {
		fileDir, _ = os.Getwd()
	}

	err = filepath.Walk(fileDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if strings.Contains(err.Error(), "Access is denied") {
				return nil
			}
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
	maxWorkerNum := runtime.NumCPU() - 2
	var wg sync.WaitGroup
	ch1 := make(chan DuplicateFiles, maxWorkerNum)
	ch2 := make(chan DuplicateFiles, len(files))
	doneCh := make(chan bool, maxWorkerNum)

	go func() {
		for size, sameSizeFiles := range files {
			ch1 <- DuplicateFiles{size: size, files: sameSizeFiles}
		}
		for i := 0; i < maxWorkerNum; i++ {
			doneCh <- true
		}
	}()

	go func(allDuplicateFiles *[]DuplicateFiles) {
		for {
			select {
			case duplicateFiles := <-ch2:
				if duplicateFiles.num > 1 {
					*allDuplicateFiles = append(*allDuplicateFiles, duplicateFiles)
				}
			default:
			}
		}
	}(&allDuplicateFiles)

	for i := 0; i < maxWorkerNum; i++ {
		wg.Add(1)
		go worker(ch1, ch2, doneCh, &wg)
	}

	wg.Wait()
	return
}

func worker(ch1, ch2 chan DuplicateFiles, doneCh chan bool, wg *sync.WaitGroup) {
	for {
		select {
		case sSizeFiles := <-ch1:
			filesMap := make(map[string][]string)
			for _, file := range sSizeFiles.files {
				m, err := GetFileMd5(file)
				if err != nil {
					continue
				}
				if _, ok := filesMap[m]; !ok {
					filesMap[m] = []string{file}
				} else {
					filesMap[m] = append(filesMap[m], file)
				}
			}
			for _, files := range filesMap {
				if len(files) > 1 {
					ch2 <- DuplicateFiles{
						size:  sSizeFiles.size,
						num:   len(files),
						files: files,
					}
				}
			}
		case <-doneCh:
			wg.Done()
		default:
		}
	}
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

func WriteResultFile(fileDir string, allDuplicateFiles []DuplicateFiles) error {
	if fileDir == "" {
		fileDir, _ = os.Getwd()
	}
	sort.Slice(allDuplicateFiles, func(i, j int) bool {
		return allDuplicateFiles[i].size > allDuplicateFiles[j].size
	})

	result := fmt.Sprintf("Get duplicate files in %s.\nTotal number of detected files groups: %d.\n", fileDir, len(allDuplicateFiles))
	for _, allDuplicateFile := range allDuplicateFiles {
		result += fmt.Sprintf("\nsize: %d\nnum: %d\nfiles: %s\n", allDuplicateFile.size, allDuplicateFile.num, strings.Join(allDuplicateFile.files, ","))
	}

	return ioutil.WriteFile("result.txt", []byte(result), 0644)
}
