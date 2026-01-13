package doctor

import (
	"fmt"
	"log"
	"path/filepath"

	"dinhphu28.com/dictionary/internal/config"
	"dinhphu28.com/dictionary/internal/database"
	"dinhphu28.com/dictionary/internal/lookup"
	"dinhphu28.com/dictionary/internal/setup"
)

func checkLookup() {
	paths := setup.DefaultPaths()
	configPath := filepath.Join(paths.ConfigDir, "config.json")
	resourcesPath := filepath.Join(paths.DataDir, "resources")
	if err := config.LoadConfig(configPath); err != nil {
		log.Fatal("failed to load config:", err)
	}
	if err := database.LoadDictionaries(resourcesPath); err != nil {
		log.Fatal("failed to load dictionaries:", err)
	}

	globalConfig := config.GetGlobalConfig()
	dictionaries := database.GetDictionaries()

	if len(dictionaries) == 0 {
		fmt.Println("⚠ Lookup skipped (no dictionaries)")
		return
	}

	dictionaryLookup := lookup.NewDictionaryLookup(dictionaries, globalConfig)
	engine := lookup.NewApproximateLookup(*dictionaryLookup)
	result, err := engine.LookupWithSuggestion("hello")
	if err != nil {
		fmt.Println("✖ Lookup failed:", err)
		return
	}

	if len(result.LookupResults) == 0 {
		fmt.Println("⚠ Lookup returned no results")
		return
	}

	fmt.Printf(
		"✔ Lookup test passed (\"hello\" → %s)\n",
		result.LookupResults[0].Dictionary,
	)
}
