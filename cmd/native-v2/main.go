package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"dinhphu28.com/dictionary/internal/native"
)

func main() {
	log.SetOutput(os.Stderr)
	log.Println("native host started (persistent)")

	for {
		raw, err := native.ReadMessage()
		if err != nil {
			if err == io.EOF {
				log.Println("chrome disconnected, exiting")
				return
			}
			log.Printf("read error: %v", err)
			return
		}

		log.Printf("received raw: %s", string(raw))

		var req map[string]any
		_ = json.Unmarshal(raw, &req)

		resp := map[string]any{
			"ok":      true,
			"echo":    req,
			"message": "Hello from persistent native host",
			"pid":     os.Getpid(),
		}

		if err := native.WriteMessage(resp); err != nil {
			log.Printf("write error: %v", err)
			return
		}

		log.Println("response sent, waiting for next request")
	}
}
