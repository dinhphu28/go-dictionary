package main

import (
	"fmt"
	"log"
	"net/http"

	"dinhphu28.com/dictionary/internal"
	"dinhphu28.com/dictionary/internal/api"
	_ "modernc.org/sqlite"
)

// ---- global list of loaded dictionaries ----

var globalConfig internal.GlobalConfig

// ---- Start of code ----
var dictionaries []internal.Dictionary

func main() {
	internal.LoadConfig("config.json")
	globalConfig = internal.GetGlobalConfig()

	if err := internal.LoadDictionaries("resources"); err != nil {
		log.Fatal("failed to load dictionaries:", err)
	}
	dictionaries = internal.GetDictionaries()

	internal.ApplyPriorityOrder(globalConfig, dictionaries)

	log.Printf("Loaded %d dictionaries\n", len(dictionaries))

	http.Handle("/lookup", api.CorsMiddleware(http.HandlerFunc(
		api.LookupHandler(dictionaries, globalConfig),
	)))

	fmt.Println("Listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
