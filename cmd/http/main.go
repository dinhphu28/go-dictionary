package main

import (
	"log"

	"dinhphu28.com/dictionary/internal/api"
	"dinhphu28.com/dictionary/internal/config"
	"dinhphu28.com/dictionary/internal/database"
	"dinhphu28.com/dictionary/internal/lookup"
	"dinhphu28.com/dictionary/internal/startup"
	_ "modernc.org/sqlite"
)

// ---- global list of loaded dictionaries ----

var globalConfig config.GlobalConfig

// ---- Start of code ----
var dictionaries []database.Dictionary

func main() {
	configPath := startup.ResolvePath("config.json")
	resourcesPath := startup.ResolvePath("resources")

	if err := config.LoadConfig(configPath); err != nil {
		log.Fatal("failed to load config:", err)
	}
	globalConfig = config.GetGlobalConfig()

	if err := database.LoadDictionaries(resourcesPath); err != nil {
		log.Fatal("failed to load dictionaries:", err)
	}
	dictionaries = database.GetDictionaries()

	log.Printf("Loaded %d dictionaries\n", len(dictionaries))

	dictionaryLookup := lookup.NewDictionaryLookup(dictionaries, globalConfig)
	approximateLookup := lookup.NewApproximateLookup(*dictionaryLookup)
	lookupHandlerV2 := api.NewLookupHandlerV2(*approximateLookup)
	router := api.NewRouter(*lookupHandlerV2)
	router.StartAPIServer()
}
