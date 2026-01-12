package main

import (
	"fmt"
	"log"

	"dinhphu28.com/dictionary/internal/setup"
)

func main() {
	fmt.Println("Installing dictionary...")

	paths := setup.DefaultPaths()

	if err := setup.Install(paths); err != nil {
		log.Fatalf("install failed: %v", err)
	}

	fmt.Println("✅ Installation complete")
	fmt.Println("Make sure ~/.local/bin is in your PATH")
}
