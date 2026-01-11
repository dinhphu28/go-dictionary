package main

import (
	"fmt"
	"log"
	"os"
	"dinhphu28.com/dictionary/internal/setup"
)

func main() {
	log.SetOutput(os.Stderr)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "setup":
			if err := setup.Run(); err != nil {
				log.Fatal(err)
			}
			return

		case "native":
			runNative()
			return

		case "http":
			runHTTP()
			return
		}
	}

	// Default behavior
	runNative()
}

func runNative() {
	// your existing native host logic
	fmt.Println("Native mode")
}

func runHTTP() {
	// your existing HTTP server logic
	fmt.Println("HTTP mode")
}
