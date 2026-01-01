package main

import (
	"fmt"
	"log"
	"net/http"

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
	config.LoadConfig("config.json")
	globalConfig = config.GetGlobalConfig()

	if err := database.LoadDictionaries("resources"); err != nil {
		log.Fatal("failed to load dictionaries:", err)
	}
	dictionaries = database.GetDictionaries()

	lookup.ApplyPriorityOrder(globalConfig, dictionaries)

	log.Printf("Loaded %d dictionaries\n", len(dictionaries))

	http.Handle("/lookup", api.CorsMiddleware(http.HandlerFunc(
		api.LookupHandler(dictionaries, globalConfig),
	)))

	fmt.Println("Listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
