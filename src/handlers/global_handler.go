package handlers

import (
	"github.com/israeleis/findUsages/src/models"
	"sync"
)

func HandleChannels(channels models.FlowChannels, wg *sync.WaitGroup, matchers []models.Matcher, regex *bool, fileFilters models.FileNameFilters, resultsStorage *map[string]*models.UsagesResults, concurrencyLimit int) {

	go collectDirectories(&channels, fileFilters, wg)
	sem := make(chan struct{}, concurrencyLimit)

	for {
		sem <- struct{}{}
		//println("runtime.NumGoroutine(): ", runtime.NumGoroutine())
		func() {
			//defer func() { <-sem }()
			//channelsListener(channels, wg, findValues, includeFile, resultsStorage)
			select {
			//case dirPath := <-channels.Directories:
			//	go func() {
			//		defer func() {
			//			wg.Done()
			//		}()
			//		nextDirectories := CollectDirFiles(channels.FilePath, dirPath, includeFile, channels.Errors)
			//		for _, dir := range nextDirectories {
			//			go func(dir string) {
			//				wg.Add(1)
			//				channels.Directories <- dir
			//			}(dir)
			//		}
			//		<-sem
			//	}()
			case filePath := <-channels.FilePath:
				go func() {
					FilePathToFileContentChannelHandler(channels.FileContent, filePath, wg, channels.Errors)

					<-sem
				}()
			case fileContent := <-channels.FileContent:
				go func() {
					FileContentToUsagesHandler(fileContent, matchers, regex, channels.UsageResult, wg, channels.Errors)

					<-sem
				}()
			case usageResult := <-channels.UsageResult:
				go func() {
					CollectResults(usageResult, resultsStorage, wg)

					<-sem
				}()
			default:
				<-sem
				//println("waiting for message")
			}
		}()
	}
}

func collectDirectories(channels *models.FlowChannels, fileFilters models.FileNameFilters, wg *sync.WaitGroup) {
	for dirPath := range channels.Directories {
		func() {
			defer wg.Done()

			nextDirectories := CollectDirFiles(channels.FilePath, dirPath, fileFilters, wg, channels.Errors)
			for _, dir := range nextDirectories {
				go func(dir string) {
					//wg.Add(1)
					channels.Directories <- dir
				}(dir)
			}
		}()
	}
}

//func channelsListener(channels models.FlowChannels, wg *sync.WaitGroup, findValues []string, fileTypes []string, resultsStorage map[string][]models.Usage) {
//	select {
//	case dirPath := <-channels.Directories:
//		func() {
//			defer func() {
//				wg.Done()
//			}()
//			nextDirectories := CollectDirFiles(channels.FilePath, dirPath, fileTypes, channels.Errors)
//			for _, dir := range nextDirectories {
//				wg.Add(1)
//				channels.Directories <- dir
//			}
//		}()
//	case filePath := <-channels.FilePath:
//		FilePathToFileContentChannelHandler(channels.FileContent, filePath, channels.Errors)
//	case fileContent := <-channels.FileContent:
//		FileContentToUsagesHandler(fileContent, findValues, channels.UsageResult, wg, channels.Errors)
//	case usageResult := <-channels.UsageResult:
//		CollectResults(usageResult, resultsStorage, wg)
//	default:
//		println("waiting for message")
//	}
//}
