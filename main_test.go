package main

import (
	"fmt"
	"testing"
)

/*
Get duplicate files success.
Total number of detected files: 69
Total number of detected files groups: 4
*/

func TestGetDuplicateFiles(t *testing.T) {
	fileDir := ""
	files, filesNum, err := FindFilesInSameSize(fileDir)
	if err != nil {
		fmt.Println("FindFilesInSameSize failed:")
		t.Error(err)
		return
	}

	allDuplicateFiles, err := FindDuplicateFiles(files)
	if err != nil {
		fmt.Println("FindDuplicateFiles failed:")
		t.Error(err)
		return
	}

	err = WriteResultFile(fileDir, allDuplicateFiles)
	if err != nil {
		fmt.Printf("WriteResultFile failed, err: %s\n", err.Error())
		return
	}

	fmt.Printf("Get duplicate files success.\nTotal number of detected files: %d\nTotal number of detected files groups: %d\n", filesNum, len(allDuplicateFiles))
}
