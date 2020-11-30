package main

import (
	"fmt"
	"testing"
)

/*
Get duplicate files success.
Total number of detected files: 12
All duplicate files: []main.DuplicateFiles{main.DuplicateFiles{size:16, num:2, files:[]string{"C:\\Users\\Administrator\\Desktop\\getDuplicateFiles\\example\\a", "C:\\Users\\Administrator\\Desktop\\getDuplicateFiles\\example\\b"}}}
*/

func TestGetDuplicateFiles(t *testing.T) {
	files, filesNum, err := FindFilesInSameSize("C:\\Users\\Administrator\\Desktop\\getDuplicateFiles")
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

	fmt.Printf("Get duplicate files success.\nTotal number of detected files: %d\nAll duplicate files: %#v\n", filesNum, allDuplicateFiles)
}

