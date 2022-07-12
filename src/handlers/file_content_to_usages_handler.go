package handlers

import (
	"github.com/israeleis/findUsages/src/models"
	"strings"
	"sync"
)

func FileContentToUsagesHandler(fileContent models.FileContent, matchers []models.Matcher, regex *bool, usagesFoundCh chan<- models.UsageResult, wg *sync.WaitGroup, errorCh chan<- error) {
	defer wg.Done()
	//for content := range fileContentCh {
	handleFile(fileContent, matchers, regex, usagesFoundCh, wg)
	//}
}

func handleFile(content models.FileContent, matchers []models.Matcher, regex *bool, usagesCh chan<- models.UsageResult, wg *sync.WaitGroup) {

	//wg.Add(1)
	//defer wg.Done()

	//log.WithFields(
	//	log.Fields{
	//		"filePath": content.Path,
	//	}).Info("new file content search for usages")

	contentStr := string(content.Content)
	for _, matcher := range matchers {
		lines := strings.Split(contentStr, "\n")
		for i, line := range lines {
			usage := matcher.FindUsage(line, content.Path, i)
			if usage == nil {
				continue
			}
			wg.Add(1)
			usagesCh <- models.UsageResult{
				FindValue: matcher.GetValue(),
				Usage:     *usage,
			}
		}
	}
}
