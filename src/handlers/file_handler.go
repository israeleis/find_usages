package handlers

import (
	"github.com/israeleis/findUsages/src/models"
	"os"
)

func FilePathToFileContentChannelHandler(fileContentCh chan<- models.FileContent, filePath string, errorCh chan<- error) {
	//log.WithFields(
	//	log.Fields{
	//		"filePath": filePath,
	//	}).Info("new file path has been found for search")

	readFileContent(filePath, errorCh, fileContentCh)

}

func readFileContent(filePath string, errorCh chan<- error, fileContentCh chan<- models.FileContent) {

	content, err := os.ReadFile(filePath)
	if err != nil {
		errorCh <- err
		return
	}
	//go func() {
	fileContentCh <- models.FileContent{
		Path:    filePath,
		Content: content,
	}
	//}()
}
