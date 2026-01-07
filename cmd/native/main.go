package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dinhphu28.com/dictionary/internal/config"
	"dinhphu28.com/dictionary/internal/database"
	"dinhphu28.com/dictionary/internal/lookup"
	"dinhphu28.com/dictionary/internal/native"

	_ "modernc.org/sqlite"
)

var (
	globalConfig config.GlobalConfig
	dictionaries []database.Dictionary
)

func main() {
	// NOTE: stdout is protocol → log to stderr instead
	log.SetOutput(os.Stderr)
	log.Println("Native host starting...")

	// NOTE: allow Ctrl+C / extension disconnect graceful exit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		start := time.Now()
		log.Println("loading config and dictionaries...")
		if err := config.LoadConfig("config.json"); err != nil {
			log.Printf("failed to load config: %v", err)
			native.WriteMessage(native.Response{
				Type:    native.Error,
				Message: err.Error(),
			})
			os.Exit(1)
		}
		globalConfig = config.GetGlobalConfig()
		log.Printf("config loaded: %+v", globalConfig)

		if err := database.LoadDictionaries("resources"); err != nil {
			log.Printf("failed to load dictionaries: %v", err)
			native.WriteMessage(native.Response{
				Type:    native.Error,
				Message: err.Error(),
			})
			os.Exit(1)
		}
		dictionaries = database.GetDictionaries()
		log.Printf("dictionaries loaded: %d", len(dictionaries))
		log.Printf("startup ready in: %s", time.Since(start))
		native.WriteMessage(native.Response{
			Type:  native.Loading,
			Ready: true,
		})
	}()

	for {
		select {
		case <-sig:
			log.Println("shutting down native host")
			return
		default:
			raw, err := native.ReadMessage()
			if err != nil {
				if err == io.EOF {
					log.Println("pipe closed - exitting")
					return
				}
				log.Printf("read error: %v", err)
				return
			}
			log.Printf("received message: %+v", raw)
			var req native.Request
			_ = json.Unmarshal(raw, &req)
			log.Printf("parsed request: %+v", req)
			ready := false
			if len(dictionaries) > 0 {
				ready = true
			}
			switch req.Type {
			case native.Ping:
				native.WriteMessage(native.Response{
					Type:  native.Pong,
					Ready: ready,
				})
			case native.Lookup:
				if !ready {
					native.WriteMessage(native.Response{
						Type:    native.Loading,
						Ready:   false,
						Message: "dictionaries are still loading",
					})
					continue
				}
				approximateLookup := lookup.NewApproximateLookup(dictionaries, globalConfig)
				result, err := approximateLookup.LookupWithSuggestion(req.Query)
				if err != nil {
					native.WriteMessage(native.Response{
						Type:    native.Error,
						Query:   req.Query,
						Message: err.Error(),
					})
					continue
				}
				native.WriteMessage(native.Response{
					Type:   native.Result,
					Query:  req.Query,
					Ready:  true,
					Result: result,
				})

			default:
				native.WriteMessage(native.Response{
					Type:    native.Error,
					Message: "unknown message type",
				})
			}
		}
	}
}
