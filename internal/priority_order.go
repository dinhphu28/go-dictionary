package internal

import (
	"sort"

	"dinhphu28.com/dictionary/internal/config"
)

func ApplyPriorityOrder(globalConfig config.GlobalConfig, dictionaries []Dictionary) {
	order := map[string]int{}
	for i, id := range globalConfig.Priority {
		order[id] = i
	}

	sort.Slice(dictionaries, func(i, j int) bool {
		return order[dictionaries[i].Manifest.ID] < order[dictionaries[j].Manifest.ID]
	})
}
