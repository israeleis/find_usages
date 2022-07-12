package handlers

import (
	"fmt"
	"github.com/israeleis/findUsages/src/models"
	"github.com/pkg/errors"
	"io/ioutil"
	"sync"
)

func CollectDirFiles(filePathCh chan<- string, dirPath string, fileFilters models.FileNameFilters, wg *sync.WaitGroup, errorsCh chan error) []string {
	var anotherDirectories []string

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		errorsCh <- errors.Wrap(err, "error found when trying to read dir files")
	}

	for _, file := range files {

		//fn := file.Name()
		//fmt.Println(fn)
		absPath := fmt.Sprintf("%s/%s", dirPath, file.Name())
		if file.IsDir() {
			wg.Add(1)
			anotherDirectories = append(anotherDirectories, absPath)
		} else if fileFilters.IsRelevant(absPath) {
			//go func() {
			wg.Add(1)
			filePathCh <- absPath
			//}()

		}
	}
	return anotherDirectories
}
