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
)

// ---------- Message schema (keep simple for Chrome test) ----------

// type Request struct {
// 	Type  string `json:"type"`
// 	Query string `json:"query,omitempty"`
// }
//
// type Response struct {
// 	Type   string      `json:"type"`
// 	Query  string      `json:"query,omitempty"`
// 	Result interface{} `json:"result,omitempty"`
// 	Error  string      `json:"error,omitempty"`
// }

func main() {
	// 🔒 CRITICAL: never write logs to stdout
	log.SetOutput(os.Stderr)
	log.Println("Native host started")

	if err := config.LoadConfig("/home/dinhphu28/ghq/github.com/dinhphu28/go-dictionary/config.json"); err != nil {
		log.Println("failed to load config:", err)
	}
	globalConfig := config.GetGlobalConfig()
	log.Printf("Config loaded: %+v\n", globalConfig)

	if err := database.LoadDictionaries("/home/dinhphu28/ghq/github.com/dinhphu28/go-dictionary/resources"); err != nil {
		log.Println("failed to load dictionaries:", err)
	}
	dictionaries := database.GetDictionaries()
	log.Printf("Loaded %d dictionaries\n", len(dictionaries))

	approximateLookup := lookup.NewApproximateLookup(dictionaries, globalConfig)

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
