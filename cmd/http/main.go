package main

import (
	"log"

	"dinhphu28.com/dictionary/internal/api"
	"dinhphu28.com/dictionary/internal/config"
	"dinhphu28.com/dictionary/internal/database"
	"dinhphu28.com/dictionary/internal/lookup"
	_ "modernc.org/sqlite"
)

// ---- global list of loaded dictionaries ----

var globalConfig config.GlobalConfig

// ---- Start of code ----
var dictionaries []database.Dictionary

func main() {
	if err := config.LoadConfig("config.json"); err != nil {
		log.Fatal("failed to load config:", err)
	}
	globalConfig = config.GetGlobalConfig()

	if err := database.LoadDictionaries("resources"); err != nil {
		log.Fatal("failed to load dictionaries:", err)
	}
	dictionaries = database.GetDictionaries()

	log.Printf("Loaded %d dictionaries\n", len(dictionaries))

	approximateLookup := lookup.NewApproximateLookup(dictionaries, globalConfig)
	lookupHandler := api.NewLookupHandler(dictionaries, globalConfig)
	lookupHandlerV2 := api.NewLookupHandlerV2(*approximateLookup)
	router := api.NewRouter(*lookupHandler, *lookupHandlerV2)
	router.StartAPIServer()
}
