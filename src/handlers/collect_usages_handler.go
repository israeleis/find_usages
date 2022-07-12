package handlers

import (
	"fmt"
	"github.com/israeleis/findUsages/src/models"
	"sync"
)

func CollectResults(usageResult models.UsageResult, resultsStorage *map[string]*models.UsagesResults, wg *sync.WaitGroup) {
	defer wg.Done()
	handleUsageStore(resultsStorage, usageResult, wg)
}

func handleUsageStore(resultsStorage *map[string]*models.UsagesResults, resultUsage models.UsageResult, wg *sync.WaitGroup) {

	findUsagesStored, exists := (*resultsStorage)[resultUsage.FindValue]
	if !exists {
		panic(fmt.Sprintf("%s is missing in result storage", resultUsage.FindValue))
	}
	findUsagesStored.AddResult(resultUsage.Usage)
}
