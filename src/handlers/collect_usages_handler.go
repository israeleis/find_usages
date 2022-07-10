package handlers

import (
	"github.com/israeleis/findUsages/src/models"
	"sync"
)

func CollectResults(usageResult models.UsageResult, resultsStorage map[string][]models.Usage, wg *sync.WaitGroup) {
	handleUsageStore(resultsStorage, usageResult, wg)
}

func handleUsageStore(resultsStorage map[string][]models.Usage, resultUsage models.UsageResult, wg *sync.WaitGroup) {
	defer wg.Done()

	findUsagesStored, exists := resultsStorage[resultUsage.FindValue]
	if !exists {
		resultsStorage[resultUsage.FindValue] = []models.Usage{resultUsage.Usage}
	} else {
		resultsStorage[resultUsage.FindValue] = append(findUsagesStored, resultUsage.Usage)
	}
}
