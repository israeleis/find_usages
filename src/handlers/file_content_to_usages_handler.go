package handlers

import (
	"github.com/israeleis/findUsages/src/models"
	"strings"
	"sync"
)

func FileContentToUsagesHandler(fileContent models.FileContent, values []string, usagesFoundCh chan<- models.UsageResult, wg *sync.WaitGroup, errorCh chan<- error) {
	//for content := range fileContentCh {
	handleFile(fileContent, values, usagesFoundCh, wg)
	//}
}

func handleFile(content models.FileContent, values []string, usagesCh chan<- models.UsageResult, wg *sync.WaitGroup) {

	//wg.Add(1)
	//defer wg.Done()

	//log.WithFields(
	//	log.Fields{
	//		"filePath": content.Path,
	//	}).Info("new file content search for usages")

	contentStr := string(content.Content)
	for _, value := range values {
		if strings.Contains(contentStr, value) {
			lines := strings.Split(contentStr, "\n")
			for i, line := range lines {
				if strings.Contains(line, value) {
					wg.Add(1)
					//go func() {

					usagesCh <- models.UsageResult{
						FindValue: value,
						Usage: models.Usage{
							FilePath:   content.Path,
							Line:       line,
							LineNumber: i,
							Index:      strings.Index(line, value),
						},
					}
					//}()
				}
			}

		}
	}
}
