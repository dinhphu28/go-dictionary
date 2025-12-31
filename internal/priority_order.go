package internal

import "sort"

func ApplyPriorityOrder(globalConfig GlobalConfig, dictionaries []Dictionary) {
	order := map[string]int{}
	for i, id := range globalConfig.Priority {
		order[id] = i
	}

	sort.Slice(dictionaries, func(i, j int) bool {
		return order[dictionaries[i].Manifest.ID] < order[dictionaries[j].Manifest.ID]
	})
}
