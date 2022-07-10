package main

import (
	"flag"
	"fmt"
	"github.com/israeleis/findUsages/src/handlers"
	"github.com/israeleis/findUsages/src/lib/flags"
	"github.com/israeleis/findUsages/src/models"
	log "github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	var includeTypes flags.ArrayFlags

	findValuesFile := flag.String("values_file", "", "File path to take search values, Each line will contains different search value")
	dirPath := flag.String("dir", "", "Directory to search in")
	flag.Var(&includeTypes, "include", "Include file type")
	flag.Var(&includeTypes, "i", "Include file type(Shortener)")

	flag.Parse()

	includeTypesLower := funk.Map(includeTypes, func(s string) string {
		return strings.ToLower(s)
	}).([]string)

	if dirPath == nil || *dirPath == "" {
		panic("parameter 'dir' is missing")
	}

	findValuesFileAbsolutePath, err := filepath.Abs(*findValuesFile)
	if err != nil {
		panic("parameter 'dir' is missing" + err.Error())
	}
	findValues := extractFindValues(findValuesFileAbsolutePath)

	channels := createChannels()

	go handleErrors(channels.Errors)
	resultsCache := initCache(findValues)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	go func() {
		channels.Directories <- *dirPath
	}()

	go handlers.HandleChannels(channels, &wg, findValues, includeTypesLower, resultsCache, 1)
	//registerHandlers(channels, &wg, findValues, includeTypesLower, resultsCache)

	wg.Wait()
	time.Sleep(5 * time.Second)

	defer channels.Close()

	res := collectResults(resultsCache)
	fmt.Println(res)
}

func collectResults(storage *map[string]*models.UsagesResults) string {
	var resLines []string

	aa := (*storage)[".GetServiceLastAuditLog("]
	println(fmt.Sprintf("%v", aa))
	for _, usagesResults := range *storage {
		resLines = append(resLines, usagesResults.String())
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
