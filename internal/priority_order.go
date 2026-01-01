package internal

import (
	"sort"

	"dinhphu28.com/dictionary/internal/config"
	"dinhphu28.com/dictionary/internal/database"
)

func ApplyPriorityOrder(globalConfig config.GlobalConfig, dictionaries []database.Dictionary) {
	order := map[string]int{}
	for i, id := range globalConfig.Priority {
		order[id] = i
	}

	sort.Slice(dictionaries, func(i, j int) bool {
		return order[dictionaries[i].Manifest.ID] < order[dictionaries[j].Manifest.ID]
	})
}
