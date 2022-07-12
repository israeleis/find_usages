package handlers

import (
	"github.com/israeleis/findUsages/src/models"
	"os"
	"sync"
)

func FilePathToFileContentChannelHandler(fileContentCh chan<- models.FileContent, filePath string, wg *sync.WaitGroup, errorCh chan<- error) {
	defer wg.Done()
	//log.WithFields(
	//	log.Fields{
	//		"filePath": filePath,
	//	}).Info("new file path has been found for search")

	readFileContent(filePath, errorCh, fileContentCh, wg)

}

func readFileContent(filePath string, errorCh chan<- error, fileContentCh chan<- models.FileContent, wg *sync.WaitGroup) {

	content, err := os.ReadFile(filePath)
	if err != nil {
		errorCh <- err
		return
	}
	//go func() {
	wg.Add(1)
	fileContentCh <- models.FileContent{
		Path:    filePath,
		Content: content,
	}
	//}()
}
