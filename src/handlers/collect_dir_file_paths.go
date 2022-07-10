package handlers

import (
	"fmt"
	"github.com/israeleis/findUsages/src/lib/flags"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

func CollectDirFiles(filePathCh chan<- string, dirPath string, fileTypes []string, errorsCh chan error) []string {

	var anotheDirectories []string

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		errorsCh <- errors.Wrap(err, "error found when trying to read dir files")
	}

	for _, file := range files {
		absPath := fmt.Sprintf("%s/%s", dirPath, file.Name())
		if file.IsDir() {
			anotheDirectories = append(anotheDirectories, absPath)
		} else if fileTypeAllowed(fileTypes, absPath) {
			//go func() {
			filePathCh <- absPath
			//}()

		}
	}
	return anotheDirectories
}

func fileTypeAllowed(fileTypes flags.ArrayFlags, path string) bool {
	lastDotInd := strings.LastIndex(path, ".")
	if lastDotInd == -1 {
		return false
	}
	ext := path[lastDotInd+1:]
	for _, fileType := range fileTypes {
		if fileType == ext {
			return true
		}
	}
	return false
}
