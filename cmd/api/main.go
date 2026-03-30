package main

import (
	"fmt"
	"log"
	"os"

	"belajar_golang/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Config loaded. Starting server on port :" + cfg.AppPort)
}
