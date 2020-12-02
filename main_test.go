package main

import (
	"fmt"
	"testing"
)

/*
Get duplicate files success.
All duplicate files: []main.DuplicateFiles{main.DuplicateFiles{size:9, num:3, files:[]string{"C:\\Users\\huchengze\\Desktop\\getDuplicateFiles\\example\\c", "C:\\Users\\huchengze\\Desktop\\getDuplicateFiles\\example\\e", "C:\\Users\\huchengze\\Desktop\\getDuplicateFiles\\example\\g"}}, main.DuplicateFiles{size:196, num:3, files:[]string{"C:\\Users\\huchengze\\Desktop\\getDuplicateFiles\\.git\\logs\\HEAD", "C:\\Users\\huchengze\\Desktop\\getDuplicateFiles\\.git\\logs\\refs\\heads\\main", "C:\\Users\\huchengze\\Desktop\\getDuplicateFiles\\.git\\logs\\refs\\remotes\\origin\\HEAD"}}, main.DuplicateFiles{size:16, num:2, files:[]string{"C:\\Users\\huchengze\\Desktop\\getDuplicateFiles\\example\\a", "C:\\Users\\huchengze\\Desktop\\getDuplicateFiles\\example\\b"}}, main.DuplicateFiles{size:0, num:2, files:[]string{"C:\\Users\\huchengze\\Desktop\\getDuplicateFiles\\example\\d", "C:\\Users\\huchengze\\Desktop\\getDuplicateFiles\\example\\f"}}}
Total number of detected files: 72
*/

func TestGetDuplicateFiles(t *testing.T) {
	files, filesNum, err := FindFilesInSameSize("C:\\Users\\huchengze\\Desktop\\getDuplicateFiles")
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

	fmt.Printf("Get duplicate files success.\nTotal number of detected files: %d\nTotal number of detected files: %d\n", filesNum, len(allDuplicateFiles))
}
