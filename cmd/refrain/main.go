package main

import (
	"fmt"
	"log"

	"github.com/gerald-lbn/refrain/internal/config"
	"github.com/gerald-lbn/refrain/internal/container"
)

func main() {
	cfg, err := config.Load("/config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	c := container.New(cfg)

	fmt.Printf("Refrain Backend Initialized.\n")
	if c.Config != nil {
		fmt.Printf("Loaded %d libraries\n", len(c.Config.Libraries))
		for _, lib := range c.Config.Libraries {
			fmt.Printf("- %s: %s\n", lib.Name, lib.Path)
		}
	}
}
