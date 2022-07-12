package main

import (
	"flag"
	"fmt"
	"github.com/israeleis/findUsages/src/handlers"
	"github.com/israeleis/findUsages/src/lib/flags"
	"github.com/israeleis/findUsages/src/models"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	var includeFileNameFilters flags.ArrayFlags
	var excludeFileNameFilters flags.ArrayFlags

	findValuesFile := flag.String("values_file", "", "File path to take search values, Each line will contains different search value")
	isRegex := flag.Bool("regex", false, "search with Regex (Default 'false')")
	dirPath := flag.String("dir", "", "Directory to search in")
	flag.Var(&includeFileNameFilters, "include", "Include file type")
	flag.Var(&includeFileNameFilters, "i", "Include file type(Shortener)")
	flag.Var(&excludeFileNameFilters, "exclude", "Exclude file type")
	flag.Var(&excludeFileNameFilters, "ex", "Exclude file type(Shortener)")
	maxUsages := flag.Int("max_usages", 1000, "max usages to be displayed")

	flag.Parse()

	if dirPath == nil || *dirPath == "" {
		panic("parameter 'dir' is missing")
	}

	findValuesFileAbsolutePath, err := filepath.Abs(*findValuesFile)
	if err != nil {
		panic("parameter 'dir' is missing" + err.Error())
	}
	findValues := extractFindValues(findValuesFileAbsolutePath)

	filesFilter := models.FileNameFilters{
		Include: models.CreateFileNameFilters(includeFileNameFilters),
		Exclude: models.CreateFileNameFilters(excludeFileNameFilters),
	}
	matchers := models.CreateMatchers(findValues, *isRegex)
	channels := createChannels()

	go handleErrors(channels.Errors)
	resultsCache := initCache(findValues)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	go func() {
		wg.Add(1)
		channels.Directories <- *dirPath
	}()

	go handlers.HandleChannels(channels, &wg, matchers, isRegex, filesFilter, resultsCache, 50)
	//registerHandlers(channels, &wg, findValues, includeTypesLower, resultsCache)

	time.Sleep(60 * time.Second)

	wg.Wait()
	println("sleeping")
	time.Sleep(3 * time.Second)

	defer channels.Close()

	res := collectResults(resultsCache, *maxUsages)
	fmt.Println(res)
}

func collectResults(storage *map[string]*models.UsagesResults, maxUsages int) string {
	var resLines []string

	for _, usagesResults := range *storage {
		if len(usagesResults.Usages) <= maxUsages {
			resLines = append(resLines, usagesResults.String())
		}
	}

	return strings.Join(resLines, "\n------------------------------------------------\n")
}

func initCache(findValues []string) *map[string]*models.UsagesResults {
	resultsCache := make(map[string]*models.UsagesResults)
	for _, value := range findValues {
		resultsCache[value] = &models.UsagesResults{
			FindValue: value,
			Usages:    []models.Usage{},
		}
	}
	return &resultsCache
}

func handleErrors(errorsCh <-chan error) {
	for err := range errorsCh {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("error found in flow")
	}
}

func extractFindValues(file string) []string {

	contentBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	content := string(contentBytes)
	return strings.Split(content, "\n")
}

func createChannels() models.FlowChannels {

	return models.FlowChannels{
		Errors:      make(chan error),
		Directories: make(chan string, 10000),
		FilePath:    make(chan string, 10000),
		FileContent: make(chan models.FileContent, 10000),
		UsageResult: make(chan models.UsageResult, 10000),
	}
}

//func registerHandlers(channels models.FlowChannels, wg *sync.WaitGroup, findValues []string, includeTypes []string, resultsStorage map[string][]models.Usage) {
//	go handlers.CollectDirFiles(channels.Directories, channels.FilePath, includeTypes, channels.Errors)
//	go handlers.FilePathToFileContentChannelHandler(channels.FilePath, channels.FileContent, wg, channels.Errors)
//	go handlers.FileContentToUsagesHandler(channels.FileContent, findValues, channels.UsageResult, wg, channels.Errors)
//	go handlers.CollectResults(channels.UsageResult, resultsStorage, wg)
//}
