package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"

	_ "modernc.org/sqlite"

	"dinhphu28.com/dictionary/internal/config"
	"dinhphu28.com/dictionary/internal/database"
	"dinhphu28.com/dictionary/internal/lookup"
	native "dinhphu28.com/dictionary/internal/native"
	"dinhphu28.com/dictionary/internal/startup"
)

func main() {
	// 🔒 CRITICAL: never write logs to stdout
	log.SetOutput(os.Stderr)
	log.Println("Native host started")

	configPath := startup.ResolvePath("config.json")
	resourcesPath := startup.ResolvePath("resources")

	if err := config.LoadConfig(configPath); err != nil {
		log.Println("failed to load config:", err)
	}
	globalConfig := config.GetGlobalConfig()
	log.Printf("Config loaded: %+v\n", globalConfig)

	if err := database.LoadDictionaries(resourcesPath); err != nil {
		log.Println("failed to load dictionaries:", err)
	}
	dictionaries := database.GetDictionaries()
	log.Printf("Loaded %d dictionaries\n", len(dictionaries))

	dictionaryLookup := lookup.NewDictionaryLookup(dictionaries, globalConfig)
	approximateLookup := lookup.NewApproximateLookup(*dictionaryLookup)

	ready := len(dictionaries) > 0

	for {
		raw, err := native.ReadMessage()
		if err != nil {
			if err == io.EOF {
				log.Println("Chrome disconnected, exiting")
				return
			}
			log.Printf("read error: %v", err)
			return
		}

		var req native.Request
		if err := json.Unmarshal(raw, &req); err != nil {
			log.Printf("bad request: %v", err)
			_ = native.WriteMessage(native.Response{
				Type:    native.Error,
				Message: "invalid request",
			})
			continue
		}

		log.Printf("received: %+v", req)

		switch req.Type {

		case native.Ping:
			_ = native.WriteMessage(native.Response{
				Type:    native.Pong,
				Ready:   ready,
				Message: "Dictionaries loaded: " + strconv.Itoa(len(dictionaries)),
			})

		case native.Lookup:
			// 🔁 TEMP: fake result to prove Chrome works
			result, err := approximateLookup.LookupWithSuggestion(req.Query)
			if err != nil {
				_ = native.WriteMessage(native.Response{
					Type:    native.Error,
					Message: "lookup error: " + err.Error(),
				})
				continue
			}
			_ = native.WriteMessage(native.Response{
				Type:   native.Result,
				Ready: true,
				Query:  req.Query,
				Result: result,
			})

		default:
			_ = native.WriteMessage(native.Response{
				Type:    native.Error,
				Message: "unknown message type",
			})
		}
	}
}
